package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"opTools/getZrtc/model"
)

func PushFalcon(data []model.FalconItem) error {
	fmt.Println("PushFalcon data:", data)
	ret, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return err
	}

	body := bytes.NewBuffer([]byte(ret))
	res, err := http.Post("http://127.0.0.1:1998/v1/push", "application/json;charset=utf-8", body)
	defer res.Body.Close()

	if err != nil {
		log.Printf("push /v1/push fail: %v", err)
		return err
	}
	return nil
}

func UpdateOdin(itemList []model.FalconItem) error {

	err := PushFalcon(itemList)
	if err != nil {
		log.Println("uploadOdin err:", err)
		return err
	}
	return nil
}
