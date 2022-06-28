package cdefs

import "time"

const (
	Nanosecond  = time.Nanosecond
	Microsecond = time.Microsecond
	Millisecond = time.Millisecond
	Second      = time.Second
	Minute      = time.Minute
	Hour        = time.Hour
	Day         = Hour * 24
	Week        = Day * 7
)

const (
	DayFormat  = "2006-01-02"
	TimeFormat = "2006-01-02 15:04:05"
)
