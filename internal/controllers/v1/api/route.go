package api

import (
	"Campus/internal/controllers/v1/api/account"
	"Campus/internal/controllers/v1/api/announcement"
	"Campus/internal/controllers/v1/api/auth"
	"Campus/internal/controllers/v1/api/classes"
	"Campus/internal/controllers/v1/api/department"
	"Campus/internal/controllers/v1/api/feedback"
	"Campus/internal/controllers/v1/api/plan"
	"Campus/internal/controllers/v1/api/student"
	"Campus/internal/controllers/v1/api/teacher"
	"Campus/internal/middleware"
	"github.com/kataras/iris"
)

func RegisterRoutes(app iris.Party) {

	//全局跨域
	//app.Use(cors.AllowAll())
	//app.AllowMethods(iris.MethodOptions)

	//注册绑定接口
	auth.RegisterRoutes(app.Party("/auth"))
	//返回登陆信息接口
	account.RegisterRoutes(app.Party("/account", middleware.JWT.Serve))
	//院系接口
	department.RegisterRoutes(app.Party("/department", middleware.JWT.Serve))
	//班级接口
	classes.RegisterRoutes(app.Party("/classes", middleware.JWT.Serve))
	//学生接口
	student.RegisterRoutes(app.Party("/student"))
	//运动计划接口
	plan.RegisterRoutes(app.Party("/plan", middleware.JWT.Serve))
	announcement.RegisterRoutes(app.Party("/announcement", middleware.JWT.Serve))
	feedback.RegisterRoutes(app.Party("/feedback", middleware.JWT.Serve))

	//教师接口
	teacher.RegisterRoutes(app.Party("/teacher", middleware.JWT.Serve))
}
