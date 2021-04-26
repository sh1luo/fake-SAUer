package cmd

import (
	"fake-SAUer/core"
	"fmt"
	"log"
	"time"
)

func init() {
	f, err := core.NewFaker()
	if err != nil {
		panic(err)
	}
	fmt.Printf("********初始化完成，共需打卡%d人********\n", f.Cnt)

	start := time.Now()
	done := f.Do()	// 执行打卡逻辑
	log.Printf("总用时%s,共需打卡%d人，成功打卡%d人\n", time.Since(start),f.Cnt,done)
	if err != nil {
		log.Fatalf("Add job err: %v", err)
		return
	}
}
