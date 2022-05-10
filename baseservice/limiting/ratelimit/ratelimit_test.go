package ratelimit

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// 使用uber开发的go.uber.org/ratelimit@v0.2.0做二次开发，增加了超时控制
// 堵塞式限流器
// 把每个请求排队，每隔1s/rate放出一个请求
// Take()接口会堵塞，直到请求允许放行，返回允许的时间
// 缺陷在于，瞬时流量过大会存在请求饿死，长时间等待【已优化】
func TestRateLimit(t *testing.T) {
	r := New(10, WithoutSlack)
	now := time.Now()
	ctx, _ := context.WithDeadline(context.TODO(), now.Add(2*time.Second))
	for i := 0; i < 50; i++ {
		go func(index int) {
			if n := r.TakeWithContext(ctx); n.IsZero() {
				fmt.Println(index, "超时")
			} else {
				fmt.Println(index, n.Sub(now))
			}
		}(i)
	}
	time.Sleep(10 * time.Second)
}

func TestSleep(t *testing.T) {
	now := time.Now()
	time.Sleep(-100 * time.Second)
	fmt.Println(time.Now().Sub(now))
}
