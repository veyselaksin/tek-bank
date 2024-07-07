package gomailer

import (
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"github.com/veyselaksin/gomailer/pkg/mailer"
	"os"
)

type Content struct {
	Subject     string
	Body        string
	To          []string
	CC          []string
	BCC         []string
	Attachments []string
}

func SendMail(content Content) error {
	auth := mailer.Authentication{
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
	}
	sender := mailer.NewPlainAuth(&auth)

	if len(content.To) == 0 {
		return errors.New("To field is required")
	}

	message := mailer.NewMessage(content.Subject, content.Body)
	message.SetFrom(os.Getenv("SMTP_USERNAME"))
	message.SetTo(content.To)

	if len(content.Attachments) > 0 {
		for _, attachment := range content.Attachments {
			err := message.SetAttachFiles(attachment)
			if err != nil {
				log.Error("Error setting attachment: ", err)
				return err
			}
		}
	}

	err := sender.SendMail(message)
	if err != nil {
		log.Error("Error sending mail: ", err)
		return err
	}
	return nil
}
