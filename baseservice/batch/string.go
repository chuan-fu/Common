package batch

import (
	"time"

	"github.com/chuan-fu/Common/zlog"
)

type StringSetFunc func(data []string) error

type StringIncrease struct {
	f        StringSetFunc
	duration time.Duration

	c chan string

	// 预设的缓冲区容量bufferLen，也为最小容量，请慎重设置
	// 容量设置如果太小会频繁触发写入，太多则会增加内存消耗
	// 单位时间内，如果缓存区真实存入个数小于bufferLen，则会把下次缓冲区容量设置为bufferLen
	// 如果真实存入个数大于bufferLen，则会设置为 真实存入个数的1.25倍
	bufferLen int
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

// 缓冲区容量
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
