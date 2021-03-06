package arrayx

import (
	"fmt"
	"testing"
)

func TestI6(t *testing.T) {
	s := []int64{2, 3, 1, 4, 6, 45, 5, 5}
	Int64(s)
	fmt.Println(s)
	Int64(s, true)
	fmt.Println(s)
}

func TestDeduplice(t *testing.T) {
	fmt.Println(DistinctInt64([]int64{}))
	fmt.Println(DistinctInt64([]int64{1}))
	fmt.Println(DistinctInt64([]int64{0, 0, 0}))
	fmt.Println(DistinctInt64([]int64{1, 1, 1}))
	fmt.Println(DistinctInt64([]int64{1, 1, 2}))
	fmt.Println(DistinctInt64([]int64{1, 2, 2}))
	fmt.Println(DistinctInt64([]int64{-1, 1, -1}))
	fmt.Println(DistinctInt64([]int64{2, 3, 1, 4, 6, 5, 5}))
	fmt.Println(DistinctFloat64([]float64{2, 3, 1, 4, 6, 5, 5}))
	fmt.Println(DistinctInt([]int{2, 3, 1, 4, 6, 5, 5}))
}

func TestI(t *testing.T) {
	s := []int{2, 3, 1, 4, 6, 45, 5, 5}
	Int(s)
	fmt.Println(s)
	Int(s, true)
	fmt.Println(s)
}
