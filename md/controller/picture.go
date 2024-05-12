package controller

import (
	"md/middleware"
	"md/model/common"
	"md/model/entity"
	"md/service"

	"github.com/kataras/iris/v12"
)

// 分页查询图片记录
func PicturePage(ctx iris.Context) {
	pageCondition := common.PageCondition[interface{}]{}
	resolveParam(ctx, &pageCondition)
	userId := middleware.CurrentUserId(ctx)
	ctx.JSON(common.NewSuccessData("查询成功", service.PicturePage(pageCondition, userId)))
}

// 删除图片
func PictureDelete(ctx iris.Context) {
	picture := entity.Picture{}
	resolveParam(ctx, &picture)
	userId := middleware.CurrentUserId(ctx)
	service.PictureDelete(picture.Id, userId)
	ctx.JSON(common.NewSuccess("删除成功"))
}

// 上传图片
func PictureUpload(ctx iris.Context) {
	userId := middleware.CurrentUserId(ctx)
	pictureFile, pictureInfo, err := ctx.FormFile("picture")
	if err != nil {
		panic(common.NewErr("图片解析失败", err))
	}
	defer pictureFile.Close()
	thumbnailFile, thumbnailInfo, err := ctx.FormFile("thumbnail")
	if err != nil {
		panic(common.NewErr("图片解析失败", err))
	}
	defer thumbnailFile.Close()
	path, message := service.PictureUpload(pictureFile, thumbnailFile, pictureInfo, thumbnailInfo, userId)
	ctx.JSON(common.NewSuccessData(message, path))
}
