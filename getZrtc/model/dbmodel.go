package model

type OnlineUser struct {
	PullNum    int64 `db:"pull_stream_num"`
	PushNum    int64 `db:"push_stream_num"`
	CreateTime int64 `db:"create_time"`
}
