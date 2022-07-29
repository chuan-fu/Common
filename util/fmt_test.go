package util

import (
	"fmt"
	"testing"
)

func TestFmtColor(t *testing.T) {
	fmt.Println(BlackArrow)
	fmt.Println(RedArrow)
	fmt.Println(GreenArrow)
	fmt.Println(YellowArrow)
	fmt.Println(BlueArrow)
	fmt.Println(PurpleRedArrow)
	fmt.Println(GreenBlueArrow)
	fmt.Println(WhiteArrow)

	fmt.Println(RedBluePrefix)
	fmt.Println(GreenBluePrefix)
	fmt.Println(FmtColor("avc", "a", RedGrep))
}
