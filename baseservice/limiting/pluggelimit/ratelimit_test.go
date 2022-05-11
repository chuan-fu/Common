package pluggelimit

import (
	"fmt"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	r := New(10, WithPer(time.Minute))
	now := time.Now()
	for i := 0; i < 50; i++ {
		go func(index int) {
			if n := r.Take(); n.IsZero() {
				fmt.Println(index, "超时")
			} else {
				fmt.Println(index, n.Sub(now))
			}
		}(i)
	}
	time.Sleep(time.Minute)
}

func TestSleep(t *testing.T) {
	now := time.Now()
	time.Sleep(-100 * time.Second)
	fmt.Println(time.Now().Sub(now))
}
