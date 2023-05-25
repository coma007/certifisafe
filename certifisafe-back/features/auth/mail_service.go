package auth

import (
	"bytes"
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

func (service *DefaultMailService) SendMail(to []string, body bytes.Buffer) error {
	from := "ftn.project.usertest@gmail.com"
	password := "zmiwmhfweojejlqy"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	go smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	return nil
}

//func (service *DefaultMailService) SendMail(to []string, body bytes.Buffer) error {
//	client := mailjet.NewMailjetClient("a4a2e88fbc90aec7ded74673efb3962e", "53f100a8e24c374d5fd9b453d7f88f73")
//	messagesInfo := []mailjet.InfoMessagesV31{
//		{
//			From: &mailjet.RecipientV31{
//				Email: "ftn.project.usertest@gmail.com",
//				Name:  "Mailjet Pilot",
//			},
//			To: &mailjet.RecipientsV31{
//				mailjet.RecipientV31{
//					Email: to[0],
//					Name:  to[0],
//				},
//			},
//			Subject:  "Your email flight plan!",
//			TextPart: "Dear passenger 1, welcome to Mailjet! May the delivery force be with you!",
//			HTMLPart: "<h3>Dear passenger 1, welcome to <a href=\"https://www.mailjet.com/\">Mailjet</a>!</h3><br />May the delivery force be with you!",
//		},
//	}
//	messages := mailjet.MessagesV31{Info: messagesInfo}
//	_, err := client.SendMailV31(&messages)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return err
//}
