module github.com/cushydigit/nanobank/mailer-service

go 1.24.2

require (
	github.com/cushydigit/nanobank/shared v0.1.0
	github.com/go-chi/chi/v5 v5.2.1
)

require github.com/golang-jwt/jwt/v5 v5.2.2 // indirect

replace github.com/cushydigit/nanobank/shared => ../shared/
