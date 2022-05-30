package signal

import "time"

const (
	defaultWaitDuration = 2 * time.Second // 延迟2秒关闭流量进入
	shutdownDelay       = 5 * time.Second // 接到信号至少沉睡5秒后关闭服务
	logPrefix           = "wf04u5ezfpxd "
)
