package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cushydigit/nanobank/shared/config"
	myredis "github.com/cushydigit/nanobank/shared/redis"
)

var (
	PORT                = os.Getenv("PORT")
	API_URL_AUTH        = os.Getenv("API_URL_AUTH")
	API_URL_ACCOUNT     = os.Getenv("API_URL_ACCOUNT")
	API_URL_TRANSACTION = os.Getenv("API_URL_TRANSACTION")
	API_URL_REDIS       = os.Getenv("API_URL_REDIS")
)

func main() {

	if PORT == "" || API_URL_AUTH == "" || API_URL_REDIS == "" || API_URL_ACCOUNT == "" || API_URL_TRANSACTION == "" {
		log.Fatal("wrong environment variable")
	}

	_ = myredis.MyRedisClientInit(context.Background(), API_URL_REDIS)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", PORT),
		Handler:           Routes(),
		IdleTimeout:       config.TIMEOUT_GATEWAY_IDLE, // Keep-alive tcp
		ReadTimeout:       config.TIMEOUT_GATEWAY_READ,
		WriteTimeout:      config.TIMEOUT_GATEWAY_WRITE,
		ReadHeaderTimeout: config.TIMEOUT_GATEWAY_READ_HEADER,
	}

	log.Printf("starting gateway service on: %s\n", PORT)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
