package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cushydigit/nanobank/account-service/internal/messaging"
	"github.com/cushydigit/nanobank/account-service/internal/repository"
	myerrors "github.com/cushydigit/nanobank/shared/errors"
	"github.com/cushydigit/nanobank/shared/internalhttp"
	"github.com/cushydigit/nanobank/shared/internalmq"
	"github.com/cushydigit/nanobank/shared/models"
	"github.com/cushydigit/nanobank/shared/redis"
	"github.com/cushydigit/nanobank/shared/types"
)

type AccountService struct {
	repo                repository.AccountRepository
	tokenCacher         redis.TokenCacher
	mq                  *internalmq.RabbitMQClient
	API_URL_TRANSACTION string
}

func NewAccountService(r repository.AccountRepository, c redis.TokenCacher, mq *internalmq.RabbitMQClient, url string) *AccountService {
	return &AccountService{
		repo:                r,
		tokenCacher:         c,
		mq:                  mq,
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
func (s *AccountService) Deposit(ctx context.Context, userID, username, email string, amount int64) error {
	if amount <= 0 {
		return myerrors.ErrAmountMustBePositive
	}
	_, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
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

	payload := types.BalanceChangePayload{
		Username: username,
		Email:    email,
		Type:     "deposit",
		Amount:   amount,
	}

	messaging.PublishNotifaction(s.mq, internalmq.QUEUE_NOTIFICATION_BALANCE, payload)

	return nil
}

// returns Errs: ErrAmountMustBePositive, ErrAccountNotFound, ErrInsufficientBalance, ErrInternalServer
func (s *AccountService) Withdraw(ctx context.Context, userID, username, email string, amount int64) error {
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

	payload := types.BalanceChangePayload{
		Username: username,
		Email:    email,
		Type:     "withdraw",
		Amount:   amount,
	}

	messaging.PublishNotifaction(s.mq, internalmq.QUEUE_NOTIFICATION_BALANCE, payload)

	return nil
}

// returns toAccount (desitination account user) and a transaction with pending, errs: ErrAmountMustBePositive, ErrInsufficientBalance, ErrAccountNotFound, ErrDestinationAccountNotFound, ErrInternalServer, ErrSelfTransder
func (s *AccountService) InitiateTransfer(ctx context.Context, fromUserID, toUserID string, amount int64) (*models.Account, string, error) {

	if amount <= 0 {
		return nil, "", myerrors.ErrAmountMustBePositive
	}

	if fromUserID == toUserID {
		return nil, "", myerrors.ErrSelfTransder
	}

	// check the from account
	from, err := s.repo.FindByUserID(ctx, fromUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", myerrors.ErrAccountNotFound
		} else {
			log.Printf("unexpected err: %v", err)
			return nil, "", myerrors.ErrInternalServer
		}
	}

	// check the from account balance
	if amount > from.Balance {
		return nil, "", myerrors.ErrInsufficientBalance
	}

	// check the to account
	to, err := s.repo.FindByUserID(ctx, toUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", myerrors.ErrDestinationAccountNotFound
		} else {
			log.Printf("unexpected err: %v", err)
			return nil, "", myerrors.ErrInternalServer
		}
	}

	// create a request to transaction service for creating a new transaction
	body := types.CreateTransactionReqBody{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
	}

	res := types.Response{}
	url := fmt.Sprintf("%s/internal", s.API_URL_TRANSACTION)

	if err := internalhttp.DoJSON(ctx, http.MethodPost, url, body, &res); err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, "", myerrors.ErrInternalServer
	}

	dataBytes, err := json.Marshal(res.Data) // convert map to json
	if err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, "", myerrors.ErrInternalServer
	}

	var t models.Transaction
	if err := json.Unmarshal(dataBytes, &t); err != nil { // convert json to struct
		log.Printf("unexpected err: %v", err)
		return nil, "", myerrors.ErrInternalServer
	}

	// setToken to cache
	if err := s.tokenCacher.SetToken(ctx, t.ConfirmationToken, t.ID); err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, "", myerrors.ErrInternalServer
	}

	return to, t.ConfirmationToken, nil
}

// returns ErrInternalServer, ErrConfirmationTokenIsNotValid, ErrAccountNotFound, ErrDestinationAccountNotFound, ErrInsufficientBalance
func (s *AccountService) ConfirmTransfer(ctx context.Context, token string) error {
	txID, err := s.tokenCacher.GetToken(ctx, token)
	if err != nil {
		return myerrors.ErrConfirmationTokenIsNotValid
	}

	res := types.Response{}
	url := fmt.Sprintf("%s/internal/%s", s.API_URL_TRANSACTION, txID)
	if err := internalhttp.DoJSON(ctx, http.MethodGet, url, nil, &res); err != nil {
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
	url = fmt.Sprintf("%s/internal/%s", s.API_URL_TRANSACTION, t.ID)
	if err := internalhttp.DoJSON(ctx, http.MethodPut, url, body, &res); err != nil {
		log.Printf("unexpected err: %v", err)
		return myerrors.ErrInternalServer
	}

	if err := s.repo.TransferAmount(ctx, t.FromUserID, t.ToUserID, t.Amount); err != nil {
		log.Printf("unexpected err: %v", err)
		return myerrors.ErrInternalServer
	}

	// delete the token from the cache
	_ = s.tokenCacher.DelToken(ctx, token)

	return nil
}
