package timex

import "time"

// 今年1月1号0点
func (t *TimeX) BeginningOfYear() time.Time {
	return time.Date(t.t.Year(), one, one, zero, zero, zero, zero, t.t.Location())
}
func (t *TimeX) BeginningOfYearUnix() int64 { return t.BeginningOfYear().Unix() }

// 明年1月1号0点
func (t *TimeX) NextYear() time.Time {
	return time.Date(t.t.Year()+one, one, one, zero, zero, zero, zero, t.t.Location())
}
func (t *TimeX) NextYearUnix() int64 { return t.NextYear().Unix() }

// 今年最后1秒
func (t *TimeX) EndOfYear() time.Time {
	return t.NextYear().Add(-1 * time.Second)
}
func (t *TimeX) EndOfYearUnix() int64 { return t.EndOfYear().Unix() }

func (t *TimeX) YearUseTime() int64  { return t.Unix() - t.BeginningOfYearUnix() } // 今年使用时间
func (t *TimeX) YearLeftTime() int64 { return t.NextYearUnix() - t.Unix() }        // 今年剩余时间
