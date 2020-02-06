package miniprogram

import (
	"Campus/internal/controllers/v1/miniprogram/auth"
	"Campus/internal/controllers/v1/miniprogram/list"
	"Campus/internal/controllers/v1/miniprogram/teacher"
	"github.com/kataras/iris"
)

func RegisterRoutes(app iris.Party) {

	auth.RegisterRoutes(app.Party("/auth"))

	list.RegisterRoutes(app.Party("/list"))

	teacher.RegisterRoutes(app.Party("/teacher"))

}
