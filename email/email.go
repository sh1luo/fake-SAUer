package email

import (
	"fake-SAUer/config"
	"gopkg.in/gomail.v2"
	"strconv"
)

type Email struct {
	*config.SMTPEmail
}

func NewEmail(info *config.SMTPEmail) *Email {
	return &Email{SMTPEmail: info}
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

func SendMail(to string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress("3450047248@qq.com", "高一宁女士您好"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return gomail.NewDialer("smtp.qq.com", 465, "3450047248@qq.com", "iyjgokhgzybudajj").DialAndSend(m)
}