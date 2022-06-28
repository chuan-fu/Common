package timex

import (
	"fmt"
	"time"
)

const (
	zero              = 0
	one               = 1
	DefaultDayFormat  = "2006-01-02"
	DefaultTimeFormat = "2006-01-02 15:04:05"
	OneDay            = 24 * time.Hour
	OneWeek           = 7 * OneDay
)

type TimeX interface {
	Unix() int64
	UnixNano() int64
	Time() time.Time

	Format(layout string) string
	FormatDay() string
	FormatTime() string

	DayFirst() time.Time // 当日0点
	DayFirstUnix() int64 // 当日0点 时间戳
	DayLast() time.Time  // 当日最后1秒 23:59:59
	DayLastUnix() int64  // 当日最后1秒 时间戳
	NextDay() time.Time  // 明日0点
	NextDayUnix() int64  // 明日0点 时间戳
	DayUseTime() int64   // 今日使用时间
	DayLeftTime() int64  // 今日剩余时间

	WeekFirst() time.Time // 本周周一0点
	WeekFirstUnix() int64 // 本周周一0点 时间戳
	NextWeek() time.Time  // 下周一0点
	NextWeekUnix() int64  // 下周一0点 时间戳
	WeekLast() time.Time  // 本周最后1秒
	WeekLastUnix() int64  // 本周最后1秒 时间戳
	WeekUseTime() int64   // 本周使用时间
	WeekLeftTime() int64  // 本周剩余时间

	MonthFirst() time.Time // 当月0点
	MonthFirstUnix() int64 // 当月0点 时间戳
	NextMonth() time.Time  // 下月1号
	NextMonthUnix() int64  // 下月1号 时间戳
	MonthLast() time.Time  // 当月最后1秒
	MonthLastUnix() int64  // 当月最后1秒 时间戳
	MonthUseTime() int64   // 当月使用时间
	MonthLeftTime() int64  // 当月剩余时间

	Add(d time.Duration) TimeX             // 添加时间，可为负数
	AddDate(years, months, days int) TimeX // 添加时间
	After(u time.Time) bool                // t是否比u大，true表示t大
	Before(u time.Time) bool               // t是否比u小，true表示t小
	Is(u time.Time) bool                   // 是否是同一时间 存在时区差异 6:00 +0200 和 4:00 UTC 是相同的
	IsZero() bool                          // 是否是空时间 time.Time{} 表示公元0001年1月1日0点

	Sub(u time.Time) time.Duration // 计算时间差 t-u
	SubPrint(u time.Time)          // 输出时间差
	SubPrintln(u time.Time)        // 输出时间差
}

type timeX struct {
	t time.Time
}

func NewNow() TimeX {
	return &timeX{t: time.Now()}
}

// 秒级时间戳
func NewTimestamp(unixSec int64) TimeX {
	return &timeX{t: time.Unix(unixSec, zero).In(time.Local)}
}

func NewTimeX(t time.Time) TimeX {
	return &timeX{t: t}
}

func (t *timeX) Time() time.Time {
	return t.t
}

func (t *timeX) Unix() int64 {
	return t.t.Unix()
}

func (t *timeX) UnixNano() int64 {
	return t.t.UnixNano()
}

func (t *timeX) Format(layout string) string {
	return t.t.Format(layout)
}

func (t *timeX) FormatDay() string {
	return t.t.Format(DefaultDayFormat)
}

func (t *timeX) FormatTime() string {
	return t.t.Format(DefaultTimeFormat)
}

// 添加时间，可为负数
func (t *timeX) Add(d time.Duration) TimeX {
	return &timeX{t: t.t.Add(d)}
}

// 添加时间
func (t *timeX) AddDate(years, months, days int) TimeX {
	return &timeX{t: t.t.AddDate(years, months, days)}
}

// t是否比u大，true表示t大
func (t *timeX) After(u time.Time) bool {
	return t.t.After(u)
}

// t是否比u小，true表示t小
func (t *timeX) Before(u time.Time) bool {
	return t.t.Before(u)
}

// 是否是同一时间 存在时区差异 6:00 +0200 和 4:00 UTC 是相同的
func (t *timeX) Is(u time.Time) bool {
	return t.t.Equal(u)
}

// 是否是空时间 time.Time{} 表示公元1年1月1日0点
func (t *timeX) IsZero() bool {
	return t.t.IsZero()
}

// 计算时间差 t-u
func (t *timeX) Sub(u time.Time) time.Duration {
	return t.t.Sub(u)
}

func (t *timeX) SubPrint(u time.Time) {
	fmt.Print(t.t.Sub(u))
}

func (t *timeX) SubPrintln(u time.Time) {
	fmt.Println(t.t.Sub(u))
}

func Now() int64 {
	return time.Now().Unix()
}

// 计算耗时 now() > t
// 等同于 time.Now().Sub(t)
func Since(t time.Time) time.Duration {
	return time.Since(t)
}

func SincePrint(t time.Time) {
	fmt.Print(time.Since(t))
}
