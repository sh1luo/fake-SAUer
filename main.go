package main

import (
	"fake-SAUer/core"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

var f *core.Faker

func main() {
	var err error
	f, err = core.NewFaker(true)
	if err != nil {
		panic(err)
	}
	
	// 两次是防止有一次出问题，1和3是避免0点高峰
	c := cron.New()
	_, err = c.AddFunc("* * * * *", func() {
		fmt.Printf("********准备开始今日打卡,共需打卡%d人********\n", f.Cnt)
		start := time.Now()
		done := f.Do()
		log.Printf("总用时%s,共需打卡%d人,成功打卡%d人\n\n", time.Since(start), f.Cnt, done)
	})
	if err != nil {
		panic(err)
	}
	
	c.Start()
	
	StartHTTPServer()
}
