package controller

import (
	"md/model/common"
	"md/model/entity"
	"md/service"

	"github.com/kataras/iris/v12"
)

// 注册
func SignUp(ctx iris.Context) {
	user := entity.User{}
	resolveParam(ctx, &user)
	service.SignUp(user)
	ctx.JSON(common.NewSuccess("注册成功"))
}

// 登录
func SignIn(ctx iris.Context) {
	user := entity.User{}
	resolveParam(ctx, &user)
	tokenResult := service.SignIn(user)
	ctx.JSON(common.NewSuccessData("登录成功", tokenResult))
}

// 退出登录
func SignOut(ctx iris.Context) {
	tokenResult := common.TokenResult{}
	resolveParam(ctx, &tokenResult)
	service.SignOut(tokenResult)
	ctx.JSON(common.NewSuccess("已退出登录"))
}

// 刷新token
func TokenRefresh(ctx iris.Context) {
	tokenResult := common.TokenResult{}
	resolveParam(ctx, &tokenResult)
	tokenResult = service.TokenRefresh(tokenResult.RefreshToken)
	ctx.JSON(common.NewSuccessData("token刷新成功", tokenResult))
}
