package service

import (
	"fmt"
	"opTools/getES/modle"
	"time"
)

var EsChan chan modle.EsTask
var PullChan = make(chan *model.Zrtclive, 1000)

func MakeTask() {
	for {

		EsChan = make(chan modle.EsTask, 500)
		now := time.now()
		t := modle.EsTask{
			Index:     fmt.Sprintf("zrtclive_info_%d", now.format("200601")),
			StartTime: time.Now().Unix() - 60,
			EndTime:   time.Now().Unix(),
		}
		EsChan <- t
		time.Sleep(time.Second * 60)
	}

}
