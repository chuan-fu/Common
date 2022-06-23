package batch

import (
	"fmt"
	"testing"
	"time"

	"github.com/chuan-fu/Common/baseservice/cast"
)

func TestString(t *testing.T) {
	s := NewStringIncrease(func(data []string) error {
		fmt.Println("set data =>", data)
		return nil
	}, time.Second, 10)

	for i := 0; i < 1000; i++ {
		time.Sleep(10 * time.Millisecond * time.Duration(i%10))
		s.Add(cast.ToString(i + 1))
	}
	time.Sleep(time.Second * 10)
	for i := 1000; i < 2000; i++ {
		time.Sleep(5 * time.Millisecond * time.Duration(i%10))
		s.Add(cast.ToString(i + 1))
	}
	time.Sleep(time.Second * 10)
}
