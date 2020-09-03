package main

import (
	"opTools/getES/es"
	"opTools/getES/service"
)

func main() {
	es.ESInit()
	go service.MakeTask()
	go es.Run()
	service.ListenPullChan()

}
