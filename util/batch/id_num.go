package batch

import (
	"time"

	"github.com/chuan-fu/Common/zlog"
)

type idNum struct {
	id  int64
	num int64
}

type IdNumSetFunc func(id, num int64) error

type IdNumIncrease struct {
	f        IdNumSetFunc // 处理函数
	duration time.Duration
	c        chan *idNum
}

// 批处理-基于id的数量变化的统计合并处理
func NewIdNumIncrease(f IdNumSetFunc, duration time.Duration, buffLen int) *IdNumIncrease {
	i := &IdNumIncrease{
		f:        f,
		duration: duration,
		c:        make(chan *idNum, buffLen),
	}
	go i.start()
	return i
}

// 每次添加数量
// 也可以小于0
func (i *IdNumIncrease) AddNum(id, num int64) {
	i.c <- &idNum{
		id:  id,
		num: num,
	}
}

func (i *IdNumIncrease) start() {
	numMap := make(map[int64]*idNum)
	handle := func() {
		for _, num := range numMap {
			if num.num != 0 {
				err := i.f(num.id, num.num)
				if err != nil {
					log.Error(err)
					break
				}
			}
			delete(numMap, num.id)
		}
	}

	ticker := time.NewTicker(i.duration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			handle()
		case num, ok := <-i.c:
			if !ok {
				handle()
				return
			}
			if numInMap, ok2 := numMap[num.id]; ok2 {
				numInMap.num += num.num
			} else {
				numMap[num.id] = num
			}
		}
	}
}

func (i *IdNumIncrease) Close() {
	close(i.c)
}
