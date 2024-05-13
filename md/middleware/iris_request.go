package middleware

import (
	"github.com/kataras/iris/v12"
)

// 创建一个自定义的中间件函数，用于打印 HTTP 请求日志和参数
func RequestLogger(ctx iris.Context) {
	// 打印请求方法和路径
	Log.Infof("http请求: %s %s", ctx.Request().Method, ctx.Request().RequestURI)

	// 打印请求参数
	params := ctx.Request().URL.Query()
	if len(params) > 0 {
		Log.Infof("请求参数:")
		for key, values := range params {
			for _, value := range values {
				Log.Infof("%s: %s", key, value)
			}
		}
	}

	// 继续处理请求
	ctx.Next()
}
