package main

import (
	"errors"
	"net/http"

	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Routes() http.Handler {
	m := chi.NewRouter()

	// Middlewares
	m.Use(middleware.Logger)
	m.Use(middleware.Recoverer)

	// specify who is allowed to connect
	m.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5134", // Dev frontend
			"https://nanobank.com",  // Prod frontend
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
	m.Route("/api", func(r chi.Router) {
		// auth service
		r.Route("/auth", func(r chi.Router) {
			r.Mount("/", http.StripPrefix("/api/auth", utils.ProxyHandler(API_URL_AUTH)))
		})

		// account service
		r.Route("/account", func(r chi.Router) {
			r.Mount("/", http.StripPrefix("/api/account", utils.ProxyHandler(API_URL_ACCOUNT)))
		})
	})

	// not allowed and not found handlers
	m.NotFound(func(w http.ResponseWriter, r *http.Request) {
		helpers.ErrorJSON(w, errors.New("route not found"), http.StatusNotFound)
	})

	m.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		helpers.ErrorJSON(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
	})

	return m

}
