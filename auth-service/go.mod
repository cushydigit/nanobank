module github.com/cushydigit/nanobank/auth-service

go 1.24.2

require (
	github.com/cushydigit/nanobank/shared v0.1.0
	github.com/go-chi/chi/v5 v5.2.1
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
)

replace github.com/cushydigit/nanobank/shared => ../shared
