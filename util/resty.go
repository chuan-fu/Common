package util

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/chuan-fu/Common/baseservice/stringx"

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
		return nil, fmt.Errorf("post: request error url:%s , code:%d , body:%s", url, resp.StatusCode(), stringx.BytesToString(resp.Body()))
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
		return nil, fmt.Errorf("get: request error url:%s , code:%d , body:%s", url, resp.StatusCode(), stringx.BytesToString(resp.Body()))
	}
	return resp.Body(), nil
}

func (c *Client) Post(ctx context.Context, url string, body interface{}) ([]byte, error) {
	return c.post(ctx, url, body, nil)
}

func (c *Client) PostResult(ctx context.Context, url string, body, result interface{}) ([]byte, error) {
	return c.post(ctx, url, body, result)
}

func (c *Client) PostCheckResult(ctx context.Context, url string, body, result interface{}, check CheckTypeService) ([]byte, error) {
	resp, err := c.post(ctx, url, body, result)
	if err == nil && check != nil {
		return resp, check.Check(result)
	}
	return resp, err
}

func (c *Client) Get(ctx context.Context, url string) ([]byte, error) {
	return c.get(ctx, url, nil)
}

func (c *Client) GetResult(ctx context.Context, url string, result interface{}) ([]byte, error) {
	return c.get(ctx, url, result)
}

func (c *Client) GetCheckResult(ctx context.Context, url string, result interface{}, check CheckTypeService) ([]byte, error) {
	resp, err := c.get(ctx, url, result)
	if err == nil && check != nil {
		return resp, check.Check(result)
	}
	return resp, err
}

type CheckTypeService interface {
	Check(resp interface{}) error
}

type CheckResp struct {
	codeKey, msgKey string
	successCode     int64
}

// 是字段名称，不是tag，具体参考下面的example
func NewCheckResp(codeKey, msgKey string, successCode int64) CheckTypeService {
	return &CheckResp{
		codeKey:     codeKey,
		msgKey:      msgKey,
		successCode: successCode,
	}
}

func (c *CheckResp) Check(resp interface{}) error {
	if resp = Indirect(resp); resp == nil {
		return errors.New("CheckResp: resp is nil")
	}
	rt := reflect.TypeOf(resp)
	rv := reflect.ValueOf(resp)

	var code int64
	var msg string
	for i := 0; i < rt.NumField(); i++ {
		switch rt.Field(i).Name {
		case c.codeKey:
			code = rv.Field(i).Int()
		case c.msgKey:
			msg = rv.Field(i).String()
		}
	}
	if code != c.successCode {
		return fmt.Errorf("CheckResp: check error code:%d , msg:%s , code should be %d", code, msg, c.successCode)
	}
	return nil
}

// example
type exampleResp struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

var exampleCheckResp = NewCheckResp("Code", "Msg", 0)
