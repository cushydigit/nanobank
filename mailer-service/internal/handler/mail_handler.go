package handler

import (
	"net/http"

	"github.com/cushydigit/nanobank/mailer-service/internal/service"
	myerrors "github.com/cushydigit/nanobank/shared/errors"
	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/types"
)

type MailHandler struct {
	service *service.MailService
}

func NewMailHandler(s *service.MailService) *MailHandler {
	return &MailHandler{
		service: s,
	}
}

func (h *MailHandler) SendMail(w http.ResponseWriter, r *http.Request) {
	req, ok := r.Context().Value(types.SendMailReqKey).(types.SendMailReqBody)
	if !ok {
		helpers.ErrorJSON(w, myerrors.ErrContextValueNotFoundInRequest, http.StatusInternalServerError)
		return
	}
	if err := h.service.SendSMPTMessage(req.From, req.To, req.Subject, req.Message); err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	payload := types.Response{
		Error:   false,
		Message: "success",
		Data:    req,
	}

	helpers.WriteJSON(w, http.StatusAccepted, payload)
}
