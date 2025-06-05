package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cushydigit/nanobank/shared/database"
	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/middlewares"
	myredis "github.com/cushydigit/nanobank/shared/redis"

	"github.com/cushydigit/nanobank/transaction-service/internal/handler"
	"github.com/cushydigit/nanobank/transaction-service/internal/repository"
	"github.com/cushydigit/nanobank/transaction-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	PORT          = os.Getenv("PORT")
	DNS           = os.Getenv("DNS")
	ROOT_EMAIL    = os.Getenv("ROOT_EMAIL")
	API_URL_REDIS = os.Getenv("API_URL_REDIS")
)

func main() {

	ctx := context.Background()

	// check environment variables
	if PORT == "" || DNS == "" || ROOT_EMAIL == "" || API_URL_REDIS == "" {
		log.Fatal("wrong environment variable")
	}

	// connect DB
	db := database.ConnectDB(DNS)

	// init redis (cacher) client
	c := myredis.MyRedisClientInit(ctx, API_URL_REDIS)

	// create repo
	r := repository.NewPostgresTransactionRepository(db)
	// create service
	s := service.NewTransactionService(r, c)
	// create handler
	h := handler.NewTransactionHandler(s)

	// create router mux
	m := chi.NewRouter()

	// setup global middlwares
	m.Use(middleware.Logger)
	m.Use(middleware.Recoverer)
	m.Use(middleware.Heartbeat("/ping"))

	// setup routes
	// internal routes
	m.Get("/internal/{id}", h.GetByID)
	m.With(middlewares.ProvideCreateTransactionReq).Post("/internal", h.Create)
	m.With(middlewares.ProvideUpdateTransactionReq).Put("/internal/{id}", h.Update)
	// require auth routes
	m.With(middlewares.RequireRoot).Get("/", h.ListAll)
	m.Get("/me", h.ListByUserID)

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
