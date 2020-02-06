package api

import (
	"Campus/internal/controllers/v2/api/plan"
	"Campus/internal/controllers/v2/api/student"
	"github.com/kataras/iris"
)

func RegisterRoutes(app iris.Party) {

	//全局跨域
	//app.Use(cors.AllowAll())
	//app.AllowMethods(iris.MethodOptions)

	//学生接口
	student.RegisterRoutes(app.Party("/student"))
	plan.RegisterRoutes(app.Party("/plan"))

}
