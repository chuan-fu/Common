package timex

import (
	"time"

	"github.com/chuan-fu/Common/cdefs"
)

// 本周周一0点
func (t *TimeX) BeginningOfWeek() time.Time {
	y, m, d := t.t.Date()
	if weekDay := t.t.Weekday(); weekDay == time.Sunday {
		d -= int(time.Saturday)
	} else {
		d = d - int(weekDay) + 1
	}
	return time.Date(y, m, d, 0, 0, 0, 0, t.t.Location())
}
func (t *TimeX) BeginningOfWeekUnix() int64 { return t.BeginningOfWeek().Unix() }

// 下周一0点
func (t *TimeX) NextWeek() time.Time {
	return t.BeginningOfWeek().Add(cdefs.Week)
}
func (t *TimeX) NextWeekUnix() int64 { return t.NextWeek().Unix() }

// 本周最后1秒
func (t *TimeX) EndOfWeek() time.Time {
	return t.NextWeek().Add(-1 * time.Second)
}
func (t *TimeX) EndOfWeekUnix() int64 { return t.EndOfWeek().Unix() }

func (t *TimeX) WeekUseTime() int64  { return t.Unix() - t.BeginningOfWeekUnix() } // 本周使用时间
func (t *TimeX) WeekLeftTime() int64 { return t.NextWeekUnix() - t.Unix() }        // 本周剩余时间
