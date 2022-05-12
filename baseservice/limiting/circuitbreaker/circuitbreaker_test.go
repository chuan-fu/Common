package circuitbreaker

import (
	"testing"

	"github.com/afex/hystrix-go/hystrix"
)

func TestHy(t *testing.T) {
	hystrix.ConfigureCommand("default", hystrix.CommandConfig{
		Timeout:                1000, // 单次请求 超时时间
		MaxConcurrentRequests:  1,    // 最大并发量
		SleepWindow:            5000, // 熔断后多久去尝试服务是否可用
		RequestVolumeThreshold: 1,    // 验证熔断的 请求数量, 10秒内采样
		ErrorPercentThreshold:  1,    // 验证熔断的 错误百分比
	})

	// hystrix.Do()
}
