package robot

import (
	"context"
	"fmt"

	"github.com/chuan-fu/Common/util"
	"github.com/pkg/errors"
)

// post post请求
func post(ctx context.Context, key string, data *Message) error {
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", key)

	var res RobotResponse
	if _, err := util.Cli(nil).PostResult(ctx, url, data, &res); err != nil {
		return errors.Wrap(err, "请求出错")
	}
	if res.ErrorCode != 0 {
		return errors.Wrap(errors.New(res.ErrorMessage), "请求出错")
	}
	return nil
}
