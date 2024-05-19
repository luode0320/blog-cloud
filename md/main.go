package main

import (
	"embed"
	"flag"
	"io/fs"
	"md/controller"
	"md/middleware"
	"md/model/common"
	"md/util"
	"net/http"
	"time"

	"github.com/kataras/iris/v12"
)

//go:embed web
var web embed.FS

func init() {
	// 解析命令行参数
	flag.StringVar(&common.Port, "p", "4001", "监听端口")
	flag.StringVar(&common.LogPath, "log", "logs", "日志目录，存放近30天的日志")
	flag.StringVar(&common.DataPath, "data", "data", "数据目录，存放数据库文件和图片")
	flag.BoolVar(&common.Register, "reg", false, "是否允许注册（即使禁止注册，在没有任何用户的情况时仍可注册）")
	flag.StringVar(&common.PostgresHost, "pg_host", "159.138.130.23", "postgres主机地址")
	flag.StringVar(&common.PostgresPort, "pg_port", "15432", "postgres端口")
	flag.StringVar(&common.PostgresUser, "pg_user", "postgres", "postgres用户")
	flag.StringVar(&common.PostgresPassword, "pg_password", "123456", "postgres密码")
	flag.StringVar(&common.PostgresDB, "pg_db", "blog-dev", "postgres数据库名")
	flag.Parse()

	// 固定配置
	common.BasicTokenKey = "md"
	common.ResourceName = ""
	common.PictureName = "picture"
	common.ThumbnailName = "thumbnail"
}

func main() {
	// 创建iris服务, Iris 是一个轻量级的、高性能的 Web 框架
	app := iris.New()

	// 初始化日志
	middleware.InitLog(common.LogPath, app.Logger())

	// 全局异常恢复
	app.Use(middleware.GlobalRecover)

	// gzip压缩
	app.Use(iris.Compression)

	// 初始化雪花算法节点
	err := util.InitSnowflake(0)
	if err != nil {
		middleware.Log.Error("初始化雪花算法节点失败：", err)
		return
	}

	// 初始化数据目录
	err = middleware.InitDataDir(common.DataPath, common.ResourceName, common.PictureName, common.ThumbnailName)
	if err != nil {
		return
	}

	// 初始化数据库连接
	err = middleware.InitDB()
	if err != nil {
		return
	}

	// 初始化API路由
	controller.InitRouter(app)

	// 网页资源路由, 启用静态文件缓存
	middleware.Log.Infof("启用静态文件缓存: 30天")
	app.Use(iris.StaticCache(time.Hour * 720))
	webFs, err := fs.Sub(web, "web")
	if err != nil {
		middleware.Log.Error("初始化网页资源失败：", err)
		return
	}
	app.HandleDir("/", http.FS(webFs))

	// 静态资源路由
	app.HandleDir("/", "")

	// 启动服务
	middleware.Log.Infof("启动服务: {%s}", common.Port)
	app.Logger().Error(app.Run(iris.Addr(":" + common.Port)))
}
