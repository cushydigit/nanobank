package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/types"
	"github.com/cushydigit/nanobank/shared/utils"
)

func ValidateRegisterUserRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReqBody
		if err := helpers.ReadJSON(w, r, &req); err != nil {
			helpers.ErrorJSON(w, errors.New("invalid request"))
			return
		}

		// basic validatons
		req.Username = strings.TrimSpace(req.Username)
		req.Email = strings.TrimSpace(req.Email)

		if req.Username == "" || req.Email == "" {
			helpers.ErrorJSON(w, errors.New("username and email are required"))
			return
		}

		if !utils.IsValidUsername(req.Username) {
			helpers.ErrorJSON(w, errors.New("invalid username, only letters, numbers, and dots allowed"))
			return
		}

		if !utils.IsValidEmail(req.Email) {
			helpers.ErrorJSON(w, fmt.Errorf("invalid email: %s", req.Email))
			return
		}

		if len(req.Password) < 6 {
			helpers.ErrorJSON(w, errors.New("passwrod must be at least 6 characters long"))
			return
		}

		// inject validated request
		ctx := context.WithValue(r.Context(), types.RegisterReqKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func ProvideAuthRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req types.AuthReqBody
		if err := helpers.ReadJSON(w, r, &req); err != nil {
			helpers.ErrorJSON(w, errors.New("invalid request"))
			return
		}
		// inject validated request
		ctx := context.WithValue(r.Context(), types.AuthReqKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ProvideRefreshRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshReqBody
		if err := helpers.ReadJSON(w, r, &req); err != nil {
			helpers.ErrorJSON(w, errors.New("invalid request"))
			return
		}
		// inject validated request
		ctx := context.WithValue(r.Context(), types.RefreshReqKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
