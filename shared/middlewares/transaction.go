package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/types"
)

func ProvideCreateTransactionReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateTransactionReqBody
		if err := helpers.ReadJSON(w, r, &req); err != nil {
			helpers.ErrorJSON(w, errors.New("invalid request"))
			return
		}

		ctx := context.WithValue(r.Context(), types.CreateTransactionReqKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ProvideUpdateTransactionReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateTransactionReqBody
		if err := helpers.ReadJSON(w, r, &req); err != nil {
			helpers.ErrorJSON(w, errors.New("invalid request"))
			return
		}

		ctx := context.WithValue(r.Context(), types.UpdateTransactionReqKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
