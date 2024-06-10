package main

import (
	"fmt"
	"log"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

func (app *application) sendConfirmationEmail(toEmail, emailToken string) error {
	server := mail.NewSMTPClient()

	server.Host = app.mailConfig.host
	server.Username = app.mailConfig.username
	server.Password = app.mailConfig.password
	server.Port = app.mailConfig.port

	server.KeepAlive = false
	server.Encryption = mail.EncryptionTLS
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	email := mail.NewMSG()
	email.SetFrom(app.mailConfig.from).
		AddTo(toEmail).
		SetSubject("Confirm you email please").
		SetBody(mail.TextHTML, fmt.Sprintf("Click the following link to confirm your email please <a href=\"http://localhost:8080/confirm-email?token=%s\">Confirm my email address</a>", emailToken))

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}
	return nil
}
