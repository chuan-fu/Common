package keyboardx

import (
	"fmt"
	"testing"
)

func TestCommandHistory(t *testing.T) {
	c := newCommandHistory([]string{"dev", "test"}, 0x1B)
	c.add("dpm")
	fmt.Println(c.setIndex(-1))
	c.add("dpm")
	fmt.Println(c.setIndex(-1))
}
