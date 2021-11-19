package notice

import (
	"gopkg.in/gomail.v2"
)

type EmailNotifier struct {
	Account string `json:"account"`
	Token   string `json:"token"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
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

	dialer := gomail.NewDialer(n.Host, n.Port, n.Account, n.Token)
	return dialer.DialAndSend(m)
}
