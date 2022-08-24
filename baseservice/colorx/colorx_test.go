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

func TestMottle(t *testing.T) {
	fmt.Printf("%c[0;32m 11111 %c[0m\n", 0x1B, 0x1B)
	fmt.Printf("%c[0;32;40m 22222 %c[0m\n", 0x1B, 0x1B)
	fmt.Printf("%c[1;32m 33333 %c[0m\n", 0x1B, 0x1B)
	fmt.Printf("%c[1;32;40m 44444 %c[0m\n", 0x1B, 0x1B)
}
