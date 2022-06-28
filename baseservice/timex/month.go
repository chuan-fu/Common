package timex

import (
	"time"
)

// 当月0点
func (t *TimeX) BeginningOfMonth() time.Time {
	y, m, _ := t.t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.t.Location())
}
func (t *TimeX) BeginningOfMonthUnix() int64 { return t.BeginningOfMonth().Unix() }

// 下个月1号0点
func (t *TimeX) NextMonth() time.Time {
	y, m, _ := t.t.Date()
	return time.Date(y, m+1, 1, 0, 0, 0, 0, t.t.Location())
}
func (t *TimeX) NextMonthUnix() int64 { return t.NextMonth().Unix() }

// 当月最后1秒
func (t *TimeX) EndOfMonth() time.Time {
	return t.NextMonth().Add(-1 * time.Second)
}
func (t *TimeX) EndOfMonthUnix() int64 { return t.EndOfMonth().Unix() }

func (t *TimeX) MonthUseTime() int64  { return t.Unix() - t.BeginningOfMonthUnix() } // 当月使用时间
func (t *TimeX) MonthLeftTime() int64 { return t.NextMonthUnix() - t.Unix() }        // 当月剩余时间
