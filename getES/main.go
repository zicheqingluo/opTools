package main

import (
	"fmt"
	"opTools/getES/es"
	"opTools/getES/service"
	"time"
)

func main() {
	es.ESInit()
	fmt.Println("启动maketask")
	go service.MakeTask()
	time.Sleep(time.Second * 2)
	fmt.Println("启动 es.run")
	go es.Run()
	//	time.Sleep(time.Second * 200)
	fmt.Println("启动listenpullchan")
	service.ListenPullChan()

}
