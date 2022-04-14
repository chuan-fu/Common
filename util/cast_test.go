package util

import (
	"fmt"
	"testing"
	"time"
)

func TestToString(t *testing.T) {
	fmt.Println(ToString(nil))
	s := 1
	bs := []byte("aa")
	fmt.Println(ToString(s))
	fmt.Println(ToString(&bs))
	fmt.Println(ToString(&s))
	fmt.Println(ToString(&map[string]string{
		"1": "2",
	}))
	fmt.Println(ToString(1))
	fmt.Println(ToString(1.2))
	fmt.Println(ToString(-1.2))
}

func TestTimeDurationToString(t *testing.T) {
	fmt.Println(timeDurationToString(2 * time.Hour))
	fmt.Println(timeDurationToString(2 * time.Minute))
	fmt.Println(timeDurationToString(120 * time.Second))
	fmt.Println(timeDurationToString(2 * time.Second))
	fmt.Println(timeDurationToString(2 * time.Millisecond))
	fmt.Println(timeDurationToString(2 * time.Microsecond))
	fmt.Println(timeDurationToString(2 * time.Nanosecond))
}

func TestToInt(t *testing.T) {
	fmt.Println(ToIntI("a"))
	fmt.Println(ToIntI("1"))
	fmt.Println(ToIntI(true))
	fmt.Println(ToIntI(false))
	fmt.Println(ToIntI(1.23))
}
