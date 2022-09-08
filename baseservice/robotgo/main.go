package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
	/*
		// 模拟按下1个键：打开开始菜单（win）
		robotgo.KeyTap(`command`)
		// 模拟按下2个键：打开资源管理器（win + e）
		robotgo.KeyTap(`e`, `command`)
		// 模拟按下3个键：打开任务管理器（Ctrl + Shift + ESC）
		robotgo.KeyTap(`esc`, `control`, `shift`)
	*/
	fmt.Println(robotgo.KeyTap(`t`, `command`))
}
