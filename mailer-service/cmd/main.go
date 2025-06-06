package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/cushydigit/nanobank/shared/helpers"
)

var (
	// JWT_SECRET    = os.Getenv("JWT_SECRET")
	PORT = os.Getenv("PORT")
	// DNS           = os.Getenv("DNS")
	// ROOT_EMAIL    = os.Getenv("ROOT_EMAIL")
	// ROOT_PASSWORD = os.Getenv("ROOT_PASSWORD")
	// API_URL_REDIS = os.Getenv("API_URL_REDIS")
)

func main() {

	// ctx := context.Background()

	// check environment variables
	// if PORT == "" || DNS == "" || ROOT_EMAIL == "" || ROOT_PASSWORD == "" || JWT_SECRET == "" || API_URL_REDIS == "" {
	// 	log.Fatal("wrong environment variable")
	// }
	// init redis (cacher) client
	// c := myredis.MyRedisClientInit(ctx, API_URL_REDIS)

	// connect DB
	// db := database.ConnectDB(DNS)

	// create repo
	// r := repository.NewPostgresUserRepository(db)
	// create service
	// s := service.NewAuthService(r, c)
	// create handler
	// h := handler.NewAuthHandler(s)

	// create router mux
	m := chi.NewRouter()

	// setup global middlwares
	m.Use(middleware.Logger)
	m.Use(middleware.Recoverer)
	m.Use(middleware.Heartbeat("/ping"))

	// setup routes
	// m.Post("/send", h.Register)

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

	log.Printf("starting auth service on: %s", PORT)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
