package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	myredis "github.com/cushydigit/nanobank/shared/redis"
)

var (
	PORT          = os.Getenv("PORT")
	API_URL_AUTH  = os.Getenv("API_URL_AUTH")
	API_URL_REDIS = os.Getenv("API_URL_REDIS")
)

func main() {

	if PORT == "" || API_URL_AUTH == "" || API_URL_REDIS == "" {
		log.Fatal("wrong environment variable")
	}

	myredis.Init(context.Background(), API_URL_REDIS)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: Routes(),
	}

	log.Printf("starting gateway service on: %s\n", PORT)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
