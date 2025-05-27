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
	r.Route("/api", func(r chi.Router) {
		// auth service
		r.Route("/auth", func(r chi.Router) {
			r.Mount("/", http.StripPrefix("/api/auth", utils.ProxyHandler(API_URL_AUTH)))
		})
	})

	return r

}
