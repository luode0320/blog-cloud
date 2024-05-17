package util

import (
	"os"
	"path/filepath"
)

// CreateDir 函数用于创建指定路径的目录。 如果目录已经存在，则不会报错，直接返回nil。
// 参数 dirPath 表示要创建的目录路径, 可以是绝对路径, 也可以是相对路径(程序启动目录)
// 返回一个错误，如果创建目录失败，则返回相应的错误信息。
func CreateDir(dirPath ...string) error {
	err := os.MkdirAll(filepath.Join(dirPath...), 0755)
	if err != nil {
		log.Errorf("创建目录失败: {%s}", err)
		return err
	}

	return nil
}

// RenameDir 函数用于重命名目录，如果目录不存在则不执行重命名操作。
// 参数 oldPath 表示原目录路径, 可以是绝对路径, 也可以是相对路径(程序启动目录)
// 参数 newPath 表示新目录路径, 可以是绝对路径, 也可以是相对路径(程序启动目录)
// 返回可能的错误
func RenameDir(oldPath string, newPath string) error {
	if !IsDirExist(oldPath) {
		// 原始目录不存在，不执行重命名操作
		return nil
	}

	err := os.Rename(oldPath, newPath)
	if err != nil {
		log.Errorf("修改目录名失败: {%s}", err)
		return err
	}

	return nil
}

// RemoveDir 函数用于删除指定目录及其所有子目录和文件
// 参数 dirPath 表示要删除的目录路径
func RemoveDir(dirPath ...string) {
	err := os.RemoveAll(filepath.Join(dirPath...))
	if err != nil {
		log.Errorf("删除目录失败: {%s}", err)
	}
}

// IsDirExist 函数用于检查目录是否存在
// 参数 dirPath 表示要检查的目录路径, 可以是绝对路径, 也可以是相对路径(程序启动目录)
// 返回目录是否存在的布尔值
func IsDirExist(dirPath string) bool {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}

	return true
}
