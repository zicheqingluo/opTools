package es

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"log"
	"opTools/getES/model"
	"opTools/getES/service"
	zlog "opTools/getES/zaplog"
	//"os"
)

var client *elastic.Client
var host = "http://192.168.148.191:9601/"

func ESInit() {
	errorlog := log.New(zlog.WarnWriter, "APP", log.LstdFlags)
	//errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
	var err error
	client, err = elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetURL(host))
	if err != nil {
		zlog.Error("创建es链接错误:%v", err)
	}
	info, code, err := client.Ping(host).Do(context.Background())
	if err != nil {
		zlog.Error("链接es失败：%s", err)
	}
	zlog.Info("ES %s returned with code %d and version %s\n", host, code, info.Version.Number)
	/*
		esVersion, err := client.ElasticsearchVersion(host)
		if err != nil {
			panic(err)
		}
	*/

}

func getPullFromES(index string, st, et int64) (int64, error) {
	pullList := []*modle.Zrtclive{}

	pullObj := elastic.NewBoolQuery()
	pullObj.Filter(elastic.NewRangeQuery("ts").Gt(st).Lt(et))
	//pullObj.Filter(elastic.NewRangeQuery("ts").Gt(1599031879))
	zlog.Info("开始时间：%d- 结束时间：%d", st, et)
	res, err := client.Search(index).Query(pullObj).Size(500).Do(context.Background())
	if err != nil {
		zlog.Error("获取%d-%d的数据失败:%s", st, et, err)
		return 0, err
	}
	totalHits := res.Hits.TotalHits.Value
	zlog.Info("总数：%d", res.Hits.TotalHits.Value)
	zlog.Debug("hist:%d", len(res.Hits.Hits))
	for _, value := range res.Hits.Hits {

		var doc *modle.Zrtclive
		json.Unmarshal(value.Source, &doc)
		//fmt.Printf("long:%t", doc.TS)
		pullList = append(pullList, doc)
		//fmt.Println(doc.Localip)

	}
	//time.Sleep(time.Second * 5)
	service.PullChan <- pullList
	return totalHits, nil

}

func Run() {

	for {

		select {
		case t := <-service.EsChan:
			zlog.Info("es.run:channel has key")
			totalHits, err := getPullFromES(t.Index, t.StartTime, t.EndTime)
			if err != nil {
				zlog.Error("获取es数据失败：%v", err)
			}
			if totalHits > 500 {

				zlog.Warn("数据大于500")
			}
			//default:
			//	time.Sleep(time.Second * 10)

		}
	}
}
