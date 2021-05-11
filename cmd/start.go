package cmd

import (
	"fake-SAUer/core"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func init() {
	f, err := core.NewFaker()
	if err != nil {
		panic(err)
	}

	// 创建一个定时任务，每天1和3点执行一次，两次是防止有一次出问题
	c := cron.New()
	_,err = c.AddFunc("0 1,3 * * *", func() {
		fmt.Printf("********初始化完成，共需打卡%d人********\n", f.Cnt)
		start := time.Now()
		// 执行打卡逻辑
		done := f.Do()
		log.Printf("总用时%s,共需打卡%d人，成功打卡%d人\n\n", time.Since(start), f.Cnt, done)
	})
	if err != nil {
		panic(err)
	}

	c.Start()
}
