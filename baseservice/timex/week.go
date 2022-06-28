package timex

import "time"

// 本周周一0点
func (t *timeX) WeekFirst() time.Time {
	year, month, day := t.t.Date()
	t.t.Weekday()
	return time.Date(year, month, day-int(t.t.Weekday())+one, zero, zero, zero, zero, t.t.Location())
}
func (t *timeX) WeekFirstUnix() int64 { return t.WeekFirst().Unix() }

// 下周一0点
func (t *timeX) NextWeek() time.Time {
	return t.WeekFirst().Add(OneWeek)
}
func (t *timeX) NextWeekUnix() int64 { return t.NextWeek().Unix() }

// 本周最后1秒
func (t *timeX) WeekLast() time.Time {
	return t.NextWeek().Add(-1 * time.Second)
}
func (t *timeX) WeekLastUnix() int64 { return t.WeekLast().Unix() }

func (t *timeX) WeekUseTime() int64  { return t.Unix() - t.WeekFirstUnix() } // 本周使用时间
func (t *timeX) WeekLeftTime() int64 { return t.NextWeekUnix() - t.Unix() }  // 本周剩余时间
