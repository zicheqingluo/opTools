package db

import (
	_ "github.com/go-sql-driver/mysql"
	"opTools/getZrtc/model"
)

func GetZrtcNums(createTime int64) (onlineList []model.OnlineUser, err error) {
	sqlstr := "select pull_stream_num,push_stream_num,create_time from tblCurrentStreams where create_time > ?"
	err = DB.Select(&onlineList, sqlstr, createTime)
	return
}

func GetFirst() (onlineList []model.OnlineUser, err error) {
	sqlstr := "select pull_stream_num,push_stream_num,create_time from tblCurrentStreams order by create_time desc limit 1;"
	err = DB.Select(&onlineList, sqlstr)

	return
}
