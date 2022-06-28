package timex

import "time"

// 当月0点
func (t *timeX) MonthFirst() time.Time {
	year, month, _ := t.t.Date()
	return time.Date(year, month, one, zero, zero, zero, zero, t.t.Location())
}
func (t *timeX) MonthFirstUnix() int64 { return t.MonthFirst().Unix() }

// 下个月1号0点
func (t *timeX) NextMonth() time.Time {
	year, month, _ := t.t.Date()
	return time.Date(year, month+one, one, zero, zero, zero, zero, t.t.Location())
}
func (t *timeX) NextMonthUnix() int64 { return t.NextMonth().Unix() }

// 当月最后1秒
func (t *timeX) MonthLast() time.Time {
	return t.NextMonth().Add(-1 * time.Second)
}
func (t *timeX) MonthLastUnix() int64 { return t.MonthLast().Unix() }

func (t *timeX) MonthUseTime() int64  { return t.Unix() - t.MonthFirstUnix() } // 当月使用时间
func (t *timeX) MonthLeftTime() int64 { return t.NextMonthUnix() - t.Unix() }  // 当月剩余时间
