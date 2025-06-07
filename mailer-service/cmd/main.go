package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/cushydigit/nanobank/mailer-service/internal/handler"
	"github.com/cushydigit/nanobank/mailer-service/internal/messaging"
	"github.com/cushydigit/nanobank/mailer-service/internal/service"
	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/internalmq"
	"github.com/cushydigit/nanobank/shared/middlewares"
)

var (
	MQ_DNS          = os.Getenv("MQ_DNS")
	PORT            = os.Getenv("PORT")
	API_URL_MAILHOG = os.Getenv("API_URL_MAILHOG")
)

func main() {

	//	check environment variables
	if PORT == "" || MQ_DNS == "" || API_URL_MAILHOG == "" {
		log.Fatal("wrong environment variable")
	}

	// create rabbitmq client
	mq, err := internalmq.NewRabbitMQClient("amqp://admin:admin@rabbitmq:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer mq.Close()
	// define the queue
	mq.DeclareQueue(internalmq.QUEUE_NOTIFICATION_BALANCE)

	// create service
	s := service.NewMailService(API_URL_MAILHOG)
	// create handler
	h := handler.NewMailHandler(s)

	// listen for upcoming messages
	if err := messaging.ListenForNotificatin(s, mq); err != nil {
		log.Fatalf("failed to start listening: %v", err)
	}

	// create router mux
	m := chi.NewRouter()

	// setup global middlwares
	m.Use(middleware.Logger)
	m.Use(middleware.Recoverer)
	m.Use(middleware.Heartbeat("/ping"))

	// setup routes
	m.With(middlewares.ProvideSendMailReq).Post("/send", h.SendMail)

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
