package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cushydigit/nanobank/shared/database"
	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/middlewares"

	"github.com/cushydigit/nanobank/transaction-service/internal/handler"
	"github.com/cushydigit/nanobank/transaction-service/internal/repository"
	"github.com/cushydigit/nanobank/transaction-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	PORT = os.Getenv("PORT")
	DNS  = os.Getenv("DNS")
)

func main() {

	// check environment variables
	if PORT == "" || DNS == "" {
		log.Fatal("wrong environment variable")
	}

	// connect DB
	db := database.ConnectDB(DNS)

	// create repo
	r := repository.NewPostgresTransactionRepository(db)
	// create service
	s := service.NewTransactionService(r)
	// create handler
	h := handler.NewTransactionHandler(s)

	// create router mux
	m := chi.NewRouter()

	// setup global middlwares
	m.Use(middleware.Logger)
	m.Use(middleware.Recoverer)
	m.Use(middleware.Heartbeat("/ping"))

	// setup routes
	m.Get("/", h.List)
	m.With(middlewares.ProvideCreateTransactionReq).Post("/internal/create", h.Create)
	m.With(middlewares.ProvideUpdateTransactionReq).Post("/internal/update", h.Create)
	m.Get("/{id}", h.GetByID)

	// not allowed and not found handlers
	m.NotFound(func(w http.ResponseWriter, r *http.Request) {
		helpers.ErrorJSON(w, errors.New("route not found"), http.StatusNotFound)
	})

	m.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		helpers.ErrorJSON(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
	})

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: m,
	}

	log.Printf("starting account service on: %s", PORT)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
