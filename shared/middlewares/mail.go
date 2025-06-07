package middlewares

import (
	"context"
	"net/http"

	myerrors "github.com/cushydigit/nanobank/shared/errors"
	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/types"
)

func ProvideSendMailReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var req types.SendMailReqBody
		if err := helpers.ReadJSON(w, r, &req); err != nil {
			helpers.ErrorJSON(w, myerrors.ErrInvalidRequest)
			return
		}

		ctx := context.WithValue(r.Context(), types.SendMailReqKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
