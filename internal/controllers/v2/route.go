package v2

import (
	"Campus/internal/controllers/v2/api"
	"Campus/internal/controllers/v2/app"
	"Campus/internal/controllers/v2/miniprogram"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
)

func RegisterRoutes(party iris.Party) {

	//全局跨域
	party.Use(cors.AllowAll())
	party.AllowMethods(iris.MethodOptions)

	app.RegisterRoutes(party.Party("/v2/app"))

	//注册绑定接口
	api.RegisterRoutes(party.Party("/v2/api"))

	//管理员登录接口
	//account.RegisterRoutes(party.Party())
	//运动记录接口
	//planRecoed.RegisterRoutes(app.Party("/planRecord"))
	////运动计划接口
	//plan.RegisterRoutes(app.Party("/plan",middleware.JWT.Serve))

	//小程序接口
	miniprogram.RegisterRoutes(party.Party("/v2/miniprogram"))

}
