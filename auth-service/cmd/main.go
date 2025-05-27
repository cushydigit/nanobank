package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	service "github.com/cushydigit/nanobank/auth-service/internal"
	"github.com/cushydigit/nanobank/auth-service/internal/handler"
	"github.com/cushydigit/nanobank/auth-service/internal/repository"
	"github.com/cushydigit/nanobank/shared/database"
	"github.com/cushydigit/nanobank/shared/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	JWT_SECRET    = os.Getenv("JWT_SECRET")
	PORT          = os.Getenv("PORT")
	DNS           = os.Getenv("DNS")
	ROOT_EMAIL    = os.Getenv("ROOT_EMAIL")
	ROOT_PASSWORD = os.Getenv("ROOT_PASSWORD")
)

func main() {

	// check environment variables
	if PORT == "" || DNS == "" || ROOT_EMAIL == "" || ROOT_PASSWORD == "" || JWT_SECRET == "" {
		log.Fatal("wrong environment variable")
	}

	// connect DB
	db := database.ConnectDB(DNS)

	// create repo
	r := repository.NewPostgresUserRepository(db)
	// create service
	s := service.NewAuthService(r)
	// create handler
	h := handler.NewAuthHandler(s)

	// init the root user or admin that has all privilages
	if admin, _ := r.FindByEmail(ROOT_EMAIL); admin == nil {
		log.Printf("the root user is not exists try to register a new root")
		if _, err := s.Register("admin", ROOT_EMAIL, ROOT_PASSWORD); err != nil {
			log.Fatalf("failed to create root user: %v", err)
		}
	}

	// create router mux
	m := chi.NewRouter()

	// setup global middlwares
	m.Use(middleware.Logger)
	m.Use(middleware.Recoverer)
	m.Use(middleware.Heartbeat("/ping"))

	// int routes routes
	m.With(middlewares.ValidateRegisterUserRequest).Post("/register", h.Register)
	m.Post("/login", h.Login)
	m.Post("/refresh", h.Refresh)

	m.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
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
