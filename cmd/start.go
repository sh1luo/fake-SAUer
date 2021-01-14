package cmd

import (
	"fake-SAUer/internal/faker"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/signal"
	"time"
)

var punchCount uint = 1

func init() {
	c := cron.New()
	f, err := faker.NewFaker()
	if err != nil {
		panic(err)
	}
	fmt.Printf("********初始化完成，共需打卡%d人********\n", f.Cnt)

	// 设置打卡两次，防止系统故障
	_, err = c.AddFunc("* * * * *", func() {
		start := time.Now()
		fmt.Printf("%s开始第%d次打卡:\n", start.String()[:19], punchCount)

		// TODO:目前先每次载入所有数据，后续改为文件监听
		// f, err = faker.NewFaker()

		// 执行打卡逻辑
		f.Do()

		fmt.Printf("第%d次打卡完毕，总用时%ss\n\n", punchCount, time.Since(start))
		punchCount++
	})

	if err != nil {
		log.Fatalf("Add job err: %v", err)
		return
	}

	//启动定时任务
	c.Start()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
