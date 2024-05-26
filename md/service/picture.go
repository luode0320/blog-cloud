package service

import (
	"io"
	"md/dao"
	"md/middleware"
	"md/model/common"
	"md/model/entity"
	"md/util"
	"mime/multipart"
	"path/filepath"
	"slices"
	"time"
)

// 分页查询图片记录
func PicturePage(pageCondition common.PageCondition[interface{}], userId string) common.PageResult[entity.PicturePageResult] {
	pictures, total, err := dao.PicturePage(middleware.Db, pageCondition.Page, userId)
	if err != nil {
		panic(common.NewErr("查询失败", err))
	}

	picturePageResults := []entity.PicturePageResult{}
	for _, v := range pictures {
		picturePageResults = append(picturePageResults, entity.PicturePageResult{
			Picture:         v,
			PicturePrefix:   "/" + filepath.ToSlash(filepath.Join(common.PictureName)) + "/",
			ThumbnailPrefix: "/" + filepath.ToSlash(filepath.Join(common.ThumbnailName)) + "/",
		})
	}

	pageResult := common.PageResult[entity.PicturePageResult]{
		Records: picturePageResults,
		Total:   total,
	}
	return pageResult
}

// 删除图片
func PictureDelete(id, userId string) {
	tx := middleware.DbW.MustBegin()
	defer tx.Rollback()

	// 查询图片
	picture, err := dao.PictureGetById(tx, id, userId)
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}

	// 查询相同大小、相同hash的图片数量
	countResult, err := dao.PictureCountBySizeHash(tx, picture.Size, picture.Hash)
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}

	// 删除记录
	err = dao.PictureDeleteById(tx, id, userId)
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}

	// 如果相同的图片只有一条记录，删除文件
	if countResult.Count == 1 {
		util.RemoveFile(filepath.Join(common.DataPath, common.ResourceName, common.PictureName), picture.Path)
		util.RemoveFile(filepath.Join(common.DataPath, common.ResourceName, common.ThumbnailName), picture.Path)
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}

	middleware.Log.Infof("成功删除图片: {%s}", id)
}

// 图片上传
func PictureUpload(pictureFile, thumbnailFile multipart.File, pictureInfo, thumbnailInfo *multipart.FileHeader, userId string) (string, string) {
	// 校验文件大小
	if pictureInfo.Size == 0 {
		panic(common.NewError("图片解析失败"))
	}
	if thumbnailInfo.Size == 0 {
		panic(common.NewError("缩略图解析失败"))
	}
	if pictureInfo.Size > 1000*1000*20 {
		panic(common.NewError("图片大小不可超过20MB"))
	}
	if thumbnailInfo.Size > 1000*100 {
		panic(common.NewError("缩略图大小不可超过100KB"))
	}

	// 获取图片,略缩图后缀
	pictureExt := util.FileExt(pictureInfo.Filename)
	thumbnailExt := util.FileExt(thumbnailInfo.Filename)

	extArr := []string{".apng", ".bmp", ".gif", ".ico", ".jfif", ".jpeg", ".jpg", ".png", ".webp"}
	extNames := "APNG,BMP,GIF,ICO,JFIF,JPEG,JPG,PNG,WebP"

	if !slices.Contains(extArr, pictureExt) {
		panic(common.NewError("仅支持以下格式的图片：" + extNames))
	}

	if !slices.Contains(extArr, thumbnailExt) {
		panic(common.NewError("仅支持以下格式的缩略图：" + extNames))
	}

	if util.StringLength(pictureInfo.Filename) > 1000 || util.StringLength(thumbnailInfo.Filename) > 1000 {
		panic(common.NewError("图片文件名称过长"))
	}

	// 读取文件
	pictureByte, err := io.ReadAll(pictureFile)
	if err != nil {
		panic(common.NewErr("图片解析失败", err))
	}

	thumbnailByte, err := io.ReadAll(thumbnailFile)
	if err != nil {
		panic(common.NewErr("缩略图解析失败", err))
	}

	// 生成sha256校验码
	sha256Str := util.EncryptSHA256(pictureByte)

	// 生成文件名
	filename := util.SnowflakeString() + pictureExt

	// 保存文件
	dirPath := filepath.Join(common.DataPath, common.ResourceName, common.PictureName)
	if err := util.CreateFile(dirPath, filename, pictureByte); err != nil {
		panic(common.NewErr("图片上传失败", err))
	}

	// 保存缩略图
	dirPath = filepath.Join(common.DataPath, common.ResourceName, common.ThumbnailName)
	if err := util.CreateFile(dirPath, filename, thumbnailByte); err != nil {
		panic(common.NewErr("图片上传失败", err))
	}

	// 添加记录
	tx := middleware.DbW.MustBegin()
	defer tx.Rollback()

	picture := entity.Picture{}
	picture.Id = util.SnowflakeString()
	picture.CreateTime = time.Now().UnixMilli()
	picture.Name = pictureInfo.Filename
	picture.Path = filename
	picture.Hash = sha256Str
	picture.Size = pictureInfo.Size
	picture.UserId = userId
	err = dao.PictureAdd(tx, picture)
	if err != nil {
		panic(common.NewErr("图片上传失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("图片上传失败", err))
	}

	path := "/" + filepath.ToSlash(filepath.Join(common.PictureName, filename))
	middleware.Log.Infof("成功上传图片: {%s}", path)
	return path, "上传成功"
}
