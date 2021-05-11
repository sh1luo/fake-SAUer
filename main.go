package main

import (
	_ "fake-SAUer/cmd"
)

func main() {
	// 由于程序内启动了定时任务，所以要阻止应用程序退出
	// 改用crontab就可以去掉了
	select {}
}
