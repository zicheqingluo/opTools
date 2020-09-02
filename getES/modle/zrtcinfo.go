package modle

import (
	"time"
)

type Host struct {
	Name string `json:"name"`
}

// Zrtclive es数据
type Zrtclive struct {
	Timestamp time.Time `json:"@timestamp"`
	Hostname  Host      `json:"host"`
	DateType  string    `json:"dataType"`
	Localip   string    `json:"local_ip"`
	Kcp       int64     `json:"cur_kcp_stream_num_online"`
	PullNum   int64     `json:"cur_pull_stream_num_online"`
	TS        int64     `json:"ts"`
}

type EsTask struct {
	Index     string
	StartTime int64
	EndTime   int64
}
