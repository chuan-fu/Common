package pushplus

import (
	"context"
	"fmt"

	"github.com/chuan-fu/Common/util"
	"github.com/pkg/errors"
)

// post post请求
func post(ctx context.Context, data *Message) error {
	var res Result
	if _, err := util.Cli(nil).PostResult(ctx, uri, data, &res); err != nil {
		return errors.Wrap(err, "请求出错")
	}
	if res.Code != 200 {
		return errors.Wrap(errors.New(fmt.Sprintf("%v", res.Msg)), "请求出错")
	}
	return nil
}
