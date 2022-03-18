package util

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestCmd(t *testing.T) {
	t1 := time.Now()
	fmt.Println(CommandContext(context.TODO(), "./main"))
	fmt.Println("use ", time.Now().Sub(t1))
}
