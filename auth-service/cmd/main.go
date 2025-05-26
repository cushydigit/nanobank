package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cushydigit/nanobank/shared/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	PORT = os.Getenv("PORT")
	DNS  = os.Getenv("DNS")
)

func main() {

	// connect DB
	_ = database.ConnectDB(DNS)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: r,
	}

	log.Printf("starting auth service on: %s", PORT)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
