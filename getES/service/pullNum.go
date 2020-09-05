package service

import (
	odinm "common/model"
	odins "common/service"
	"fmt"
	"opTools/getES/model"
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
			StartTime: time.Now().Unix() - 60,
			EndTime:   time.Now().Unix(),
		}
		EsChan <- t
		fmt.Printf("maketask:put key to chan %v", EsChan)
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
		case strings.HasPrefix(value.Localip, "10.2"):
			bdsz[value.Localip] = value.PullNum
		case strings.HasPrefix(value.Localip, "10.3"):
			bdwh[value.Localip] = value.PullNum
		case strings.HasPrefix(value.Localip, "10.0"):
			//fmt.Println("广州ip：", value.Localip)
			bdgz[value.Localip] = value.PullNum

		default:
			fmt.Printf("%s 未找到所属机房\n", value.Localip)
		}

	}
	for _, num := range bdgz {
		bdgzSum += num
	}
	for _, num := range bdsz {
		bdszSum += num
	}
	for _, num := range bdbj {
		bdbjSum += num
	}
	for _, num := range bdwh {
		bdwhSum += num
	}
	fmt.Println("bdbj: ", bdbj)
	fmt.Println("bdsz: ", bdsz)
	fmt.Println("bdgz: ", bdgz)
	fmt.Println("bdwh: ", bdwh)
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
	err := odins.UpdateOdin(itemList)
	if err != nil {
		fmt.Println("上传odin失败：", err)
	}
}

func ListenPullChan() {
	for {
		select {
		case d := <-PullChan:
			fmt.Println(time.Now())
			makeData(d)

		}
	}
}
