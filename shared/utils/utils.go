package utils

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// isValidEmail checks if the email has a valid format
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// isValidUsername allows alphanumerics and dots, with no leading/trailing/consecutive dots
func IsValidUsername(username string) bool {
	// Must start and end with letter or number, dots allowed in between
	re := regexp.MustCompile(`^[a-zA-Z0-9]+(?:\.[a-zA-Z0-9]+)*$`)
	return re.MatchString(username)
}

func ProxyHandler(target, block string) http.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	return func(w http.ResponseWriter, r *http.Request) {
		// block routes
		if strings.HasPrefix(strings.ToLower(r.URL.Path), block) {
			helpers.ErrorJSON(w, errors.New("forbidden"), http.StatusForbidden)
		}
		r.Host = url.Host
		proxy.ServeHTTP(w, r)
	}
}

func GenerateTransactionToken() (string, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return uid.String(), nil
}
