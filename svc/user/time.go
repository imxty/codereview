package user

import "time"

// LocBeijing 北京时间
var LocBeijing *time.Location

var (
	// 初始化时间
	InitTime = time.Date(1970, 1, 2, 0, 0, 1, 0, time.UTC)
)

func init() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	LocBeijing = loc
}
