package stringx

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

func TestPoint(t *testing.T) {
	fmt.Println("Hello, 世界", len("世界"), utf8.RuneCountInString("世界"))
}

func TestTrimByte(t *testing.T) {
	fmt.Println(TrimByte("", 'a'))
	fmt.Println(TrimByte("aaabcav", 'a'))
	fmt.Println(ReplaceByte("aaabcav", 'a', 'd'))
}
