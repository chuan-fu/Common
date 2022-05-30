package signal

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chuan-fu/Common/baseservice/mr"
	"github.com/chuan-fu/Common/util"
	"github.com/chuan-fu/Common/zlog"
)

var (
	sh        *shutDown
	signalMap = map[os.Signal]string{
		syscall.SIGTERM: "syscall.SIGTERM",
		syscall.SIGINT:  "syscall.SIGINT",
		syscall.SIGQUIT: "syscall.SIGQUIT",
	}
)

func init() {
	sh = NewShutDown()
}

func AddShutDownFunc(f func() error) {
	sh.Add(f)
}

// 该方法要作为主函数main
func Run(opts ...SignalOption) {
	defer util.DeferFunc()

	c := buildConfig(opts)

	log.Info(logPrefix, "Application SignalClose Running")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	s := <-ch // 等待关闭信号量
	log.Info(logPrefix, "Application trigger close signal: ", signalMap[s])

	mr.FinishVoid(
		func() {
			sh.Run() // 调用所有注册关闭接口
			log.Info(logPrefix, "Application shut down register close")
		},
		func() {
			if c.middleware != nil { // 等待中间件关闭，所有残余流量请求结束
				c.middleware.ShutDown()
				log.Info(logPrefix, "Application shut down middleware close")
			}
		},
		func() { // 确保等待时间至少超过shutdownDelay
			time.Sleep(c.shutdownDelay)
		},
	)

	log.Info(logPrefix, "Application is shut down")
	os.Exit(0)
}

type config struct {
	shutdownDelay time.Duration
	middleware    ServerClose
}

func buildConfig(opts []SignalOption) *config {
	c := &config{
		shutdownDelay: shutdownDelay,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type SignalOption func(c *config)

func WithShotDownDelay(t time.Duration) SignalOption {
	return func(c *config) {
		if t > 0 {
			c.shutdownDelay = t
		}
	}
}

func WithMiddleware(s ServerClose) SignalOption {
	return func(c *config) {
		c.middleware = s
	}
}
