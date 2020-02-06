package miniprogram

import (
	"Campus/internal/controllers/v2/miniprogram/auth"
	"Campus/internal/controllers/v2/miniprogram/list"
	"Campus/internal/controllers/v2/miniprogram/student"
	"Campus/internal/controllers/v2/miniprogram/teacher"
	"Campus/internal/middleware"
	"github.com/kataras/iris"
)

func RegisterRoutes(app iris.Party) {
	//登录注册
	auth.RegisterRoutes(app.Party("/auth", middleware.JWT.Serve))
	//院系班级计划列表
	list.RegisterRoutes(app.Party("/list", middleware.JWT.Serve))
	//学生运动记录列表
	student.RegisterRoutes(app.Party("/student", middleware.JWT.Serve))
	//教师信息
	teacher.RegisterRoutes(app.Party("/teacher"))

}
