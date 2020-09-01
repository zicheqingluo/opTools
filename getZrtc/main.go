package main

import (
	"opTools/getZrtc/db"
	"opTools/getZrtc/service"
)

func main() {
	connInfo := "smt_op:smt_op@tcp(192.168.148.191:3306)/homework_smt?parseTime=true"
	err := db.Init(connInfo)
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()
	service.Run()

}
