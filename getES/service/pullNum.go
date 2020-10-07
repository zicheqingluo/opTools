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
	var bdgzSum int64
	var bdszSum int64
	var bdbjSum int64
	var bdwhSum int64
	for _, value := range data {

		switch {
		case strings.HasPrefix(value.Localip, "192.168"):
			bdbj[value.Localip] = value.PullNum
			bdbjSum += value.PullNum
		case strings.HasPrefix(value.Localip, "10.2"):
			bdsz[value.Localip] = value.PullNum
			bdszSum += value.PullNum
		case strings.HasPrefix(value.Localip, "10.3"):
			bdwh[value.Localip] = value.PullNum
			bdwhSum += value.PullNum
		case strings.HasPrefix(value.Localip, "10.0"):
			//fmt.Println("广州ip：", value.Localip)
			bdgz[value.Localip] = value.PullNum
			bdgzSum += value.PullNum

		default:
			zlog.Warn("%s 未找到所属机房", value.Localip)
		}

	}
	zlog.Info("bdbj: %v", bdbj)
	zlog.Info("bdsz: %v", bdsz)
	zlog.Info("bdgz: %v", bdgz)
	zlog.Info("bdwh: %v", bdwh)
	odinbdbj := odinm.FalconItem{
		Endpoint:    "rdqa-rd-test666.bjdd.zybang.com",
		Metric:      "bdbjpullnum",
		Timestamp:   time.Now().Unix(),
		Step:        60,
		Value:       bdbjSum,
		CounterType: "GAUGE",
		Tags:        "",
	}
	odinbdsz := odinm.FalconItem{
		Endpoint:    "rdqa-rd-test666.bjdd.zybang.com",
		Metric:      "bdszpullnum",
		Timestamp:   time.Now().Unix(),
		Step:        60,
		Value:       bdszSum,
		CounterType: "GAUGE",
		Tags:        "",
	}
	odinbdwh := odinm.FalconItem{
		Endpoint:    "rdqa-rd-test666.bjdd.zybang.com",
		Metric:      "bdwhpullnum",
		Timestamp:   time.Now().Unix(),
		Step:        60,
		Value:       bdwhSum,
		CounterType: "GAUGE",
		Tags:        "",
	}
	odinbdgz := odinm.FalconItem{
		Endpoint:    "rdqa-rd-test666.bjdd.zybang.com",
		Metric:      "bdgzpullnum",
		Timestamp:   time.Now().Unix(),
		Step:        60,
		Value:       bdgzSum,
		CounterType: "GAUGE",
		Tags:        "",
	}
	itemList := []odinm.FalconItem{}
	itemList = append(itemList, odinbdwh, odinbdgz, odinbdbj, odinbdsz)
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
