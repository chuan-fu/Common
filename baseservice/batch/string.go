package batch

import (
	"time"

	"github.com/chuan-fu/Common/zlog"
)

type StringSetFunc func(data []string) error

type StringIncrease struct {
	f        StringSetFunc
	duration time.Duration

	c         chan string
	bufferLen int // chan容量大小
}

func NewStringIncrease(f StringSetFunc, duration time.Duration, buffLen int) *StringIncrease {
	s := &StringIncrease{
		f:         f,
		duration:  duration,
		c:         make(chan string, buffLen),
		bufferLen: buffLen,
	}
	go s.start()
	return s
}

// 添加
func (s *StringIncrease) Add(data string) {
	s.c <- data
}

func (s *StringIncrease) start() {
	num, max := 0, s.calcMax(s.bufferLen)
	list := make([]string, 0, max)

	handle := func() {
		if num == 0 {
			return
		}
		if err := s.f(list); err != nil {
			log.Error(err)
			return
		}
		// 容量预设，列表重置，下标重置
		max = s.calcMax(num)
		list = make([]string, 0, max)
		num = 0
	}

	ticker := time.NewTicker(s.duration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			handle()
		case data, ok := <-s.c:
			if !ok {
				handle()
				return
			}
			list = append(list, data)
			num++
			if num >= max { // 如果积存数据达到阙值，触发写入
				handle()
			}
		}
	}
}

// 预设为容量/上次间隔内写入个数的1.25倍
// 如果小于初始bufferLen，设为bufferLen
func (s *StringIncrease) calcMax(num int) int {
	if num < s.bufferLen {
		return s.bufferLen
	}
	return num + num/4
}

func (s *StringIncrease) Close() {
	close(s.c)
}
