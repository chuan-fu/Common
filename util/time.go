package util

import "time"

const (
	DefaultDayFormat  = "2006-01-02"
	DefaultTimeFormat = "2006-01-02 15:04:05"
	OneDay            = 60 * 60 * 24
)

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
	return time.Now().Unix()/OneDay*OneDay - 8*3600
	// times, _ := time.Parse(DefaultDayFormat, time.Now().Format(DefaultDayFormat))
	// return times.Unix() - 3600*8
}

// 获取明天0点的时间戳
func GetTomorrow() int64 {
	return time.Now().Unix()/OneDay*OneDay + 16*3600
	// times, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	// return times.Unix() + 3600*16
}

func GetTodayFormat() string {
	return time.Now().Format(DefaultDayFormat)
}

func GetNowFormat() string {
	return time.Now().Format(DefaultTimeFormat)
}
