package keyboardx

const (
	prefix            = "-> # "
	redGrep           = 0x1B // 红色
	defaultBufferSize = 10
)

type config struct {
	prefix                string
	grep                  int
	bufferSize            int      // 缓冲区
	cmdList               []string // 默认指令
	taskSvc               TaskService
	emptyEnter            Task // 空回车逻辑处理
	preHandle, postHandle func(s string)
	checkInHistoryHandle  func(s string) (string, bool) // 校验是否加入历史
}

type Option func(c *config)

func buildConfig(opts []Option) *config {
	c := &config{
		prefix:     prefix,
		grep:       redGrep,
		bufferSize: defaultBufferSize,
		cmdList:    make([]string, 0),
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
