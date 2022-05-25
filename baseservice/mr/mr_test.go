package mr

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/mr"

	"github.com/pkg/errors"

	log "github.com/chuan-fu/Common/zlog"
)

func TestMr(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	now := time.Now()
	err := FinishWithCtx(ctx, func() (err error) {
		time.Sleep(time.Second)
		fmt.Println("1")
		return
	}, func() (err error) {
		time.Sleep(2 * time.Second)
		fmt.Println("2")
		return errors.New("err2")
	}, func() (err error) {
		time.Sleep(3 * time.Second)
		fmt.Println("3")
		return
	})
	if err != nil {
		log.Error(err)
	}
	fmt.Println(time.Now().Sub(now))
	time.Sleep(2 * time.Second)
}

func TestMr2(t *testing.T) {
	s := sync.WaitGroup{}
	s.Add(1)
	s.Done()
	now := time.Now()
	err := mr.Finish(func() (err error) {
		time.Sleep(time.Second)
		fmt.Println("1")
		return
	}, func() (err error) {
		time.Sleep(2 * time.Second)
		fmt.Println("2")
		return errors.New("err2")
	}, func() (err error) {
		time.Sleep(3 * time.Second)
		fmt.Println("3")
		return
	})
	if err != nil {
		log.Error(err)
	}
	fmt.Println(time.Now().Sub(now))
	time.Sleep(2 * time.Second)
}
