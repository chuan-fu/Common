package retryx

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/chuan-fu/Common/util"

	"github.com/chuan-fu/Common/baseservice/stringx"

	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

type Client struct {
	*resty.Client
	headers map[string]string
	check   CheckRespService
}

type Option func(c *Client)

func WithResty(client *resty.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.Client = client
		}
	}
}

func WithHeaders(h map[string]string) Option {
	return func(c *Client) {
		c.headers = h
	}
}

// CheckRespService 使用请谨慎，使用之后，所有带Result的请求都会校验该规则
func WithCheckResp(check CheckRespService) Option {
	return func(c *Client) {
		c.check = check
	}
}

var globalClient = &Client{Client: resty.New()}

func SetClient(c *Client) {
	globalClient = c
}

func GetClient() *Client {
	return globalClient
}

func NewClient(opts ...Option) *Client {
	c := &Client{}
	for _, opt := range opts {
		opt(c)
	}
	if c.Client == nil {
		c.Client = resty.New()
	}
	return c
}

const (
	defaultTimeOut       = 5 * time.Second
	defaultRetryCount    = 3
	defaultRetryWaitTime = time.Second
)

func NewResty() *resty.Client {
	return resty.New()
}

func NewRestyWithThreeRetry() *resty.Client {
	client := resty.New()
	client.SetTimeout(defaultTimeOut)
	client.SetRetryCount(defaultRetryCount)
	client.SetRetryWaitTime(defaultRetryWaitTime)
	client.AddRetryCondition(func(resp *resty.Response) (isRetry bool, err error) {
		if resp == nil || resp.StatusCode() != http.StatusOK {
			return true, nil
		}
		return false, nil
	})
	return client
}

func (c *Client) post(ctx context.Context, url string, body, result interface{}) ([]byte, error) {
	resp, err := c.R().
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
	resp, err := c.R().
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
	resp, err := c.post(ctx, url, body, result)
	if err == nil && c.check != nil {
		return resp, c.check.Check(result)
	}
	return resp, err
}

// 不使用全局check，单独传入
func (c *Client) PostCheckResult(ctx context.Context, url string, body, result interface{}, check CheckRespService) ([]byte, error) {
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
	resp, err := c.get(ctx, url, result)
	if err == nil && c.check != nil {
		return resp, c.check.Check(result)
	}
	return resp, err
}

func (c *Client) GetCheckResult(ctx context.Context, url string, result interface{}, check CheckRespService) ([]byte, error) {
	resp, err := c.get(ctx, url, result)
	if err == nil && check != nil {
		return resp, check.Check(result)
	}
	return resp, err
}

type CheckRespService interface {
	Check(resp interface{}) error
}

type CheckResp struct {
	codeKey, msgKey string
	successCode     int64
}

// 是字段名称，不是tag，具体参考下面的example
func NewCheckResp(codeKey, msgKey string, successCode int64) CheckRespService {
	return &CheckResp{
		codeKey:     codeKey,
		msgKey:      msgKey,
		successCode: successCode,
	}
}

func (c *CheckResp) Check(resp interface{}) error {
	if resp = util.Indirect(resp); resp == nil {
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
