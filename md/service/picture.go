package service

import (
	"io"
	"md/dao"
	"md/middleware"
	"md/model/common"
	"md/model/entity"
	"md/util"
	"mime/multipart"
	"os"
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
		picturePageResults = append(picturePageResults, entity.PicturePageResult{Picture: v, PicturePrefix: "/" + common.ResourceName + "/" + common.PictureName + "/", ThumbnailPrefix: "/" + common.ResourceName + "/" + common.ThumbnailName + "/"})
	}
	pageResult := common.PageResult[entity.PicturePageResult]{Records: picturePageResults, Total: total}
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
		os.Remove(common.DataPath + common.ResourceName + "/" + common.PictureName + "/" + picture.Path)
		os.Remove(common.DataPath + common.ResourceName + "/" + common.ThumbnailName + "/" + picture.Path)
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}
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

	// 获取文件后缀
	pictureExt := util.FileExt(pictureInfo.Filename)
	thumbnailExt := util.FileExt(thumbnailInfo.Filename)
	extArr := []string{".apng", ".bmp", ".gif", ".ico", ".jfif", ".jpeg", ".jpg", ".png", ".webp"}
	extNames := "APNG、BMP、GIF、ICO、JPEG、PNG、WebP"
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

	// 查询相同大小和校验码的文件
	pictures, err := dao.PictureBySizeHash(middleware.Db, pictureInfo.Size, sha256Str)
	if err != nil {
		panic(common.NewErr("图片上传失败", err))
	}

	// 生成文件名
	filename := util.SnowflakeString() + pictureExt

	needAddRecord := true
	message := "上传成功"
	if len(pictures) == 0 {
		// 无相同文件，保存文件
		savePictureFile, err := os.Create(common.DataPath + common.ResourceName + "/" + common.PictureName + "/" + filename)
		if err != nil {
			panic(common.NewErr("图片上传失败", err))
		}
		defer savePictureFile.Close()
		_, err = savePictureFile.Write(pictureByte)
		if err != nil {
			panic(common.NewErr("图片上传失败", err))
		}
		// 保存缩略图
		saveThumbnailFile, err := os.Create(common.DataPath + common.ResourceName + "/" + common.ThumbnailName + "/" + filename)
		if err != nil {
			panic(common.NewErr("图片上传失败", err))
		}
		defer saveThumbnailFile.Close()
		_, err = saveThumbnailFile.Write(thumbnailByte)
		if err != nil {
			panic(common.NewErr("图片上传失败", err))
		}
	} else {
		// 存在相同文件，文件名使用原记录中的名字
		for _, v := range pictures {
			filename = v.Path
			// 判断是否自己也存在相同文件，如存在则不添加记录
			if v.UserId == userId {
				needAddRecord = false
				message = "图片已存在"
				break
			}
		}
	}

	// 添加记录
	if needAddRecord {
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
	}

	return "/" + common.ResourceName + "/" + common.PictureName + "/" + filename, message
}
