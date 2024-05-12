package middleware

import (
	"fmt"
	"md/model/common"
	"reflect"
	"runtime"
	"strings"

	"github.com/kataras/iris/v12"
)

// 全局异常恢复
func GlobalRecover(ctx iris.Context) {
	defer func() {
		if err := recover(); err != nil {
			if ctx.IsStopped() {
				return
			}

			// 详细堆栈信息
			var callers []string
			for i := 1; ; i++ {
				_, file, line, got := runtime.Caller(i)
				if !got {
					break
				}
				callers = append(callers, fmt.Sprintf("%s:%d", file, line))
			}

			// 日志记录使用的详细错误信息
			var errMessage string
			// 返回信息
			var errResponse common.ErrorResponse

			// 判断异常类型是否为主动抛出
			if reflect.TypeOf(err) == reflect.TypeOf(common.ErrorResponse{}) {
				errResponse = err.(common.ErrorResponse)
				errMessage = errResponse.Message
			} else {
				// 非主动抛出，使用默认异常信息
				errResponse = common.NewError("内部错误")
				errMessage = fmt.Sprintf("%s", err)
			}

			// 1.通用信息
			logMessage := fmt.Sprintf("[全局异常] 状态码：%d，异常信息：%s\n", errResponse.Code, errMessage)

			// 2.error信息
			if errResponse.Err != nil {
				logMessage += fmt.Sprintf("%s\n", errResponse.Err.Error())
			}

			// 3.最外层方法及行数
			fileName, line := ctx.HandlerFileLine()
			logMessage += fmt.Sprintf("%s:%d (%s)\n", fileName, line, ctx.HandlerName())

			// 4.详细堆栈信息
			logMessage += strings.Join(callers, "\n")

			Log.Error(logMessage)

			ctx.JSON(errResponse)
			ctx.StopExecution()
		}
	}()
	ctx.Next()
}
