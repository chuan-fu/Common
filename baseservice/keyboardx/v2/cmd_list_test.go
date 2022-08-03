package keyboardx

import (
	"fmt"
	"testing"
)

func TestNewCmdListServ(t *testing.T) {
	cmd := newCmdListServ([]string{"test", "dev", "devd"})
	fmt.Println(cmd.Last()) // devd
	fmt.Println(cmd.Last()) // dev
	cmd.Add("pre")
	cmd.Add("pre")
	fmt.Println(cmd.Last()) // pre
	fmt.Println(cmd.Next()) // 空
	fmt.Println(cmd.Last()) // pre
	fmt.Println(cmd.Last()) // devd
	fmt.Println(cmd.Last()) // dev
	fmt.Println(cmd.Last()) // test
	fmt.Println(cmd.Last()) // test
	fmt.Println(cmd.FindStart("d"))
	fmt.Println(cmd.FindStart("d"))
	fmt.Println(cmd.FindStart("d"))
	cmd.FindClose()
	runTask(cmd.HistoryPrefixTask(""), "")
	fmt.Println("------")
	runTask(cmd.HistoryPrefixTask(""), "d")
}
