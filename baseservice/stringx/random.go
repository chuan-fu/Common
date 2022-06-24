package stringx

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

type arrInfo struct {
	data []rune
	len  int
}

var (
	randInstance                                              *rand.Rand
	upperLetterArr, lowerLetterArr, numberArr, specialCharArr *arrInfo
)

func init() {
	randSource := rand.NewSource(time.Now().UnixNano())
	randInstance = rand.New(randSource)

	upperLetterArr = &arrInfo{
		data: []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'},
	}
	upperLetterArr.len = len(upperLetterArr.data)

	lowerLetterArr = &arrInfo{
		data: []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'},
	}
	lowerLetterArr.len = len(lowerLetterArr.data)

	numberArr = &arrInfo{
		data: []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'},
	}
	numberArr.len = len(numberArr.data)

	specialCharArr = &arrInfo{
		data: []rune{'!', '@', '#', '$', '%', '^', '&', '*'},
	}
	specialCharArr.len = len(specialCharArr.data)
}

type RandomParam interface {
	GenRandomKey(strLen int) string
}

type randomParam struct {
	arrList []*arrInfo
	len     int
}

func NewRandomParam(char int) RandomParam {
	r := &randomParam{
		arrList: make([]*arrInfo, 0),
		len:     0,
	}
	if char&UpperLetter > 0 {
		r.arrList = append(r.arrList, upperLetterArr)
		r.len += upperLetterArr.len
	}
	if char&LowerLetter > 0 {
		r.arrList = append(r.arrList, lowerLetterArr)
		r.len += lowerLetterArr.len
	}
	if char&Number > 0 {
		r.arrList = append(r.arrList, numberArr)
		r.len += numberArr.len
	}
	if char&SpecialChar > 0 {
		r.arrList = append(r.arrList, specialCharArr)
		r.len += specialCharArr.len
	}
	return r
}

// 生成随机字符串
func (r *randomParam) GenRandomKey(strLen int) string {
	if r.len == 0 || strLen == 0 {
		return ""
	}
	sb := strings.Builder{}
	for i := 0; i < strLen; i++ {
		sb.WriteRune(r.getRune())
	}
	return sb.String()
}

func (r *randomParam) getRune() rune {
	randInt := randInstance.Intn(r.len)
	for k := range r.arrList {
		if randInt >= r.arrList[k].len {
			randInt -= r.arrList[k].len
			continue
		}
		return r.arrList[k].data[randInt]
	}
	return 'A'
}
