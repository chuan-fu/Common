package util

import (
	"fmt"
	"testing"

	"github.com/ahmetb/go-linq"
)

func TestGoLink(t *testing.T) {
	s := linq.From([]string{
		"AA",
		"BB",
	})
	fmt.Println(s.Last())
}
