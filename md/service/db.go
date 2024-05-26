package service

import (
	"md/middleware"
	"md/model/common"
	"md/model/entity"
	"md/util"
	"os"
	"path/filepath"
	"strings"
)

// RefreshDb 刷新数据库数据
func RefreshDb() {
	path := filepath.Join(common.DataPath, common.ResourceName)

	// 定义一个回调函数，用于处理Walk过程中遇到的每一个文件或目录
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // 如果遇到错误，直接返回错误
		}

		split := strings.Split(path, "\\")
		if len(split) <= 1 {
			return nil
		}
		if split[1] == "picture" || split[1] == "thumbnail" {
			return nil
		}
		// 刷新一级目录
		if len(split) == 2 {
			if info.IsDir() {
				RefreshBook(entity.Book{
					ParentId: "",
					Name:     info.Name(),
				}, "")
			}
			return nil
		}
		// 刷新二级目录+文件
		if len(split) == 3 {
			if info.IsDir() {
				RefreshBook(entity.Book{
					Name: info.Name(),
				}, split[1])
			} else {
				RefreshDocument(entity.Document{
					Name:    strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())),
					Content: util.ReadFileContent(path),
				}, split[1])
			}
			return nil
		}
		// 刷新文件
		if len(split) == 4 {
			if !info.IsDir() {
				RefreshDocument(entity.Document{
					Name:    strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())),
					Content: util.ReadFileContent(path),
				}, split[2])
			}
			return nil
		}
		return nil
	})
	if err != nil {
		middleware.Log.Fatal(err)
	}
}
