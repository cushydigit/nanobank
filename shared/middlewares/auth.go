package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/types"
	"github.com/cushydigit/nanobank/shared/utils"

	myerrors "github.com/cushydigit/nanobank/shared/errors"
)

var (
	ROOT_EMAIL = os.Getenv("ROOT_EMAIL")
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

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// require env variable set correctly

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helpers.ErrorJSON(w, myerrors.ErrMissingAuthorizationHeader, http.StatusUnauthorized)
			return
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		if accessToken == authHeader {
			helpers.ErrorJSON(w, myerrors.ErrInvalidTokenFormat, http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateToken(accessToken)
		if err != nil {
			helpers.ErrorJSON(w, myerrors.ErrInvalidOrExpiredToken, http.StatusUnauthorized)
			return
		}

		// inject into header for downstream services
		r.Header.Set(string(types.XUserID), claims.UserID)
		r.Header.Set(string(types.XUserEmail), claims.Email)
		r.Header.Set(string(types.XUsername), claims.Username)

		// inject into context for internal use
		ctx := context.WithValue(r.Context(), types.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, types.UserEmailKey, claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// this middleware require the RequireAuth to come before this middleware
func RequireRoot(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check enironment vairable
		if ROOT_EMAIL == "" {
			log.Fatalf("wrong environment variable for require root middleware")
		}

		email, ok := r.Context().Value(types.UserEmailKey).(string)
		if !ok {
			log.Fatalf("the require root middleware must come after the require auth middleware")
		}

		if email != ROOT_EMAIL {
			helpers.ErrorJSON(w, myerrors.ErrPermissionDenied, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
