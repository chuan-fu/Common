package log

import (
	"testing"
)

func TestZlog(t *testing.T) {
	ReloadLogger(ZlogConf{
		SysName:  "dpm",
		Encoding: "console",
	})
	Debug("aa")
	Info("aa")
	Warn("aa")
	Error("aa")
	Fatal("aaa")
}
