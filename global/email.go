package global

import (
	"gopkg.in/gomail.v2"
	"strconv"
)

type Email struct {
	*smtpEmail
}

func NewEmail(info *smtpEmail) *Email {
	return &Email{smtpEmail: info}
}

func (e *Email) SendMail(to string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.Account, "I'm a Robot"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	p, err := strconv.Atoi(e.Port)
	if err != nil {
		return err
	}
	dialer := gomail.NewDialer(e.Host, p, e.Account, e.Token)
	return dialer.DialAndSend(m)
}
