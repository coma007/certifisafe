package auth

import (
	"bytes"
	"github.com/mailgun/mailgun-go"
	"net/smtp"
)

type MailService interface {
	SendMail(to []string, body bytes.Buffer) error
}

type DefaultMailService struct {
}

func NewDefaultMailService() *DefaultMailService {
	return &DefaultMailService{}
}

func (service *DefaultMailService) SendMail1(to []string, body bytes.Buffer) error {
	from := "ftn.project.usertest@gmail.com"
	password := "zmiwmhfweojejlqy"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	go smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	return nil
}

func (service *DefaultMailService) SendMail(to []string, body bytes.Buffer) error {
	mg := mailgun.NewMailgun(MailgunApiDomain, MailgunApiKey)
	mg.SetAPIBase(MailgunApiBase)
	m := mg.NewMessage(
		"Certifisafe <certifisafe@mailgun.org>",
		"Hello",
		"",
		to[0],
	)

	m.SetHtml(body.String())
	go mg.Send(m)
	return nil
}
