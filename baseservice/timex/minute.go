package timex

import (
	"time"
)

// 当前分钟0秒
func (t *TimeX) BeginningOfMinute() time.Time {
	y, m, d := t.t.Date()
	hour, min, _ := t.t.Clock()
	return time.Date(y, m, d, hour, min, 0, 0, t.t.Location())
}
func (t *TimeX) BeginningOfMinuteUnix() int64 { return t.BeginningOfMinute().Unix() }

// 当前分钟59秒
func (t *TimeX) EndOfMinute() time.Time {
	return t.BeginningOfMinute().Add(time.Minute - time.Second)
}
func (t *TimeX) EndOfMinuteUnix() int64 { return t.EndOfMinute().Unix() }

// 下一分钟0秒
func (t *TimeX) NextMinute() time.Time {
	return t.BeginningOfMinute().Add(time.Minute)
}
func (t *TimeX) NextMinuteUnix() int64 { return t.NextMinute().Unix() }

func (t *TimeX) MinuteUseTime() int64  { return t.Unix() - t.BeginningOfMinuteUnix() } // 当前分钟使用时间
func (t *TimeX) MinuteLeftTime() int64 { return t.NextMinuteUnix() - t.Unix() }        // 当前分钟剩余时间
