package keyboardx

import (
	"fmt"
	"testing"

	"github.com/chuan-fu/Common/baseservice/colorx"
)

func TestCommandHistory(t *testing.T) {
	c := newCommandHistory([]string{"dev", "test"}, colorx.WordRed)
	c.add("dpm")
	fmt.Println(c.setIndex(-1))
	c.add("dpm")
	fmt.Println(c.setIndex(-1))
	fmt.Println(c.find("d", 0))
}
