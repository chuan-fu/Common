package robot

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

// post post请求
func post(ctx context.Context, key string, data Message) error {
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", key)

	var res RobotResponse

	client := resty.New().SetTimeout(10 * time.Second)
	resp, err := client.R().
		SetHeader("content-type", "application/json;charset=UTF-8").
		SetContext(ctx).
		SetBody(data).
		SetResult(&res).
		Post(url)
	if err != nil {
		return errors.Wrap(err, "请求出错")
	}

	if resp.StatusCode() == 200 {
		if res.ErrorCode != 0 {
			return errors.Wrap(errors.New(res.ErrorMessage), "请求出错")
		}
		return nil
	}
	return errors.Wrap(errors.New(fmt.Sprintf("连接企微服务报错 status:%v", resp.StatusCode())), "请求出错")
}
