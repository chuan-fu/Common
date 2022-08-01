package colorx

import "testing"

func TestColorPrint(t *testing.T) {
	Print(WordGreen, "AAA")
	Println(WordRed, "AAA")
	Printf(WordYellow, "A%dB", 1)
}
