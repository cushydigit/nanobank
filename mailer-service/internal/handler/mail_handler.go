package handler

import (
	"net/http"

	"github.com/cushydigit/nanobank/mailer-service/internal/service"
	"github.com/cushydigit/nanobank/shared/helpers"
)

type MailHandler struct {
	service *service.MailService
}

func NewMailHandler(s *service.MailService) *MailHandler {
	return &MailHandler{service: s}
}

func (h *MailHandler) SendMail(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, http.StatusAccepted, nil)
}
