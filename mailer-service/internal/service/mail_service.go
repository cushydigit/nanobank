package service

type MailService struct {
}

func NewMailService() *MailService {
	return &MailService{}
}

func (s *MailService) SendSMPTMessage(from, to, subject, message string) error {
	return nil
}
