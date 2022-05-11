package tokenlimit

import (
	"time"

	"github.com/pkg/errors"
)

const (
	defaultPer      = time.Second
	pingInterval    = time.Second
	tokenFormat     = "%s:tokens"
	timestampFormat = "%s:ts"
)

type config struct {
	per                          time.Duration
	pingInterval                 time.Duration // ping间隔
	tokenFormat, timestampFormat string
}

// 计算、校验token生成间隔
func (c *config) calcInterval(rate, burst int) (time.Duration, error) {
	if rate < 1 || burst < 1 || burst < rate {
		return 0, errors.New("per、rate、burst有误")
	}
	interval := c.per / time.Second * time.Duration(rate)
	if fillTime := float64(burst) * interval.Seconds(); fillTime < 1 { // 填满桶时间不到1s，小于redis的ttl最小时间单位）
		return 0, errors.New("填满桶时间不到1s，小于redis的ttl最小时间单位")
	}
	return interval, nil
}

type Option func(c *config)

// buildConfig combines defaults with options.
func buildConfig(opts []Option) config {
	c := &config{
		per:             defaultPer,
		pingInterval:    pingInterval,
		tokenFormat:     tokenFormat,
		timestampFormat: timestampFormat,
	}
	for _, opt := range opts {
		opt(c)
	}
	return *c
}

// 允许的最小单位为1s
// 部分原因为：redis最小存储时间单位为1s
func WithPer(per time.Duration) Option {
	if per < defaultPer {
		per = defaultPer
	}
	return func(c *config) {
		c.per = per
	}
}

func WithPingInterval(pingInterval time.Duration) Option {
	return func(c *config) {
		c.pingInterval = pingInterval
	}
}

func WithTokenFormat(tokenFormat string) Option {
	return func(c *config) {
		c.tokenFormat = tokenFormat
	}
}

func WithTimestampFormat(timestampFormat string) Option {
	return func(c *config) {
		c.timestampFormat = timestampFormat
	}
}
