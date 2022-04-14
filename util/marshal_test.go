package util

import (
	"fmt"
	"testing"
)

type A struct {
	AA string
}

func TestMarshalObj(t *testing.T) {
	a := &A{AA: "11"}
	aP := &a
	fmt.Println(MarshalObj(&aP))
}
