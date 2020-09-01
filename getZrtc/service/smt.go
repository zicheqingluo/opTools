package service

import (
	"fmt"
	"opTools/getZrtc/db"
	"opTools/getZrtc/model"
	"time"
)

var ct int64
var infoList []model.OnlineUser
var err error

func Run() {
	for {
		if ct == 0 {
			fmt.Println("INFO:0ct=", ct)
			infoList, err = db.GetFirst()
			if err != nil {
				panic(err)
			}
			updateOdin()
		} else {
			fmt.Println("INFO:ct=", ct)
			infoList, err = db.GetZrtcNums(ct)
			if err != nil {
				fmt.Println("error:", err)
			}
			updateOdin()
		}

	}
}

func updateOdin() {
	if infoList != nil {
		fmt.Println("INFO:infoList=", infoList)
		var itemList = []model.FalconItem{}
		var lastID int
		for i, v := range infoList {
			pull := model.FalconItem{
				Endpoint:    "rdqa-rd-test666.bjdd.zybang.com",
				Metric:      "pullnum",
				Timestamp:   v.CreateTime,
				Step:        60,
				Value:       v.PullNum,
				CounterType: "GAUGE",
				Tags:        "",
			}
			push := model.FalconItem{
				Endpoint:    "rdqa-rd-test666.bjdd.zybang.com",
				Metric:      "pushnum",
				Timestamp:   v.CreateTime,
				Step:        60,
				Value:       v.PullNum,
				CounterType: "GAUGE",
				Tags:        "",
			}
			itemList = append(itemList, pull, push)

			lastID = i
			fmt.Println("INFO:lastID=", lastID)
		}
		fmt.Println("itemlist:", itemList)
		err := UpdateOdin(itemList)
		if err != nil {
			return
		}

		ct = itemList[lastID].Timestamp
		fmt.Println("INFO:new ct=", ct)
		time.Sleep(time.Second * 60)

	} else {
		time.Sleep(time.Second * 60)

	}
}
