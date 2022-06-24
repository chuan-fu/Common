package syncx

import (
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestSingleFightDo(t *testing.T) {
	g := NewSingleFlightWithTimeout(3 * time.Second)
	// g := NewSingleFlight()

	for i := 0; i < 5; i++ {
		go func(index int) {
			val, share, err := g.DoEx("key1", func() (interface{}, error) {
				fmt.Println(index, "key1 do it")
				time.Sleep(2 * time.Second)
				return "result1", nil
			})
			fmt.Println(index, "key1", val, err, share)
		}(i)

		go func(index int) {
			val, share, err := g.DoEx("key2", func() (interface{}, error) {
				fmt.Println(index, "key2 do it")
				time.Sleep(2 * time.Second)
				panic(2)
				return "result2", nil
			})
			fmt.Println(index, "key2", val, err, share)
		}(i)

		go func(index int) {
			val, share, err := g.DoEx("key3", func() (interface{}, error) {
				fmt.Println(index, "key3 do it")
				time.Sleep(2 * time.Second)
				return nil, errors.New("errResult3")
			})
			fmt.Println(index, "key3", val, err, share)
		}(i)
	}
	time.Sleep(5 * time.Second)
}

func TestSingleFightDoWithTimeout(t *testing.T) {
	g := NewSingleFlightWithTimeout(3 * time.Second)
	_ = g
}
