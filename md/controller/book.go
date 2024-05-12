package controller

import (
	"md/middleware"
	"md/model/common"
	"md/model/entity"
	"md/service"

	"github.com/kataras/iris/v12"
)

// 添加文集
func BookAdd(ctx iris.Context) {
	book := entity.Book{}
	resolveParam(ctx, &book)
	book.UserId = middleware.CurrentUserId(ctx)
	service.BookAdd(book)
	ctx.JSON(common.NewSuccess("添加成功"))
}

// 修改文集
func BookUpdate(ctx iris.Context) {
	book := entity.Book{}
	resolveParam(ctx, &book)
	book.UserId = middleware.CurrentUserId(ctx)
	service.BookUpdate(book)
	ctx.JSON(common.NewSuccess("更新成功"))
}

// 删除文集
func BookDelete(ctx iris.Context) {
	book := entity.Book{}
	resolveParam(ctx, &book)
	userId := middleware.CurrentUserId(ctx)
	service.BookDelete(book.Id, userId)
	ctx.JSON(common.NewSuccess("删除成功"))
}

// 查询文集列表
func BookList(ctx iris.Context) {
	userId := middleware.CurrentUserId(ctx)
	ctx.JSON(common.NewSuccessData("查询成功", service.BookList(userId)))
}
