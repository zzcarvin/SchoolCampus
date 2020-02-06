package v1

import (
	"Campus/internal/controllers/v1/api"
	"Campus/internal/controllers/v1/app"
	"Campus/internal/controllers/v1/miniprogram"
	"Campus/internal/controllers/v1/open"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
)

func RegisterRoutes(party iris.Party) {

	//全局跨域
	party.Use(cors.AllowAll())
	party.AllowMethods(iris.MethodOptions)

	app.RegisterRoutes(party.Party("/app"))

	//注册绑定接口
	api.RegisterRoutes(party.Party("/api"))

	//计划进度接口
	open.RegisterRoutes(party.Party("/open"))

	//小程序接口
	miniprogram.RegisterRoutes(party.Party("/miniprogram"))

	//管理员登录接口
	//account.RegisterRoutes(party.Party())
	//运动记录接口
	//planRecoed.RegisterRoutes(app.Party("/planRecord"))
	////运动计划接口
	//plan.RegisterRoutes(app.Party("/plan",middleware.JWT.Serve))

}
