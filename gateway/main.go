package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	API_URL_AUTH = os.Getenv("API_URL_AUTH")
	PORT         = os.Getenv("PORT")
)

func main() {

	if PORT == "" || API_URL_AUTH == "" {
		log.Fatal("wrong environment variable")
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: Routes(),
	}

	log.Printf("starting gateway service on: %s\n", PORT)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
