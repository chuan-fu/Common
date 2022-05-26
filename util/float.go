package util

import (
	"fmt"
	"math"
	"strconv"
)

const (
	amountEpsilon = 0.00001
)

// 保留两位小数，精度保护
// 直接舍弃后面位数
// 2.0499999 =》 2.05
// 2.0590001 =》 2.05
func RoundingTwo(f float64) float64 {
	f += amountEpsilon
	return math.Floor(f*100) / 100
}

func RoundingTwoStr(f float64) string {
	return fmt.Sprintf("%.2f", RoundingTwo(f))
}

func RoundTwo2(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
