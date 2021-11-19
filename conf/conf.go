package conf

import (
	"encoding/json"
	"os"
)

var (
	GlobalConfig *Config
)

type Config struct {
	NotifierInfo NotifierInfo `json:"notifier_info"`
	StusInfo     []*StuInfo   `json:"stu_info"`
}

type NotifierInfo struct {
	Method string    `json:"method"`
	Email  SMTPEmail `json:"email"`
}

type StuInfo struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Province string `json:"province"`
	City     string `json:"city"`
	College  string `json:"college"`
	Account  string `json:"account"`
	Passwd   string `json:"passwd"`
	To       string `json:"to"`

	Uuid string `json:"uuid"`
}

type SMTPEmail struct {
	Account string `json:"account"`
	Token   string `json:"token"`
	Host    string `json:"host"`
	Port    string `json:"port"`
}

func ReadConfig() error {
	GlobalConfig = &Config{
		StusInfo: make([]*StuInfo, 5),
	}

	f, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)
	
	err = json.NewDecoder(f).Decode(GlobalConfig)
	if err != nil {
		return err
	}
	
	return nil
}

func FilterInvalidStudents() {

}