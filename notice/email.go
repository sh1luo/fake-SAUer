package notice

import (
	"gopkg.in/gomail.v2"
	"strconv"
)

type EmailNotifier struct {
	Account string `json:"account"`
	Token   string `json:"token"`
	Host    string `json:"host"`
	Port    string `json:"port"`
}

func NewEmailNotifier(notifier interface{}) *EmailNotifier {
	e := notifier.(*EmailNotifier)
	return e
}

func (n *EmailNotifier) Notice(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(n.Account, "Punch Message"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	
	p, err := strconv.Atoi(n.Port)
	if err != nil {
		return err
	}
	dialer := gomail.NewDialer(n.Host, p, n.Account, n.Token)
	return dialer.DialAndSend(m)
}
