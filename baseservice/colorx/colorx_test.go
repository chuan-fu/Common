package colorx

import (
	"fmt"
	"testing"
)

func TestColorPrint(t *testing.T) {
	Print(WordGreen, "AAA")
	Println(WordRed, "AAA")
	Printf(WordYellow, "A%dB", 1)
}

func TestKeywordsSprintf(t *testing.T) {
	fmt.Println(KeywordsSprintf(WordGreen, "ABCA", "A"))
	fmt.Println(KeywordsSprintf(WordRed, "ABCA", "A"))
}
