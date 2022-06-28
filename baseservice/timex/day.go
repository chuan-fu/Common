package timex

import "time"

// 当日0点
func (t *timeX) DayFirst() time.Time {
	year, month, day := t.t.Date()
	return time.Date(year, month, day, zero, zero, zero, zero, t.t.Location())
}
func (t *timeX) DayFirstUnix() int64 { return t.DayFirst().Unix() }

// 当日最后1秒
func (t *timeX) DayLast() time.Time {
	return t.DayFirst().Add(OneDay - time.Second)
}
func (t *timeX) DayLastUnix() int64 { return t.DayLast().Unix() }

// 明日0点
func (t *timeX) NextDay() time.Time {
	return t.DayFirst().Add(OneDay)
}
func (t *timeX) NextDayUnix() int64 { return t.NextDay().Unix() }

func (t *timeX) DayUseTime() int64  { return t.Unix() - t.DayFirstUnix() } // 当日使用时间
func (t *timeX) DayLeftTime() int64 { return t.NextDayUnix() - t.Unix() }  // 当日剩余时间
