package util

import (
	"encoding/json"
	"fmt"
	"log"
)

type RespMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewRespMsg(code int, msg string, data interface{}) *RespMsg {
	return &RespMsg{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func (resp *RespMsg) JSONBytes() []byte {
	r, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	return r
}

func (resp *RespMsg) JSONString() string {
	r, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	return string(r)
}

func GenSimpleRespStream(code int, msg string) []byte {
	return []byte(fmt.Sprintf(`{"code":&d, "msg":"%s"}`, code, msg))
}

func GenSimpleRespString(code int, msg string) string {
	return fmt.Sprintf(`{"code":&d, "msg":"%s"}`, code, msg)
}
