module github.com/cushydigit/nanobank/account-service

go 1.24.2

replace github.com/cushydigit/nanobank/shared => ../shared

require (
	github.com/cushydigit/nanobank/shared v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.2.1
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.38.0 // indirect
)
