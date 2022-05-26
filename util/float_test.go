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
