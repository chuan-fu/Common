package util

import (
	"fmt"
	"testing"
)

func TestRounding(t *testing.T) {
	fmt.Println(RoundingTwo(2))
	fmt.Println(RoundingTwo(2.01))
	fmt.Println(RoundingTwo(2.00999999))
	fmt.Println(RoundingTwo(2.01000001))
	fmt.Println(RoundingTwo(2.0399999))
	fmt.Println(RoundingTwo(2.0419999))
	fmt.Println(RoundingTwo(2.0590001))
	fmt.Println(RoundingTwo(2.0599999))

	fmt.Println(RoundingTwo(-2))
	fmt.Println(RoundingTwo(-2.01))
	fmt.Println(RoundingTwo(-2.00999999))
	fmt.Println(RoundingTwo(-2.01000001))

	fmt.Println(RoundingTwo(0))
}

func TestRound(t *testing.T) {
	fmt.Println(RoundFloat64(0))
	fmt.Println(RoundFloat64(0.00001))
	fmt.Println(RoundFloat64(0.5))
	fmt.Println(RoundFloat64(0.4))
	fmt.Println(RoundFloat64(1.499999999))
	fmt.Println(RoundFloat64(1.500000001))
	fmt.Println(RoundFloat64(1.45))
	fmt.Println(RoundFloat64(1.55))
}

func TestRound2(t *testing.T) {
	fmt.Println(RoundFloat64(-0.0001))
	fmt.Println(RoundFloat64(-0.5))
	fmt.Println(RoundFloat64(-0.4))
	fmt.Println(RoundFloat64(-1.499999999))
	fmt.Println(RoundFloat64(-1.500000001))
	fmt.Println(RoundFloat64(-1.45))
	fmt.Println(RoundFloat64(-1.55))
}
