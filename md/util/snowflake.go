// 雪花算法工具类
package util

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

// 初始化雪花算法节点
func InitSnowflake(nodeCode int64) error {
	var err error
	node, err = snowflake.NewNode(nodeCode)
	return err
}

// 生成int64类型的雪花id
func SnowflakeInt() int64 {
	return node.Generate().Int64()
}

// 生成string类型的雪花id
func SnowflakeString() string {
	return node.Generate().String()
}

// 生成[]byte类型的雪花id
func SnowflakeBytes() []byte {
	return node.Generate().Bytes()
}
