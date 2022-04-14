package util

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

var globalClient *resty.Client

const (
	defaultTimeOut       = 5 * time.Second
	defaultRetryCount    = 3
	defaultRetryWaitTime = time.Second
)

func init() {
	globalClient = resty.New()
	globalClient.SetTimeout(defaultTimeOut)
	globalClient.SetRetryCount(defaultRetryCount)
	globalClient.SetRetryWaitTime(defaultRetryWaitTime)
	globalClient.AddRetryCondition(func(resp *resty.Response) (isRetry bool, err error) {
		if resp == nil || resp.StatusCode() != http.StatusOK {
			return true, nil
		}
		return false, nil
	})
}

// 新建空client
func NewResty() *resty.Client {
	return resty.New()
}

// 覆盖全局client
func SetGlobalResty(c *resty.Client) {
	globalClient = c
}

// 空client覆盖全局 PS:不重试
func SetNewGlobalResty() {
	globalClient = NewResty()
}

type Client struct {
	c       *resty.Client
	headers map[string]string
}

// 新建client
func NewClient(c *resty.Client, headers map[string]string) *Client {
	return &Client{c: c, headers: headers}
}

// 默认client
// 重试3次，每次间隔1s，超时5s
func Cli(headers map[string]string) *Client {
	return &Client{c: globalClient, headers: headers}
}

func (c *Client) post(ctx context.Context, url string, body, result interface{}) ([]byte, error) {
	resp, err := c.c.R().
		SetContext(ctx).
		SetHeaders(c.headers).
		SetBody(body).
		SetResult(result).
		Post(url)
	if err != nil {
		return nil, errors.Wrap(err, "Post")
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("post: request error code:%d , body:%s .", resp.StatusCode(), BytesToString(resp.Body()))
	}
	return resp.Body(), nil
}

func (c *Client) get(ctx context.Context, url string, result interface{}) ([]byte, error) {
	resp, err := c.c.R().
		SetContext(ctx).
		SetHeaders(c.headers).
		SetResult(result).
		Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "Get")
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Get: request error code:%d , body:%s", resp.StatusCode(), BytesToString(resp.Body()))
	}
	return resp.Body(), nil
}

func (c *Client) Post(ctx context.Context, url string, body interface{}) ([]byte, error) {
	return c.post(ctx, url, body, nil)
}

func (c *Client) PostResult(ctx context.Context, url string, body, result interface{}) ([]byte, error) {
	return c.post(ctx, url, body, result)
}

func (c *Client) PostCheckResult(ctx context.Context, url string, body, result interface{}, f CheckRespFunc) ([]byte, error) {
	resp, err := c.post(ctx, url, body, result)
	if err == nil && f != nil {
		return resp, f(result)
	}
	return resp, err
}

func (c *Client) Get(ctx context.Context, url string) ([]byte, error) {
	return c.get(ctx, url, nil)
}

func (c *Client) GetResult(ctx context.Context, url string, result interface{}) ([]byte, error) {
	return c.get(ctx, url, result)
}

func (c *Client) GetCheckResult(ctx context.Context, url string, result interface{}, f CheckRespFunc) ([]byte, error) {
	resp, err := c.get(ctx, url, result)
	if err == nil && f != nil {
		return resp, f(result)
	}
	return resp, err
}

type CheckRespFunc func(resp interface{}) error

func GetCheckRespFunc(codeKey, msgKey string, successCode int64) CheckRespFunc {
	return func(resp interface{}) error {
		if resp = Indirect(resp); resp == nil {
			return errors.New("CheckRespFunc: resp is nil")
		}
		rt := reflect.TypeOf(resp)
		rv := reflect.ValueOf(resp)

		var code int64
		var msg string
		for i := 0; i < rt.NumField(); i++ {
			switch rt.Field(i).Name {
			case codeKey:
				code = rv.Field(i).Int()
			case msgKey:
				msg = rv.Field(i).String()
			}
		}
		if code != successCode {
			return fmt.Errorf("CheckRespFunc: check error code:%d , msg:%s ", code, msg)
		}
		return nil
	}
}
