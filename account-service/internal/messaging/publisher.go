package messaging

import (
	"encoding/json"
	"log"

	"github.com/cushydigit/nanobank/shared/internalmq"
	"github.com/rabbitmq/amqp091-go"
)

func PublishNotifaction(mq *internalmq.RabbitMQClient, queueName string, payload any) error {
	q, err := mq.DeclareQueue(queueName)
	if err != nil {
		log.Fatal(err)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = mq.Channel.Publish("", q.Name, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	},
	)

	return err

}
