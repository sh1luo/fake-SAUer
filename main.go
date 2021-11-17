package main

import (
	"fake-SAUer/conf"
	"fake-SAUer/core"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func main() {
	err := conf.ReadConfig()
	if err != nil {
		panic(err)
	}
	
	f, err := core.NewFaker(true)
	if err != nil {
		panic(err)
	}
	log.Printf("Loading finished, totally %d students...\n", f.Cnt)

	// 两次是防止有一次出问题，1和3是避免0点高峰
	c := cron.New()
	_, err = c.AddFunc("0 1,3 * * *", func() {
		log.Printf("start...")
		start := time.Now()
		done := f.Do()
		log.Printf("Completed today, %s sec. %d in needing to sign-in,successfully %d\n", time.Since(start), f.Cnt, done)
	})
	if err != nil {
		panic(err)
	}
	
	c.Start()
	select {
	
	}
	core.StartHTTPServer()
}
