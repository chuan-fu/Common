package retryx

import (
	"context"
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	cli := NewClient(WithHeaders(map[string]string{
		"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTY2NDEyMTQsImlhdCI6MTY1NjU1NDgxNCwidXNlciI6eyJ0aW1lc3RhbXAiOiIxNjU2NTU0ODE0ODUyNzc0MjQ1IiwidXNlcl9pZCI6NTU4Njg0OH19.pLHQmRpNelktD-3Ch_Z-WJcvXjdY3yXDOwZJguHN_GI",
	}), WithCheckResp(NewCheckResp("Code", "Msg", 0)))

	url := "http://goapi.dlab.cn/dpm/bms/station_letter/list"

	var r Tools14
	data, err := cli.PostResult(context.TODO(), url, &Tools13{
		PageIndex: 1,
		PageSize:  10,
	}, &r)
	fmt.Println(err)
	if err == nil {
		fmt.Println(string(data))
		fmt.Printf("%+v", r)
	}
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
