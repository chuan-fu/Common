package retryx

type (
	retryConfig struct {
		retryNum int                    // 最小为1，表示不重试
		isRetry  func(interface{}) bool // 可以默认返回false，表示，非报错则不重试
	}

	Option func(r *retryConfig)
)

var (
	DefaultRetryNum = 3                                       // 默认重试三次
	DefaultIsRetry  = func(interface{}) bool { return false } // 默认 非报错不重试
)

func WithRetryNum(num int) Option {
	return func(r *retryConfig) {
		if num > 0 {
			r.retryNum = num
		}
	}
}

func WithIsRetryFunc(f func(interface{}) bool) Option {
	return func(r *retryConfig) {
		r.isRetry = f
	}
}

func buildConfig(opts []Option) *retryConfig {
	r := &retryConfig{
		retryNum: DefaultRetryNum,
		isRetry:  DefaultIsRetry,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}
