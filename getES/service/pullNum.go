package service

import (
	"fmt"
	"opTools/getES/modle"
	"time"
)

var EsChan chan modle.EsTask
var PullChan = make(chan *modle.Zrtclive, 1000)

func MakeTask() {
	for {

		EsChan = make(chan modle.EsTask, 500)
		now := time.Now()
		t := modle.EsTask{
			Index:     fmt.Sprintf("zrtclive_info_%d", now.Format("200601")),
			StartTime: time.Now().Unix() - 60,
			EndTime:   time.Now().Unix(),
		}
		EsChan <- t
		time.Sleep(time.Second * 60)
	}

}
