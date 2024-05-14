package middleware

import (
	"md/model/common"
	"strings"
	"unicode/utf8"

	"github.com/kataras/iris/v12"
	"github.com/muesli/cache2go"
)

// 数据接口授权
func DataAuth(ctx iris.Context) {
	token := resolveHeader(ctx, "Bearer")

	// 检验缓存中是否存在此token
	if !cache2go.Cache(common.AccessTokenCache).Exists(token) {
		panic(common.NewErrorCode(common.HttpAuthFailure, "认证失败"))
	}

	ctx.Next()
}

// TokenAuth 函数用于进行 token 相关接口的认证授权
// 参数 ctx 表示 Iris 的上下文对象
func TokenAuth(ctx iris.Context) {
	//// 从请求头中解析出 token
	//token := resolveHeader(ctx, "Basic")
	//
	//// 计算当前时间戳的10分钟为基准的值，并进行 SHA256 加密
	//current := time.Now().UnixMilli() / 600000
	//t1 := util.EncryptSHA256([]byte(common.BasicTokenKey + strconv.FormatInt(current, 10)))
	//t2 := util.EncryptSHA256([]byte(common.BasicTokenKey + strconv.FormatInt(current-1, 10)))
	//t3 := util.EncryptSHA256([]byte(common.BasicTokenKey + strconv.FormatInt(current+1, 10)))
	//
	//// 验证 token 是否在允许的范围内
	//if token != t1 && token != t2 && token != t3 {
	//	// 如果 token 不在允许的范围内，则抛出认证失败的错误
	//	panic(common.NewErrorCode(common.HttpAuthFailure, "认证失败"))
	//}

	// 继续处理请求
	ctx.Next()
}

// 获取当前登录用户id
func CurrentUserId(ctx iris.Context) string {
	token := resolveHeader(ctx, "Bearer")
	res, err := cache2go.Cache(common.AccessTokenCache).Value(token)
	if err != nil {
		panic(common.NewErrorCode(common.HttpAuthFailure, "认证失败"))
	}
	tokenCache := res.Data().(*common.TokenCache)
	if tokenCache.Id == "" {
		panic(common.NewErrorCode(common.HttpAuthFailure, "认证失败"))
	}
	return tokenCache.Id
}

// resolveHeader 函数用于解析头信息中的认证信息
// 参数 ctx 表示 Iris 的上下文对象
// 参数 prefix 表示认证信息的前缀
// 返回解析出的认证信息字符串
func resolveHeader(ctx iris.Context, prefix string) string {
	// 获取请求头中的 Authorization 信息
	header := ctx.GetHeader("Authorization")

	// 计算前缀的长度
	prefixLen := utf8.RuneCountInString(prefix) + 1

	// 判断 Authorization 信息是否以指定前缀开头，并且长度大于前缀长度
	if header != "" && strings.Index(header, prefix) == 0 && utf8.RuneCountInString(header) > prefixLen {
		// 返回去除前缀后的认证信息
		return string([]rune(header)[prefixLen:])
	}

	// 如果认证信息不符合要求，则抛出认证失败的错误
	panic(common.NewErrorCode(common.HttpAuthFailure, "认证失败"))
}
