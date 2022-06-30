package timex

import (
	"fmt"
	"time"

	"github.com/chuan-fu/Common/cdefs"
)

type TimeX struct {
	t time.Time
}

func NewNow() *TimeX {
	return &TimeX{t: time.Now()}
}

// 秒级时间戳
func NewTimestamp(unixSec int64) *TimeX {
	return &TimeX{t: time.Unix(unixSec, 0)}
}

func NewTimestampWithLocation(unixSec int64, loc *time.Location) *TimeX {
	return &TimeX{t: time.Unix(unixSec, 0).In(loc)}
}

func NewTimeX(t time.Time) *TimeX {
	return &TimeX{t: t}
}

func (t *TimeX) Time() time.Time {
	return t.t
}

func (t *TimeX) Unix() int64 {
	return t.t.Unix()
}

func (t *TimeX) UnixNano() int64 {
	return t.t.UnixNano()
}

func (t *TimeX) Format(layout string) string {
	return t.t.Format(layout)
}

func (t *TimeX) FormatToday() string {
	return t.t.Format(cdefs.DayFormat)
}

func (t *TimeX) FormatTime() string {
	return t.t.Format(cdefs.TimeFormat)
}

// 添加时间，可为负数
func (t *TimeX) Add(d time.Duration) *TimeX {
	return &TimeX{t: t.t.Add(d)}
}

// 添加一天
func (t *TimeX) AddDay() *TimeX {
	return &TimeX{t: t.t.Add(cdefs.Day)}
}

// 添加几天
func (t *TimeX) AddSomeDay(d int64) *TimeX {
	return &TimeX{t: t.t.Add(time.Duration(d) * cdefs.Day)}
}

// 添加一周
func (t *TimeX) AddWeek() *TimeX {
	return &TimeX{t: t.t.Add(cdefs.Day)}
}

// 添加时间
func (t *TimeX) AddDate(years, months, days int) *TimeX {
	return &TimeX{t: t.t.AddDate(years, months, days)}
}

// t是否比u大，true表示t大
func (t *TimeX) After(u time.Time) bool {
	return t.t.After(u)
}

// t是否比u小，true表示t小
func (t *TimeX) Before(u time.Time) bool {
	return t.t.Before(u)
}

// 是否是同一时间 存在时区差异 6:00 +0200 和 4:00 UTC 是相同的
func (t *TimeX) Same(u time.Time) bool {
	return t.t.Equal(u)
}

// 是否是空时间 time.Time{} 表示公元1年1月1日0点
func (t *TimeX) Zero() bool {
	return t.t.IsZero()
}

// 计算时间差 t-u
func (t *TimeX) Sub(u time.Time) time.Duration {
	return t.t.Sub(u)
}

func (t *TimeX) SubPrint(u time.Time) {
	fmt.Print(t.t.Sub(u))
}

func (t *TimeX) SubPrintln(u time.Time) {
	fmt.Println(t.t.Sub(u))
}

func Now() int64 {
	return time.Now().Unix()
}

func TodayFormat() string {
	return time.Now().Format(cdefs.DayFormat)
}

func TimeFormat() string {
	return time.Now().Format(cdefs.TimeFormat)
}

// 计算耗时 now() > t
// 等同于 time.Now().Sub(t)
func Since(t time.Time) time.Duration {
	return time.Since(t)
}

func SincePrint(t time.Time) {
	fmt.Print(time.Since(t))
}
