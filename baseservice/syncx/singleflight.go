package syncx

import (
	"sync"
	"time"

	"github.com/chuan-fu/Common/util"
	"github.com/pkg/errors"
)

var TimeOutErr = errors.New("SingleFight Run timeOut")

// 使用单飞时，如使用ctx超时控制
// 可能会因为超时导致所有等待任务都失败
type SingleFight interface {
	Do(key string, fn SingleFightFunc) (v interface{}, err error)
	DoEx(key string, fn SingleFightFunc) (v interface{}, fresh bool, err error)
}

type SingleFightFunc func() (interface{}, error)

type call struct {
	wg sync.WaitGroup

	val  interface{} // 调用结果
	err  error       // 错误信息
	dups int         // 请求等待个数
}

type singleFlight struct {
	mu      sync.Mutex
	m       map[string]*call
	timeout time.Duration
}

func NewSingleFlight() SingleFight {
	return &singleFlight{
		m: make(map[string]*call),
	}
}

func NewSingleFlightWithTimeout(t time.Duration) SingleFight {
	return &singleFlight{
		m:       make(map[string]*call),
		timeout: t,
	}
}

func (g *singleFlight) Do(key string, fn SingleFightFunc) (v interface{}, err error) {
	c, needDo := g.checkCall(key)
	if needDo {
		g.doCall(c, key, fn)
	}
	return c.val, c.err
}

func (g *singleFlight) DoEx(key string, fn SingleFightFunc) (v interface{}, fresh bool, err error) {
	c, isFresh := g.checkCall(key)
	if isFresh {
		g.doCall(c, key, fn)
	}
	return c.val, isFresh, c.err
}

// 返回是否需要新建 needDo
func (g *singleFlight) checkCall(key string) (*call, bool) {
	g.mu.Lock()
	if c, ok := g.m[key]; ok {
		c.dups++
		g.mu.Unlock()
		c.wg.Wait()
		return c, false
	}
	c := new(call)
	c.wg.Add(1) // 要在未解锁前加，不然小概率出现，未Add() 直接Wait()返回的可能
	g.m[key] = c
	g.mu.Unlock()
	return c, true
}

func (g *singleFlight) doCall(c *call, key string, fn SingleFightFunc) {
	defer func() {
		g.mu.Lock()
		delete(g.m, key)
		g.mu.Unlock()

		c.wg.Done() // 结束任务
	}()

	g.runFn(c, fn)
}

func (g *singleFlight) runFn(c *call, fn SingleFightFunc) {
	if g.timeout <= 0 {
		g.safeCallFn(c, fn, nil)
		return
	}

	timeout := time.NewTimer(g.timeout) // 超时控制
	defer timeout.Stop()
	dChan := NewDoneChan() // 结束控制

	// 执行任务
	go g.safeCallFn(c, fn, dChan)

	select {
	case <-timeout.C: // 超时关闭
		c.err = TimeOutErr
	case <-dChan.Done(): // 结束关闭
	}
}

func (g *singleFlight) safeCallFn(c *call, fn SingleFightFunc, dChan DoneChan) {
	defer func() {
		if e := recover(); e != nil { // 捕获异常
			c.err = util.NewPanicError(e)
			if dChan != nil { // 关闭
				dChan.Close()
			}
		}
	}()
	c.val, c.err = fn()
	if dChan != nil { // 关闭
		dChan.Close()
	}
}
