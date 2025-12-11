package time

import (
	"time"
)

// LocBeijing 北京时区
var LocBeijing *time.Location

var (
	// 初始化时间
	InitTime = time.Date(1970, 1, 1, 0, 0, 2, 0, time.UTC)
)

func init() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	LocBeijing = loc
}

// DayBeginUTCTime 根据指定时刻 t，返回 t 所在指定时区的当天的起始UTC时间。
func DayBeginUTCTime(t time.Time, loc *time.Location) time.Time {
	tm := t.In(loc)
	return time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, loc).UTC()
}

// DayEndUTCTime 根据指定时刻 t，返回 t 所在指定时区的的结束UTC时间。
func DayEndUTCTime(t time.Time, loc *time.Location) time.Time {
	tm := t.In(loc)
	return time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 59, 59, 0, loc).UTC()
}

//获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

//获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, LocBeijing)
}
