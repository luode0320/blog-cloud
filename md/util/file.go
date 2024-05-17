package util

import (
	"md/middleware"
	"os"
	"path/filepath"
)

// CreateFile 函数用于创建文件并写入内容
// 参数 dirPath 表示文件所在目录路径
// 参数 fileName 表示要创建的文件名
// 参数 content 表示要写入的文件内容
// 返回可能的错误
func CreateFile(dirPath string, fileName string, content string) error {
	if err := CreateDir(dirPath); err != nil {
		return err
	}

	saveMdFile, err := os.Create(filepath.Join(dirPath, fileName))
	if err != nil {
		log.Errorf("创建文件失败: {%s}", err)
		return err
	}
	defer saveMdFile.Close()

	_, err = saveMdFile.Write([]byte(content))
	if err != nil {
		log.Errorf("写入文件错误: {%s}", err)
		return err
	}

	return nil
}

// RenameFile 函数用于重命名文件
// 参数 dirPath 表示文件所在目录路径
// 参数 oldFileName 表示原文件名
// 参数 newFileName 表示新文件名
// 返回可能的错误
func RenameFile(dirPath string, oldFileName string, newFileName string) error {
	// 创建目录
	if err := CreateDir(dirPath); err != nil {
		return err
	}

	oldFile := filepath.Join(dirPath, oldFileName)
	newFile := filepath.Join(dirPath, newFileName)
	err := os.Rename(oldFile, newFile)
	if err != nil {
		log.Errorf("重命名文件失败: {%s}", err)
		return err
	}

	return nil
}

// RemoveFile 函数用于删除文件
// 参数 filePath 表示文件所在目录路径
// 参数 fileName 表示要删除的文件名
// 返回可能的错误
func RemoveFile(filePath string, fileName string) {
	err := os.Remove(filepath.Join(filePath, fileName))
	if err != nil {
		middleware.Log.Errorf("删除文件失败: {%s}", err)
	}
}
