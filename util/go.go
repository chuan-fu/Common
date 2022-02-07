package util

import "github.com/chuan-fu/Common/log"

func DeferFunc() {
	if e := recover(); e != nil {
		log.Errorf("DeferFunc panic : %v", e)
	}
}

func Go(f func()) {
	go func() {
		defer DeferFunc()
		f()
	}()
}
