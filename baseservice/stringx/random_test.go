package stringx

import (
	crand "crypto/rand"
	"fmt"
	"testing"
)

func TestTT(*testing.T) {
	//for i := 0; i < 100; i++ {
	//	fmt.Println(randInstance.Intn(5))
	//}

	r := NewRandomParam(AllChar)
	for i := 0; i < 100; i++ {
		fmt.Println(r.GenRandomKey(8))
	}
}

func TestCrypto(t *testing.T) {
	b := make([]byte, 8)
	crand.Read(b)
	fmt.Printf("%x%x%x%x", b[0:2], b[2:4], b[4:6], b[6:8])
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
