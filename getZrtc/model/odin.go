package model

type FalconItem struct {
	Endpoint    string      `json:"endpoint"`
	Metric      string      `json:"metric"`
	Value       interface{} `json:"value"`
	CounterType string      `json:"counterType"`
	Step        int64       `json:"step"`
	Timestamp   int64       `json:"timestamp"`
	Tags        string      `json:"tags"`
}
