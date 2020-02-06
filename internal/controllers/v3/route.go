package v2

import (
	"Campus/internal/controllers/v3/app"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
)

func RegisterRoutes(party iris.Party) {

	//全局跨域
	party.Use(cors.AllowAll())
	party.AllowMethods(iris.MethodOptions)

	app.RegisterRoutes(party.Party("/v3/app"))



	//管理员登录接口
	//account.RegisterRoutes(party.Party())
	//运动记录接口
	//planRecoed.RegisterRoutes(app.Party("/planRecord"))
	////运动计划接口
	//plan.RegisterRoutes(app.Party("/plan",middleware.JWT.Serve))

}
