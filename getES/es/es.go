package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
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

func getPullFromES() {
	pullObj := elastic.NewBoolQuery()
	pullObj.Must(elastic.NewMatchQuery("cur_pull_stream_num_online"))
	pullObj.Filter(elastic.NewRangeQuery("ts").Gt(1598959320).Lte(1598959380))
	res, err := client.Search("zrtclive_info_202009").Query(pullObj).Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

}
