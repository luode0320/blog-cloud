package middleware

import (
	"os"
)

// 初始化数据目录
func InitDataDir(dataPath, resourceName, pictureName, thumbnailName string) error {
	path := dataPath + resourceName

	picturePath := path + "/" + pictureName
	err := os.MkdirAll(picturePath, 0777)
	if err != nil {
		Log.Error("创建图片目录失败：", err)
		return err
	}

	Log.Infof("创建图片目录：{%s}", picturePath)

	thumbnailPath := path + "/" + thumbnailName
	err = os.MkdirAll(thumbnailPath, 0777)
	if err != nil {
		Log.Error("创建缩略图目录失败：", err)
		return err
	}

	Log.Infof("创建缩略图目录：{%s}", thumbnailPath)
	return nil
}
