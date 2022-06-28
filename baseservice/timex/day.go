package timex

import (
	"time"

	"github.com/chuan-fu/Common/cdefs"
)

// 当日0点
func (t *TimeX) BeginningOfDay() time.Time {
	y, m, d := t.t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.t.Location())
}
func (t *TimeX) BeginningOfDayUnix() int64 { return t.BeginningOfDay().Unix() }

// 当日最后1秒
func (t *TimeX) EndOfDay() time.Time {
	return t.BeginningOfDay().Add(cdefs.Day - time.Second)
}
func (t *TimeX) EndOfDayUnix() int64 { return t.EndOfDay().Unix() }

// 明日0点
func (t *TimeX) NextDay() time.Time {
	return t.BeginningOfDay().Add(cdefs.Day)
}
func (t *TimeX) NextDayUnix() int64 { return t.NextDay().Unix() }

func (t *TimeX) DayUseTime() int64  { return t.Unix() - t.BeginningOfDayUnix() } // 当日使用时间
func (t *TimeX) DayLeftTime() int64 { return t.NextDayUnix() - t.Unix() }        // 当日剩余时间
