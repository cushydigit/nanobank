package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cushydigit/nanobank/account-service/internal/repository"
	myerrors "github.com/cushydigit/nanobank/shared/errors"
	"github.com/cushydigit/nanobank/shared/models"
	"github.com/cushydigit/nanobank/shared/types"
)

type AccountService struct {
	repo                repository.AccountRepository
	API_URL_TRANSACTION string
}

func NewAccountService(r repository.AccountRepository, url string) *AccountService {
	return &AccountService{
		repo:                r,
		API_URL_TRANSACTION: url,
	}
}

// returns Errs: ErrAccountNotFound, ErrInternalServer
func (s *AccountService) Get(ctx context.Context, userID string) (*models.Account, error) {
	a, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		// the account not found
		if err == sql.ErrNoRows {
			return nil, myerrors.ErrAccountNotFound
		} else {
			log.Printf("unexpected error: %v", err)
			return nil, myerrors.ErrInternalServer
		}
	}
	return a, nil
}

// returns Errs: ErrAccountAlreadyExists, ErrInternalServer
func (s *AccountService) Create(ctx context.Context, userID, username string) (*models.Account, error) {
	if exists, _ := s.repo.FindByUserID(ctx, userID); exists != nil {
		return nil, myerrors.ErrAccountAlreadyExists
	}
	newAccount := models.NewAccount(userID, username)
	if err := s.repo.Create(ctx, newAccount); err != nil {
		log.Printf("unexpected error: %v", err)
		return nil, myerrors.ErrInternalServer
	}
	return newAccount, nil
}

// returns Errs: ErrAmountMustBePositive, ErrAccountNotFound, ErrInternalServer
func (s *AccountService) Deposit(ctx context.Context, userID string, amount int64) error {
	if amount <= 0 {
		return myerrors.ErrAmountMustBePositive
	}
	if _, err := s.repo.FindByUserID(ctx, userID); err != nil {
		if err == sql.ErrNoRows {
			return myerrors.ErrAccountNotFound
		} else {
			log.Printf("unexpected error: %v", err)
			return myerrors.ErrInternalServer
		}
	}
	if err := s.repo.UpdateBalance(ctx, userID, amount); err != nil {
		log.Printf("unexpected error: %v", err)
		return err
	}
	return nil
}

// returns Errs: ErrAmountMustBePositive, ErrAccountNotFound, ErrInsufficientBalance, ErrInternalServer
func (s *AccountService) Withdraw(ctx context.Context, userID string, amount int64) error {
	if amount <= 0 {
		return myerrors.ErrAmountMustBePositive
	}
	a, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return myerrors.ErrAccountNotFound
		} else {
			log.Printf("unexpected error: %v", err)
			return myerrors.ErrInternalServer
		}
	}

	if a.Balance < amount {
		return myerrors.ErrInsufficientBalance
	}
	if err := s.repo.UpdateBalance(ctx, userID, -amount); err != nil {
		log.Printf("unexpected error: %v", err)
		return err
	}
	return nil
}

// returns toAccount (desitination account user) and a transaction with pending, errs: ErrAmountMustBePositive, ErrInsufficientBalance, ErrAccountNotFound, ErrDestinationAccountNotFound, ErrInternalServer
func (s *AccountService) InitiateTransfer(ctx context.Context, fromUserID, toUserID string, amount int64) (*models.Account, *models.Transaction, error) {

	if amount <= 0 {
		return nil, nil, myerrors.ErrAmountMustBePositive
	}

	// check the from account
	from, err := s.repo.FindByUserID(ctx, fromUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, myerrors.ErrAccountNotFound
		} else {
			log.Printf("unexpected err: %v", err)
			return nil, nil, myerrors.ErrInternalServer
		}
	}

	// check the from account balance
	if amount > from.Balance {
		return nil, nil, myerrors.ErrInsufficientBalance
	}

	// check the to account
	to, err := s.repo.FindByUserID(ctx, toUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, myerrors.ErrDestinationAccountNotFound
		} else {
			log.Printf("unexpected err: %v", err)
			return nil, nil, myerrors.ErrInternalServer
		}
	}

	// create a request to transaction service for creating a new transaction
	body := types.CreateTransactionReqBody{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, nil, myerrors.ErrInternalServer
	}

	resp, err := http.Post(fmt.Sprintf("%s/internal", s.API_URL_TRANSACTION), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, nil, myerrors.ErrInternalServer
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("unexpected status code in creating the trnasaction: %d", resp.StatusCode)
		return nil, nil, myerrors.ErrInternalServer
	}

	var res types.Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, nil, myerrors.ErrInternalServer
	}

	dataBytes, err := json.Marshal(res.Data) // convert map to json
	if err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, nil, myerrors.ErrInternalServer
	}

	var t models.Transaction
	if err := json.Unmarshal(dataBytes, &t); err != nil { // convert json to struct
		log.Printf("unexpected err: %v", err)
		return nil, nil, myerrors.ErrInternalServer
	}

	return to, &t, nil
}

// returns ErrInternalServer, ErrConfirmationTokenIsNotValid, ErrAccountNotFound, ErrDestinationAccountNotFound, ErrInsufficientBalance
func (s *AccountService) ConfirmTransfer(ctx context.Context, txID, token string) error {
	resp, err := http.Get(fmt.Sprintf("%s/internal/%s", s.API_URL_TRANSACTION, txID))
	if err != nil {
		log.Printf("unexpected err: %v", err)
		return myerrors.ErrInternalServer
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("unexpected status code in getting the transaction: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	var res types.Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Printf("unexpected err: %v", err)
		return myerrors.ErrInternalServer
	}

	dataBytes, err := json.Marshal(res.Data)
	if err != nil {
		log.Printf("unexpected err: %v", err)
		return myerrors.ErrInternalServer
	}

	var t models.Transaction
	if err := json.Unmarshal(dataBytes, &t); err != nil {
		log.Printf("unexpected err: %v", err)
		return myerrors.ErrInternalServer
	}

	// check if transaction already has some other status
	if t.Status != models.StatusPending {
		return myerrors.ErrConfirmationTokenIsNotValid
	}

	// TODO chekc if the casher has this value (is not expired)

	// check if token and txID is the same
	if token != t.ConfirmationToken {
		return myerrors.ErrConfirmationTokenIsNotValid
	}

	// check the source account
	sa, err := s.repo.FindByUserID(ctx, t.FromUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return myerrors.ErrAccountNotFound
		}
		log.Printf("unexpected err: %v", err)
		return myerrors.ErrInternalServer
	}

	// check the balance of source account
	if sa.Balance < t.Amount {
		return myerrors.ErrInsufficientBalance
	}

	// check the destination account
	_, err = s.repo.FindByUserID(ctx, t.ToUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return myerrors.ErrDestinationAccountNotFound
		}
		log.Printf("unexpected err: %v", err)
		return myerrors.ErrInternalServer
	}

	body := types.UpdateTransactionReqBody{
		ID:     t.ID,
		Status: models.StatusConfirmed,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		log.Printf("unexpected err: %v", err)
		return myerrors.ErrInternalServer
	}

	resp, err = http.Post(fmt.Sprintf("%s/internal/%s", s.API_URL_TRANSACTION, t.ID), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("unexpected err: %v", err)
		return myerrors.ErrInternalServer
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("unexpected status code: %d", resp.StatusCode)
		return myerrors.ErrInternalServer
	}

	if err := s.repo.TransferAmount(ctx, t.FromUserID, t.ToUserID, t.Amount); err != nil {
		log.Printf("unexpected err: %v", err)
		return myerrors.ErrInternalServer
	}

	return nil
}
