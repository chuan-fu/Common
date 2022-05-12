package batch

import (
	"fmt"
	"testing"
	"time"
)

func TestNum(t *testing.T) {
	n := NewIdNumIncrease(func(id, num int64) error {
		fmt.Println("add", id, num)
		return nil
	}, time.Second, 100)
	for i := 0; i < 200; i++ {
		time.Sleep(100 * time.Millisecond)
		if i%2 == 0 {
			n.AddNum(int64(i%10), int64(i))
		} else {
			n.AddNum(int64(i%10), -int64(i))
		}

	}
	time.Sleep(time.Minute)
}
