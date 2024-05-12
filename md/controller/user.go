package controller

import (
	"md/middleware"
	"md/model/common"
	"md/model/entity"
	"md/service"

	"github.com/kataras/iris/v12"
)

// 更新用户密码
func UserUpdatePassword(ctx iris.Context) {
	userCondition := entity.UserCondition{}
	resolveParam(ctx, &userCondition)
	userCondition.Id = middleware.CurrentUserId(ctx)
	service.UserUpdatePassword(userCondition)
	ctx.JSON(common.NewSuccess("更新成功"))
}
