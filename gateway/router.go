package main

import (
	"github.com/cushydigit/nanobank/shared/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Routes() http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// specify who is allowed to connect
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5134",  // Dev frontend
			"https://microstore.com", // Prod frontend
		},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{
			"Authorization",
			"X-User-ID",
		},
		AllowCredentials: true,
		MaxAge:           300, // seconds
	}))

	// routes
	// auth service
	r.Post("/login", utils.ProxyHandler(API_URL_AUTH))
	r.Post("/register", utils.ProxyHandler(API_URL_AUTH))

	return r

}
