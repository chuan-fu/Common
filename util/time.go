package util

import "time"

const (
	DefaultDayFormat  = "2006-01-02"
	DefaultTimeFormat = "2006-01-02 15:04:05"
	OneDay            = 60 * 60 * 24
)

// 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// TimeWithUnixSec 将时间戳转为time, 时区为本地
func TimeWithUnixSec(unixSec int64) time.Time {
	return time.Unix(unixSec, 0).In(time.Local)
}

// 获取参数的时间与今天剩余时间的较少值
func GetTodayMinTime(t time.Duration) time.Duration {
	t2 := GetTodayLastSecond()
	if t < t2 {
		return t
	}
	return t2
}

// 获取今天剩余时间
func GetTodayLastSecond() time.Duration {
	now := time.Now().Unix()
	return time.Duration(GetTomorrow()-now) * time.Second
}

// 获取今天0点的时间戳
func GetToday() int64 {
	return GetZeroTime(time.Now()).Unix()
}

// 获取明天0点的时间戳
func GetTomorrow() int64 {
	return GetZeroTime(time.Now()).Unix() + OneDay
}

// 获取那天0点的时间戳
func GetThatDay(t int64) int64 {
	return GetZeroTime(TimeWithUnixSec(t)).Unix()
}

func GetTodayFormat() string {
	return time.Now().Format(DefaultDayFormat)
}

func GetNowFormat() string {
	return time.Now().Format(DefaultTimeFormat)
}
