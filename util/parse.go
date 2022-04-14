package util

import (
	"strconv"

	"github.com/pkg/errors"
)

const (
	Int64Base  = 10
	NumBitSize = 64
)

func ToInt(s string) (int, error) {
	if s == "" {
		return 0, errors.New(`ToInt: unable to ToInt("")`)
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.Wrap(err, "ToInt")
	}
	return i, nil
}

func ToInt64(s string) (int64, error) {
	if s == "" {
		return 0, errors.New(`ToInt64: unable to ToInt64("")`)
	}
	i, err := strconv.ParseInt(s, Int64Base, NumBitSize)
	if err != nil {
		return 0, errors.Wrap(err, "ToInt64")
	}
	return i, nil
}

func ToFloat(s string) (float64, error) {
	if s == "" {
		return 0, errors.New(`ToFloat: unable to ToFloat("")`)
	}
	i, err := strconv.ParseFloat(s, NumBitSize)
	if err != nil {
		return 0, errors.Wrap(err, "ToFloat")
	}
	return i, nil
}
