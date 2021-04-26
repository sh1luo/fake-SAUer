package global

import (
	"encoding/json"
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
)

type Config struct {
	E         *smtp      `json:"withEmail"`
	StusInfos []*StuInfo `json:"stu_info"`
}

type StuInfo struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Province string `json:"province"`
	City     string `json:"city"`
	College  string `json:"college"`
	Account  string `json:"account"`
	Passwd   string `json:"passwd"`
	Email    string `json:"email"`

	UUID string `json:"uuid"`
}

type smtp struct {
	Enabled bool   `json:"enabled"`
	Account string `json:"account"`
	Token   string `json:"token"`
	Host    string `json:"host"`
	Port    string `json:"port"`
}

func (e *smtp) SendMail(to string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.Account, "Punch MSG"))
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

func ReadConfig() (err error) {
	c := &Config{
		E:         new(smtp),
		StusInfos: make([]*StuInfo, 1),
	}

	f, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		return err
	}
	G_CONF = c
	return nil
}
