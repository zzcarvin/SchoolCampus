package app

import (
	"Campus/internal/controllers/v3/app/run"

	"Campus/internal/middleware"
	"github.com/kataras/iris"
)

func RegisterRoutes(app iris.Party) {

	run.RegisterRoutes(app.Party("/run", middleware.JWT.Serve))

	//var student []Student
	//err = lib.Engine.Table("student").Select("id ,name").Where("id in (?)", lib.Engine.Table("studentinfo").Select("id").
	//	Where("status = ?", 2)).Find(&student)
	//SELECT id ,name FROM `student` WHERE (id in (SELECT id FROM `studentinfo` WHERE (status = 2)))

}
