package service

import (
	"fmt"
	"net/smtp"
)

type MailService struct {
	API_URL_MAILHOG string
}

func NewMailService(url string) *MailService {
	return &MailService{
		API_URL_MAILHOG: url,
	}
}

func (s *MailService) SendSMPTMessage(from, to, subject, message string) error {

	toArr := []string{to}
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, message))

	return smtp.SendMail(s.API_URL_MAILHOG, nil, from, toArr, msg)
}
