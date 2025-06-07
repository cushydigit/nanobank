package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cushydigit/nanobank/mailer-service/internal/service"
	"github.com/cushydigit/nanobank/shared/internalmq"
	"github.com/cushydigit/nanobank/shared/types"
)

func ListenForNotificatin(s *service.MailService, mq *internalmq.RabbitMQClient) error {

	msgs, err := mq.Channel.Consume(internalmq.QUEUE_NOTIFICATION_BALANCE, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var payload types.BalanceChangePayload
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				log.Printf("invalid notification payload: %v", err)
				continue
			}

			log.Printf("sending email for balance update...")
			subject := fmt.Sprintf("Your account has been %s", payload.Type)
			body := fmt.Sprintf("Hi %s,\n\nAn amount of %d has been %s in your account.\n\nCheers.", payload.Username, payload.Amount, payload.Type)
			if err := s.SendSMPTMessage("no-reply@nanobank.com", payload.Email, subject, body); err != nil {
				log.Printf("failed to email the notification: %v", err)
				continue
			}
		}
	}()

	return nil
}
