package pushplus

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

// post post请求
func post(ctx context.Context, data *Message) error {
	var res Result

	client := resty.New().SetTimeout(5 * time.Second)
	resp, err := client.R().
		SetHeader("content-type", "application/json;charset=UTF-8").
		SetContext(ctx).
		SetBody(data).
		SetResult(&res).
		Post(uri)
	if err != nil {
		return errors.Wrap(err, "请求出错")
	}

	if resp.StatusCode() == http.StatusOK {
		if res.Code != 200 {
			return errors.Wrap(errors.New(fmt.Sprintf("%v", res.Msg)), "请求出错")
		}
		return nil
	}
	return errors.Wrap(errors.New(fmt.Sprintf("连接pushplus服务报错 status:%v", resp.StatusCode())), "请求出错")
}
