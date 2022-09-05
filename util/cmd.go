package util

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/chuan-fu/Common/baseservice/cast"
	"github.com/chuan-fu/Common/baseservice/colorx"
	"github.com/chuan-fu/Common/baseservice/keyboardx"

	"github.com/pkg/errors"
)

func Command(name string, arg ...string) (resp string, err error) {
	return command(exec.Command(name, arg...))
}

func CommandBash(arg string) (string, error) {
	return command(exec.Command("bash", "-c", arg))
}

func CommandContext(ctx context.Context, name string, arg ...string) (resp string, err error) {
	return command(exec.CommandContext(ctx, name, arg...))
}

func CommandBashContext(ctx context.Context, arg string) (string, error) {
	return command(exec.CommandContext(ctx, "bash", "-c", arg))
}

func command(cmd *exec.Cmd) (resp string, err error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = errors.Wrap(errors.New(stderr.String()), err.Error())
		return
	}
	resp = out.String()
	return
}

func CheckStrList(list []string) (resp string) {
	switch len(list) {
	case 0:
		fmt.Println()
	case 1:
		fmt.Printf("%s\n\n", list[0])
		resp = list[0]
	default:
		f := NewFmtList()
		for k := range list {
			f.Add(list[k])
		}
		fmt.Print(f.String())

		_ = keyboardx.KeyboardX(
			func(s string) (isEnd bool, err error) {
				index, err2 := cast.ToInt(s)
				if err2 == nil {
					index -= 1
					if len(list) > index && index >= 0 {
						resp = list[index]
						return true, nil
					}
				}
				fmt.Println("---下标有误---")
				return true, nil
			},
			keyboardx.WithEmptyEnter(func(s string) (isEnd bool, err error) {
				resp = list[0]
				return true, nil
			}),
			keyboardx.WithPrefix(colorx.PurpleRedArrow),
		)
	}
	return
}
