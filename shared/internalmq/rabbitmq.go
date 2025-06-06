package internalmq

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// this structure works fine if a service is just a listener or consumer
// for cases that the service is listener and consumer it should use seprate channel using mustex

type RabbitMQClient struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQClient(amqpURL string) (*RabbitMQClient, error) {
	conn, err := retryConnection(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQClient{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (r *RabbitMQClient) DeclareQueue(queueName string) (amqp.Queue, error) {
	return r.Channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
}

func (r *RabbitMQClient) Close() {
	_ = r.Channel.Close()
	_ = r.Conn.Close()
}

func retryConnection(url string) (*amqp.Connection, error) {
	var err error
	var conn *amqp.Connection
	for count := range 10 {
		conn, err = amqp.Dial(url)
		if err != nil {
			log.Printf("Retrying to connect ot RabbitMQ after %v.. Attempt %d", 5*time.Second, count)
			time.Sleep(5 * time.Second)
		} else {
			return conn, err
		}
	}
	return conn, err

}
