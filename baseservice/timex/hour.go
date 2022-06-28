package timex

import (
	"time"
)

// 当前小时0分0秒
func (t *TimeX) BeginningOfHour() time.Time {
	y, m, d := t.t.Date()
	return time.Date(y, m, d, t.t.Hour(), 0, 0, 0, t.t.Location())
}
func (t *TimeX) BeginningOfHourUnix() int64 { return t.BeginningOfHour().Unix() }

// 当前小时59分59秒
func (t *TimeX) EndOfHour() time.Time {
	return t.BeginningOfHour().Add(time.Hour - time.Second)
}
func (t *TimeX) EndOfHourUnix() int64 { return t.EndOfHour().Unix() }

// 下一个小时0分0秒
func (t *TimeX) NextHour() time.Time {
	return t.BeginningOfHour().Add(time.Hour)
}
func (t *TimeX) NextHourUnix() int64 { return t.NextHour().Unix() }

func (t *TimeX) HourUseTime() int64  { return t.Unix() - t.BeginningOfHourUnix() } // 当前小时使用时间
func (t *TimeX) HourLeftTime() int64 { return t.NextHourUnix() - t.Unix() }        // 当前小时剩余时间
