// 随机字符串工具类
package util

import (
	"math/rand"
	"strings"

	uuid "github.com/iris-contrib/go.uuid"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	numbers = []rune("0123456789")
)

// 生成uuid
func Uuid() string {
	result, err := uuid.NewV4()
	if err != nil {
		return RandomString(32)
	}
	return result.String()
}

// 生成没有连字符的uuid
func UuidNoHyphen() string {
	result, err := uuid.NewV4()
	if err != nil {
		return RandomString(32)
	}
	return strings.ReplaceAll(result.String(), "-", "")
}

// 随机生成字符串
func RandomString(length int) string {
	return random(length, letters)
}

// 随机生成数字字符串
func RandomNumber(length int) string {
	return random(length, numbers)
}

// 从对应字符集中随机生成指定长度字符串
func random(length int, arr []rune) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = arr[rand.Intn(len(arr))]
	}
	return string(b)
}
