package util

import (
	"fmt"
	"testing"
)

func TestTT(*testing.T) {
	r := NewRandomParam(AllChar)
	for i := 0; i < 100; i++ {
		fmt.Println(r.GenRandomKey(6))
	}
}

const (
	MessageSubBitNewItem        = 0b001
	MessageSubBitDailyRecommend = 0b010
	MessageSubBitVipDay         = 0b100
)

func TestHH(t *testing.T) {
	fmt.Println(LetterChar & UpperLetter)
	fmt.Println(LetterChar & LowerLetter)
	fmt.Println(LetterChar & Number)
	fmt.Println(LetterChar & SpecialChar)
	fmt.Println(MessageSubBitNewItem, MessageSubBitDailyRecommend, MessageSubBitVipDay)
}
