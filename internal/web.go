package internal

import (
	"Campus/configs"
	v1 "Campus/internal/controllers/v1"
	v2 "Campus/internal/controllers/v2"
	v3 "Campus/internal/controllers/v3"
	"Campus/internal/middleware"
	"github.com/kataras/iris"
)

func WebServe() {
	//获取配置
	cfg := configs.Conf.Web

	//默认iris app
	app := iris.Default()

	//初始化JWT
	middleware.JwtInit()

	//总路由
	v1.RegisterRoutes(app)
	v2.RegisterRoutes(app)
	v3.RegisterRoutes(app)
	//静态目录
	app.StaticWeb("/", cfg.Static)

	//启动监听
	app.Run(iris.Addr(cfg.Addr))
}
