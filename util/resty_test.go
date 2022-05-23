package util

import (
	"context"
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	SetGlobalResty(NewResty())
	cli := Cli(map[string]string{
		"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDk5MzEyMjEsImlhdCI6MTY0OTg0NDgyMSwidXNlciI6eyJ0aW1lc3RhbXAiOiIxNjQ5ODQ0ODIxMjI2MjMxODA1IiwidXNlcl9pZCI6NTU4Njg0OH19.sbcm4t318D2qywy7jAWlM_EoRrYyFi5nRRlIjU6N-3A",
	})
	url := "http://goapi.dlab.cn/dpm/bms/station_letter/list"

	var r Tools14
	f := NewCheckResp("Code", "Msg", 0)
	_ = f
	data, err := cli.PostCheckResult(context.TODO(), url, &Tools13{
		PageIndex: 1,
		PageSize:  10,
	}, &r, f)
	fmt.Println(err)
	fmt.Println(string(data))
	fmt.Printf("%+v", r)
}

type Tools13 struct {
	PageIndex int64 `json:"page_index"`
	PageSize  int64 `json:"page_size"`
}

type Tools14 struct {
	Flag bool   `json:"flag"`
	Msg  string `json:"msg"`
	Code int64  `json:"code"`
	Data Data   `json:"data"`
}

type Data struct {
	PageIndex int64   `json:"page_index"`
	PageSize  int64   `json:"page_size"`
	Total     int64   `json:"total"`
	Datas     []Datas `json:"datas"`
}

type Datas struct {
	ID        int64     `json:"id"`
	MsgType   int64     `json:"msg_type"`
	Message   []Message `json:"message"`
	IsRead    int64     `json:"is_read"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

type Message struct {
	Msgtype int64  `json:"msgtype"`
	Color   string `json:"color"`
	Content string `json:"content"`
}

func TestExampleCheck(t *testing.T) {
	fmt.Println(exampleCheckResp.Check(&exampleResp{Code: 0, Msg: "00"}))
	fmt.Println(exampleCheckResp.Check(&exampleResp{Code: 1, Msg: "11"}))
}
