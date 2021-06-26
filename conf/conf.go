package conf

import (
	"encoding/json"
	"os"
)

var (
	G_Conf *Config
)

type Config struct {
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
	To       string `json:"to"`
	
	Uuid     string
}

func ReadConfig() error {
	c := &Config{
		StusInfos: make([]*StuInfo, 0),
	}
	
	f, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer f.Close()
	
	if err = json.NewDecoder(f).Decode(&c); err != nil {
		return err
	}
	
	G_Conf = c
	return nil
}
