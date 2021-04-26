package global

import (
	"encoding/json"
	"os"
)

type Config struct {
	WithEmail *smtpEmail `json:"withEmail"`
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

type smtpEmail struct {
	Account string `json:"account"`
	Token   string `json:"token"`
	Host    string `json:"host"`
	Port    string `json:"port"`
}

func ReadConfig() (err error) {
	c := &Config{
		WithEmail: new(smtpEmail),
		StusInfos: make([]*StuInfo, 1),
	}

	f, err := os.Open("global.json")
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
