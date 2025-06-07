package messaging

import (
	"encoding/json"

	"github.com/cushydigit/nanobank/shared/internalmq"
	"github.com/rabbitmq/amqp091-go"
)

func PublishNotifaction(mq *internalmq.RabbitMQClient, queueName string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = mq.Channel.Publish("", queueName, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	},
	)

	return err

}
