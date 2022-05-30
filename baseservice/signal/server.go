package signal

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerClose interface {
	ShutDown()
}

const (
	stateRun   = 0
	stateWait  = 1
	stateClose = 2
)

type ServerCloseMiddleware struct {
	reqNum       int32
	state        int32
	waitTime     time.Time
	waitDuration time.Duration
}

type ServerOption func(s *ServerCloseMiddleware)

func WithWaitDuration(waitDuration time.Duration) ServerOption {
	return func(s *ServerCloseMiddleware) {
		if s.waitDuration < 0 {
			s.waitDuration = 0
		}
		s.waitDuration = waitDuration
	}
}

func NewServerCloseMiddleware(opts ...ServerOption) *ServerCloseMiddleware {
	s := &ServerCloseMiddleware{
		state:        stateRun,
		waitDuration: defaultWaitDuration,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (g *ServerCloseMiddleware) ShutDown() {
	if atomic.CompareAndSwapInt32(&g.state, stateRun, stateWait) {
		g.waitTime = time.Now().Add(g.waitDuration)
	}
	for {
		if atomic.LoadInt32(&g.state) == stateWait && time.Now().After(g.waitTime) { // 超时则切为close状态
			atomic.CompareAndSwapInt32(&g.state, stateWait, stateClose)
		}
		if atomic.LoadInt32(&g.state) == stateClose && atomic.LoadInt32(&g.reqNum) == 0 {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (g *ServerCloseMiddleware) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !g.check() {
			c.AbortWithStatus(http.StatusNotImplemented)
			return
		}

		atomic.AddInt32(&g.reqNum, 1)
		defer atomic.AddInt32(&g.reqNum, -1)

		c.Next()
	}
}

func (g *ServerCloseMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !g.check() {
		rw.WriteHeader(http.StatusNotImplemented)
		return
	}
	atomic.AddInt32(&g.reqNum, 1)
	defer atomic.AddInt32(&g.reqNum, -1)

	next(rw, r)
}

// 不为关闭状态则通过
func (g *ServerCloseMiddleware) check() bool {
	return atomic.LoadInt32(&g.state) != stateClose
}
