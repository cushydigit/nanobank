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

	q, err := mq.DeclareQueue("balance.notification")
	if err != nil {
		return err
	}

	msgs, err := mq.Channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var notif types.BalanceChangePayload
			if err := json.Unmarshal(msg.Body, &notif); err != nil {
				log.Printf("invalid notification payload: %v", err)
				continue
			}

			log.Printf("sending email for balance update...")
			subject := fmt.Sprintf("Your account has been %s", notif.Type)
			body := fmt.Sprintf("Hi %s,\n\nAn amount of %d has been %s in your account.\n\nCheers.", "dear user", notif.Amount, notif.Type)
			if err := s.SendSMPTMessage("no-reply@nanobank.com", notif.Email, subject, body); err != nil {
				log.Printf("failed to email the notification: %v", err)
				continue
			}
		}
	}()

	return nil
}
