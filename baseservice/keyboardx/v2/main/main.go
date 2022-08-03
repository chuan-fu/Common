package main

import (
	"fmt"

	"github.com/chuan-fu/Common/baseservice/keyboardx/v2"
	"github.com/chuan-fu/Common/zlog"
)

type DevServ struct {
	keyboardx.BaseTask
}

func (d *DevServ) PreHandle(s string) error {
	fmt.Println("PreHandle:", s)
	return nil
}

func (d *DevServ) PostHandle(s string) error {
	fmt.Println("PostHandle:", s)
	return nil
}

func (d *DevServ) Handle(s string) error {
	fmt.Println("\nKeyboardX 1 start =>", s)
	return nil
}

func main() {
	keyboardxSvc := keyboardx.NewKeyboardX().
		ResetCmdList([]string{"dev", "test", "pre"}).
		AddHelp().
		AddHistory().
		AddExit()

	keyboardxSvc.AddFullyTask("dev", "dev desc", &DevServ{})

	if err := keyboardxSvc.Run(); err != nil {
		log.Error(err)
		return
	}
}
