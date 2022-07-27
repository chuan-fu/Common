package main

import (
	"fmt"

	"github.com/chuan-fu/Common/baseservice/keyboardx"
	"github.com/chuan-fu/Common/zlog"
)

func main() {
	err := keyboardx.KeyboardX(func(s string) (isEnd bool, err error) {
		fmt.Println("\nKeyboardX 1 start =>", s)
		err = keyboardx.KeyboardX(func(s string) (isEnd bool, err error) {
			fmt.Println("\nKeyboardX 2 =>", s)
			return
		})
		if err != nil {
			log.Error(err)
		}
		fmt.Println("\nKeyboardX 1 end =>", s)
		return
	})
	if err != nil {
		log.Error(err)
		return
	}
}
