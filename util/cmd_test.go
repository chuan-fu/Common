package util

import (
	"fmt"
	"testing"
	"time"
)

func TestCmd(t *testing.T) {
	t1 := time.Now()
	fmt.Println(CommandContext("./main"))
	fmt.Println("use ", time.Now().Sub(t1))
}
