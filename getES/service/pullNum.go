package service

import (
	odinm "common/model"
	odins "common/service"
	zlog "common/zaplog"
	"fmt"
	"opTools/getES/model"
	"opTools/getES/util"
	"strings"
	"time"
)

var EsChan chan modle.EsTask
var PullChan = make(chan []*modle.Zrtclive, 10)

func MakeTask() {
	EsChan = make(chan modle.EsTask, 50)
	for {

		now := time.Now()
		t := modle.EsTask{
			Index:     fmt.Sprintf("zrtclive_info_%v", now.Format("200601")),
			StartTime: now.Unix() - 60,
			EndTime:   now.Unix(),
		}
		EsChan <- t
		zlog.Info("生成任务:%s", now.Format("200601"))
		time.Sleep(time.Second * 60)
	}

}

func makeData(data []*modle.Zrtclive) {
	bdgz := make(map[string]int64, 100)
	bdsz := make(map[string]int64, 100)
	bdbj := make(map[string]int64, 100)
	bdwh := make(map[string]int64, 100)
	allidc := map[string]*map[string]int64{
		"bdbj": &bdbj,
		"bdsz": &bdsz,
		"bdgz": &bdgz,
		"bdwh": &bdwh,
	}
	for _, value := range data {

		switch {
		case strings.HasPrefix(value.Localip, "192.168"):
			/*
				v, ok := idc["bdbj"]
				if ok {
					v.IdcData[value.Localip] = value.PullNum
					idc["bdbj"] = v

				} else {
					v := modle.IdcNum{
						IdcData[value.Localip]: value.PullNum,
						Sum:                    0,
						Metric:                 "bdbjpullnum",
					}
					idc["bdbj"] = v
				}
			*/
			bdbj[value.Localip] = value.PullNum
		case strings.HasPrefix(value.Localip, "10.2"):
			bdsz[value.Localip] = value.PullNum
		case strings.HasPrefix(value.Localip, "10.3"):
			bdwh[value.Localip] = value.PullNum
		case strings.HasPrefix(value.Localip, "10.0"):
			//fmt.Println("广州ip：", value.Localip)
			bdgz[value.Localip] = value.PullNum

		default:
			zlog.Warn("%s 未找到所属机房", value.Localip)
		}

	}
	zlog.Info("bdbj: %v", bdbj)
	zlog.Info("bdsz: %v", bdsz)
	zlog.Info("bdgz: %v", bdgz)
	zlog.Info("bdwh: %v", bdwh)
	itemList := []odinm.FalconItem{}
	for k, v := range allidc {
		var sum int64
		for _, num := range *v {
			sum += num
		}
		odin := odinm.FalconItem{
			Endpoint:    "rdqa-rd-test666.bjdd.zybang.com",
			Metric:      k + "pullnum",
			Timestamp:   time.Now().Unix(),
			Step:        60,
			Value:       sum,
			CounterType: "GAUGE",
			Tags:        "",
		}
		itemList = append(itemList, odin)
		zlog.Info("%s 人数:%d", k, sum)

	}
	upOdin := func(itemList []odinm.FalconItem) func() error {
		return func() error {
			err := odins.UpdateOdin(itemList)
			return err
		}
	}
	err := util.Retry(3, 5*time.Second, upOdin(itemList))
	if err != nil {
		zlog.Error("上传odin失败：%s", err)
	}
	zlog.Info("上传至odin")
}

func ListenPullChan() {
	for {

		select {
		case d := <-PullChan:
			zlog.Info("处理数据")
			makeData(d)

		}
	}
}
