package zlog

import (
	"testing"
)

func TestZlog(t *testing.T) {
	ReloadLogger(ZlogConf{
		SysName: "dpm",
	})
	Debug("aa")
	Info("aa")
	Warn("aa")
	Error("aa")
	Panic("aaa")
}
