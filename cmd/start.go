package cmd

import (
	"fake-SAUer/internal/faker"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

var punchCount uint = 1

func init() {
	c := cron.New()
	f, err := faker.NewFaker()
	if err != nil {
		panic(err)
	}
	fmt.Println(f.Cf.StusInfos[0].UUID)
	fmt.Printf("********初始化完成，共需打卡%d人********\n", f.Cnt)
	rubbish := ""

	// 设置打卡两次，防止系统故障
	//_, err = c.AddFunc("* * * * *", func() {
		start := time.Now()
		fmt.Printf("%s开始第%d次打卡:\n", start.String()[:19], punchCount)

		// TODO:目前先每次载入所有数据，后续改为文件监听
		// f, err = faker.NewFaker()

		// 执行打卡逻辑
		ok := f.Do()

		fmt.Printf("第%d次打卡完毕，总用时%s\n\n", punchCount, time.Since(start))
		punchCount++
		if !ok {
			rubbish += "又"
			// _ = email.SendMail("255327317@qq.com", "高一宁小朋友早上好呀", "今天打卡"+rubbish+"出问题了呢，赶快手动打一下卡，然后告诉你蓝朋友吧:)\n\n")
		}
	//})

	if err != nil {
		log.Fatalf("Add job err: %v", err)
		return
	}

	//启动定时任务
	c.Start()
}
