package middleware

import (
	"bytes"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"net/http"
	"strings"
)

// 创建一个自定义的中间件函数，用于打印 HTTP 请求日志和参数
func RequestLogger(ctx iris.Context) {
	method := ctx.Request().Method
	path := ctx.Request().URL.String()

	// 如果是POST/PUT请求，并且内容类型为JSON，则读取内容体
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
	}
	params := ""
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err == nil {
		defer ctx.Request().Body.Close()

		buf := bytes.NewBuffer(body)
		ctx.Request().Body = ioutil.NopCloser(buf)
		params = string(body)
		if strings.Contains(params, "\r\n") {
			params = strings.ReplaceAll(params, "\r\n", "")
		}
		if strings.Contains(params, "\n") {
			params = strings.ReplaceAll(params, "\n", "")
		}
		params = strings.ReplaceAll(params, " ", "")
	}
	ctx.Next()

	divider := "────────────────────────────────────────────────────────────────────────────────────────────────────────────────"
	topDivider := "┌" + divider
	middleDivider := "├" + divider
	bottomDivider := "└" + divider
	outputStr :=
		"\n" + topDivider +
			"\n│ 请求地址:" + path + "\n" + middleDivider +
			"\n│ 请求参数:" + params +
			"\n" + bottomDivider
	Log.Infof(outputStr)

	// 继续处理请求
	ctx.Next()
}
