package util

import (
	"fmt"
	"testing"

	"github.com/chuan-fu/Common/log"
)

func TestReadFileLine(t *testing.T) {
	c, err := ReadFileLine("./file.go")
	if err != nil {
		log.Error(err)
		return
	}
	for {
		s, ok := <-c
		if !ok {
			break
		}
		fmt.Println(string(s))
	}
}

func TestWriteFile(t *testing.T) {
	WriteFile("./11.txt", []byte("asasd11"))
}

func TestAppendFile(t *testing.T) {
	fmt.Println(AppendFile("./11.txt", "asasd11"))
}

func TestReadPath(t *testing.T) {
	c := ReadPath("/tmp")
	for {
		s, ok := <-c
		if !ok {
			return
		}
		fmt.Println(s)
	}
}
