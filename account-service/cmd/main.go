package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cushydigit/nanobank/account-service/internal/handler"
	"github.com/cushydigit/nanobank/account-service/internal/repository"
	"github.com/cushydigit/nanobank/account-service/internal/service"
	"github.com/cushydigit/nanobank/shared/database"
	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/internalmq"
	"github.com/cushydigit/nanobank/shared/middlewares"
	myredis "github.com/cushydigit/nanobank/shared/redis"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	PORT                = os.Getenv("PORT")
	DNS                 = os.Getenv("DNS")
	API_URL_TRANSACTION = os.Getenv("API_URL_TRANSACTION")
	API_URL_REDIS       = os.Getenv("API_URL_REDIS")
)

func main() {

	ctx := context.Background()

	// check environment variables
	if PORT == "" || DNS == "" || API_URL_TRANSACTION == "" || API_URL_REDIS == "" {
		log.Fatal("wrong environment variable")
	}

	// init redis (cacher) client
	c := myredis.MyRedisClientInit(ctx, API_URL_REDIS)

	mq, err := internalmq.NewRabbitMQClient("amqp://admin:admin@rabbitmq:5672/")
	if err != nil {
		log.Fatal(err)
	}

	// connect DB
	db := database.ConnectDB(DNS)

	// create repo
	r := repository.NewPostgresAccountRepository(db)
	// create service
	s := service.NewAccountService(r, c, mq, API_URL_TRANSACTION)
	// create handler
	h := handler.NewAccountHandler(s)

	// create router mux
	m := chi.NewRouter()

	// setup global middlwares
	m.Use(middleware.Logger)
	m.Use(middleware.Recoverer)
	m.Use(middleware.Heartbeat("/ping"))

	// setup routes
	m.Get("/", h.Get)
	m.Post("/", h.Create)
	m.With(middlewares.ProvideUpdateBalanceReq).Post("/deposit", h.Deposit)
	m.With(middlewares.ProvideUpdateBalanceReq).Post("/withdraw", h.Withdraw)
	m.With(middlewares.ProvideInitiateTransferReq).Post("/transfer/initiate", h.InitiateTransfer)
	m.With(middlewares.ProvideConfirmTransferReq).Post("/transfer/confirm", h.ConfirmTransfer)

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
