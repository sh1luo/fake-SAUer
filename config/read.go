package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	WithEmail WithEmail  `json:"withEmail"`
	StusInfos []*StuInfo `json:"stu_info"`
}

type WithEmail struct {
	On      bool   `json:"on"`
	Account string `json:"account"`
	Token   string `json:"token"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	IsSSL   bool   `json:"is_ssl"`
	From    string `json:"from"`
}

type StuInfo struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Province string `json:"province"`
	City     string `json:"city"`
	College  string `json:"college"`
	Account  string `json:"account"`
	Passwd   string `json:"passwd"`
}

func ReadConfig() *Config {
	var c Config

	f, err := os.Open("config.json")
	if err != nil {
		return nil
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		return nil
	}

	return &c
}
