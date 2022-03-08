package util

import (
	"net/http"
	"time"

	"gopkg.in/resty.v1"
)

var client *resty.Client

func init() {
	client = resty.New()
	client.SetTimeout(5 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(time.Second)
	client.AddRetryCondition(func(resp *resty.Response) (isRetry bool, err error) {
		if resp == nil || resp.StatusCode() != http.StatusOK {
			return true, nil
		}
		return false, nil
	})
}

func Client() *resty.Client {
	return client
}
