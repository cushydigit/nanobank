package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/types"
)

func ProvideUpdateBalanceReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateBalanceReqBody
		if err := helpers.ReadJSON(w, r, &req); err != nil {
			helpers.ErrorJSON(w, errors.New("invalid request"))
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
			helpers.ErrorJSON(w, errors.New("invalid request"))
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
			helpers.ErrorJSON(w, errors.New("invalid request"))
			return
		}
		ctx := context.WithValue(r.Context(), types.ConfirmTransferReqKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
