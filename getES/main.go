package main

import (
	"github.com/google/gops/agent"
	//"go.uber.org/zap"
	"opTools/getES/es"
	"opTools/getES/service"
	zlog "opTools/getES/zaplog"
	"time"
)

func main() {
	zlog.InitLog("./logs/info.log", "./logs/error.log", "info")
	if err := agent.Listen(agent.Options{Addr: "0.0.0.0:8000"}); err != nil {

		zlog.Warn("gops初始化失败:%s", err)
	}
	zlog.Info("gops运行")
	//logger.AtomicLevel.SetLevel(zap.InfoLevel)
	//logger.Logger.Debug("启动。。")
	es.ESInit()
	zlog.Debug("ES init ...")
	go service.MakeTask()
	time.Sleep(time.Second * 2)
	zlog.Debug("make task ...")
	go es.Run()
	//	time.Sleep(time.Second * 200)
	zlog.Debug("es.Run listenPullChan ...")
	service.ListenPullChan()

}
