package controller

import (
	"md/middleware"
	"md/model/common"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
)

// 初始化iris路由
func InitRouter(app *iris.Application) {
	middleware.Log.Info("初始化iris路由")

	// 允许跨域
	app.UseRouter(cors.AllowAll())

	app.PartyFunc("/api", func(api iris.Party) {
		// 开放接口
		api.PartyFunc("/open", func(open iris.Party) {
			open.Use(middleware.RequestLogger)

			open.Get("/doc/get/{id}", DocumentGetPublished)
			open.Post("/doc/page", DocumentPagePulished)
		})

		// token相关接口
		api.PartyFunc("/token", func(token iris.Party) {
			token.Use(middleware.TokenAuth)
			token.Use(middleware.RequestLogger)

			token.Post("/sign-up", SignUp)
			token.Post("/sign-in", SignIn)
			token.Post("/sign-out", SignOut)
			token.Post("/refresh", TokenRefresh)
		})

		// 数据接口
		api.PartyFunc("/data", func(data iris.Party) {
			data.Use(middleware.DataAuth)

			// 更新密码
			data.PartyFunc("/user", func(user iris.Party) {
				user.Use(middleware.RequestLogger)
				user.Post("/update-password", UserUpdatePassword)
			})

			// 目录
			data.PartyFunc("/book", func(book iris.Party) {
				book.Use(middleware.RequestLogger)
				book.Post("/add", BookAdd)
				book.Post("/update", BookUpdate)
				book.Post("/delete", BookDelete)
				book.Post("/list", BookList)
			})

			// 文档
			data.PartyFunc("/doc", func(doc iris.Party) {
				doc.Use(middleware.RequestLogger)
				doc.Post("/add", DocumentAdd)
				doc.Post("/update", DocumentUpdate)
				doc.Post("/update-content", DocumentUpdateContent)
				doc.Post("/delete", DocumentDelete)
				doc.Post("/list", DocumentList)
				doc.Post("/get", DocumentGet)
			})

			// 图片
			data.PartyFunc("/pic", func(pic iris.Party) {
				pic.Post("/page", PicturePage)
				pic.Post("/delete", PictureDelete)
				pic.Post("/upload", PictureUpload)
			})

			// 非对称密钥
			data.PartyFunc("/rsa", func(rsa iris.Party) {
				rsa.Use(middleware.RequestLogger)
				rsa.Post("/generate", RSAGenerateKey)
				rsa.Post("/encrypt", RSAEncrypt)
				rsa.Post("/decrypt", RSADecrypt)
				rsa.Post("/sign", RSASign)
				rsa.Post("/verify", RSAVerify)
			})
		})
	})
}

// resolveParam 函数用于解析参数
// 参数 ctx 表示 Iris 的上下文对象
// 参数 con 表示要解析的参数对象
func resolveParam(ctx iris.Context, con interface{}) {
	// 使用 ReadJSON 方法将请求的 JSON 数据解析到 con 参数中
	err := ctx.ReadJSON(&con)
	if err != nil {
		// 如果解析失败，则抛出参数解析失败的错误
		panic(common.NewErr("参数解析失败", err))
	}
}
