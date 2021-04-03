package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	WithEmail *SMTPEmail `json:"withEmail"`
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

type SMTPEmail struct {
	Account string `json:"account"`
	Token   string `json:"token"`
	Host    string `json:"host"`
	Port    string `json:"port"`
}

func ReadConfig() (*Config, error) {
	c := &Config{
		WithEmail: new(SMTPEmail),
		StusInfos: make([]*StuInfo, 0),
	}

	f, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
