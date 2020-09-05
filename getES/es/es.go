package es

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"opTools/getES/model"
	"opTools/getES/service"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
)

var client *elastic.Client
var host = "http://192.168.148.191:9601/"

func ESInit() {
	errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
	var err error
	client, err = elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetURL(host))
	if err != nil {
		panic(err)
	}
	info, code, err := client.Ping(host).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("ES returned with code %d and version %s\n", code, info.Version.Number)
	esVersion, err := client.ElasticsearchVersion(host)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s ES version %s\n", host, esVersion)

}

func getPullFromES(index string, st, et int64) (int64, error) {
	pullList := []*modle.Zrtclive{}

	pullObj := elastic.NewBoolQuery()
	pullObj.Filter(elastic.NewRangeQuery("ts").Gt(st).Lt(et))
	//pullObj.Filter(elastic.NewRangeQuery("ts").Gt(1599031879))
	res, err := client.Search(index).Query(pullObj).Size(500).Do(context.Background())
	if err != nil {
		fmt.Println("getpullfromes", err)
		return 0, err
	}
	fmt.Println("开始时间：", st, "结束时间：", et)
	totalHits := res.Hits.TotalHits.Value
	fmt.Println("总数：", res.Hits.TotalHits.Value)
	fmt.Println("hist:", len(res.Hits.Hits))
	for _, value := range res.Hits.Hits {

		var doc *modle.Zrtclive
		json.Unmarshal(value.Source, &doc)
		//fmt.Printf("long:%t", doc.TS)
		pullList = append(pullList, doc)
		//fmt.Println(doc.Localip)

	}
	time.Sleep(time.Second * 5)
	service.PullChan <- pullList
	return totalHits, nil

}

func Run() {

	for {

		select {
		case t := <-service.EsChan:
			fmt.Println("es.run:channel has key")
			totalHits, err := getPullFromES(t.Index, t.StartTime, t.EndTime)
			if err != nil {
				fmt.Printf("获取es数据失败：%v", err)
			}
			if totalHits > 500 {

				fmt.Println("数据大于500")
			}
			//default:
			//	time.Sleep(time.Second * 10)

		}
	}
}
