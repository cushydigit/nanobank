package middlewares

import (
	"context"
	"net/http"

	myerrors "github.com/cushydigit/nanobank/shared/errors"
	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/types"
)

func ProvideUpdateBalanceReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateBalanceReqBody
		if err := helpers.ReadJSON(w, r, &req); err != nil {
			helpers.ErrorJSON(w, myerrors.ErrInvalidRequest)
			return
		}

		// basic validate the req
		if req.Amount <= 0 {
			helpers.ErrorJSON(w, myerrors.ErrAmountMustBePositive)
			return
		}

		ctx := context.WithValue(r.Context(), types.UpdateBlanceReqKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ProvideInitiateTransferReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req types.InitiateTransferReqBody
		if err := helpers.ReadJSON(w, r, &req); err != nil {
			helpers.ErrorJSON(w, myerrors.ErrInvalidRequest)
			return
		}
		// baisic validatation
		if req.Amount <= 0 {
			helpers.ErrorJSON(w, myerrors.ErrAmountMustBePositive)
			return
		}
		if req.ToUserID == "" {
			helpers.ErrorJSON(w, myerrors.ErrToUserRequiredInRequest)
			return
		}
		ctx := context.WithValue(r.Context(), types.InitiateTransferReqKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ProvideConfirmTransferReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req types.ConfirmTransferReqBody
		if err := helpers.ReadJSON(w, r, &req); err != nil {
			helpers.ErrorJSON(w, myerrors.ErrInvalidRequest)
			return
		}

		if req.Token == "" {
			helpers.ErrorJSON(w, myerrors.ErrInvalidRequest)
		}

		ctx := context.WithValue(r.Context(), types.ConfirmTransferReqKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
