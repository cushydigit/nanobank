package utils

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func GenerateUserID(prefix string) string {
	first := 1
	rest := ""
	for range 7 {
		rest += strconv.Itoa(rand.Intn(10))
	}
	return fmt.Sprintf("%s-%d%s", prefix, first, rest)
}

func GenerateAccountNumber(prefix string) string {
	first := 9910
	rest := ""
	for range 12 {
		rest += strconv.Itoa(rand.Intn(10))
	}
	return fmt.Sprintf("%s-%d%s", prefix, first, rest)
}

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

func ProxyHandler(target string) http.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	return func(w http.ResponseWriter, r *http.Request) {
		r.Host = url.Host
		proxy.ServeHTTP(w, r)
	}
}
