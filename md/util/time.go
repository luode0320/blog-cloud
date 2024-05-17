package util

import "time"

// CreateStamp 函数用于生成当前时间的 Unix 毫秒时间戳, 它是一个13位的数字
// 返回当前时间的 Unix 毫秒时间戳
func CreateStamp() int64 {
	return time.Now().UnixMilli()
}
