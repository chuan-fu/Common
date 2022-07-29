package keyboardx

import "github.com/chuan-fu/Common/util"

const (
	DefaultBufferSize = 10
	CmdHistory        = "history"
)

type config struct {
	prefix                string
	grep                  int      // 颜色
	bufferSize            int      // 缓冲区
	cmdList               []string // 默认指令
	taskSvc               TaskService
	emptyEnter            Task // 空回车逻辑处理
	preHandle, postHandle func(s string)
	checkInHistoryHandle  func(s string) (string, bool) // 校验是否加入历史
	needHistory           bool
}

type Option func(c *config)

func buildConfig(opts []Option) *config {
	c := &config{
		prefix:      util.GreenBluePrefix,
		grep:        util.RedGrep,
		bufferSize:  DefaultBufferSize,
		cmdList:     make([]string, 0),
		needHistory: true,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithPrefix(prefix string) Option {
	return func(c *config) {
		c.prefix = prefix
	}
}

func WithGrep(grep int) Option {
	return func(c *config) {
		c.grep = grep
	}
}

func WithBufferSize(bufferSize int) Option {
	return func(c *config) {
		c.bufferSize = bufferSize
	}
}

func WithCmdList(cmdList []string) Option {
	return func(c *config) {
		c.cmdList = cmdList
	}
}

func WithTask(taskSvc TaskService) Option {
	return func(c *config) {
		c.taskSvc = taskSvc
	}
}

func WithEmptyEnter(t Task) Option {
	return func(c *config) {
		c.emptyEnter = t
	}
}

func WithPreHandle(f func(s string)) Option {
	return func(c *config) {
		c.preHandle = f
	}
}

func WithPostHandle(f func(s string)) Option {
	return func(c *config) {
		c.postHandle = f
	}
}

func WithCheckInHistory(f func(s string) (string, bool)) Option {
	return func(c *config) {
		c.checkInHistoryHandle = f
	}
}

func WithNeedHistory(b bool) Option {
	return func(c *config) {
		c.needHistory = b
	}
}
