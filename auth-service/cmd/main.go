package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cushydigit/nanobank/auth-service/internal/handler"
	"github.com/cushydigit/nanobank/auth-service/internal/repository"
	"github.com/cushydigit/nanobank/auth-service/internal/service"
	"github.com/cushydigit/nanobank/shared/database"
	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/middlewares"
	myredis "github.com/cushydigit/nanobank/shared/redis"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	JWT_SECRET    = os.Getenv("JWT_SECRET")
	PORT          = os.Getenv("PORT")
	DNS           = os.Getenv("DNS")
	ROOT_EMAIL    = os.Getenv("ROOT_EMAIL")
	ROOT_PASSWORD = os.Getenv("ROOT_PASSWORD")
	API_URL_REDIS = os.Getenv("API_URL_REDIS")
)

func main() {

	ctx := context.Background()

	// check environment variables
	if PORT == "" || DNS == "" || ROOT_EMAIL == "" || ROOT_PASSWORD == "" || JWT_SECRET == "" || API_URL_REDIS == "" {
		log.Fatal("wrong environment variable")
	}
	// init redis (cacher) client
	c := myredis.MyRedisClientInit(ctx, API_URL_REDIS)

	// connect DB
	db := database.ConnectDB(DNS)

	// create repo
	r := repository.NewPostgresUserRepository(db)
	// create service
	s := service.NewAuthService(r, c)
	// create handler
	h := handler.NewAuthHandler(s)

	// init the root user or admin that has all privilages
	if admin, _ := r.FindByEmail(ctx, ROOT_EMAIL); admin == nil {
		log.Printf("the root user is not exists try to register a new root")
		if _, err := s.Register(ctx, "admin", ROOT_EMAIL, ROOT_PASSWORD); err != nil {
			log.Fatalf("failed to create root user: %v", err)
		}
	}

	// create router mux
	m := chi.NewRouter()

	// setup global middlwares
	m.Use(middleware.Logger)
	m.Use(middleware.Recoverer)
	m.Use(middleware.Heartbeat("/ping"))

	// setup routes
	m.With(middlewares.ValidateRegisterUserRequest).Post("/register", h.Register)
	m.With(middlewares.ProvideAuthRequest).Post("/login", h.Login)
	m.With(middlewares.ProvideRefreshRequest).Post("/refresh", h.Refresh)
	m.With(middlewares.ProvideRefreshRequest).Post("/logout", h.Logout)

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
