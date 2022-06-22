package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	UpperLetter = 1 << iota // 1 大写字母
	LowerLetter             // 2 小写字母
	Number                  // 4 数字
	SpecialChar             // 8 特殊字符
)

const (
	LetterChar    = UpperLetter | LowerLetter
	LetterNumChar = UpperLetter | LowerLetter | Number
	AllChar       = UpperLetter | LowerLetter | Number | SpecialChar
)

var (
	randSeed     = time.Now().UnixNano()
	randSource   = rand.NewSource(randSeed)
	randInstance = rand.New(randSource)

	upperLetterArr = []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
	lowerLetterArr = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	numberArr      = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	specialCharArr = []rune{'!', '@', '#', '$', '%', '^', '&', '*'}
)

type RandomParam struct {
	arrList []rune
	len     int
}

func NewRandomParam(char int) *RandomParam {
	r := &RandomParam{}
	if char&UpperLetter > 0 {
		r.arrList = append(r.arrList, upperLetterArr...)
	}
	if char&LowerLetter > 0 {
		r.arrList = append(r.arrList, lowerLetterArr...)
	}
	if char&Number > 0 {
		r.arrList = append(r.arrList, numberArr...)
	}
	if char&SpecialChar > 0 {
		r.arrList = append(r.arrList, specialCharArr...)
	}
	r.len = len(r.arrList)
	return r
}

// 生成随机字符串
func (r *RandomParam) GenRandomKey(strLen int) string {
	if r.len == 0 || strLen == 0 {
		return ""
	}
	sb := strings.Builder{}
	for i := 0; i < strLen; i++ {
		randInt := randInstance.Intn(r.len)
		sb.WriteRune(r.arrList[randInt])
	}
	return sb.String()
}
