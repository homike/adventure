package util

import "time"

func ByteToInt32(bytes []byte) int32 {
	return 0
}

// 两个日期，相差天数
func TimeSub(t1, t2 time.Time) int {
	t1 = t1.UTC().Truncate(24 * time.Hour)
	t2 = t2.UTC().Truncate(24 * time.Hour)
	return int(t1.Sub(t2).Hours() / 24)
}

func TimeSubNow(date int64) int32 {
	return int32(time.Now().Sub(time.Unix(date, 0)).Seconds())
}
