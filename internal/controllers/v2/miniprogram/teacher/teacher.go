package teacher

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"github.com/kataras/iris"
)

func get(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	teacher := models.Teacher{}
	//根据id查询
	b, err := lib.Engine.Table("teacher").ID(id).Get(&teacher)

	if err != nil {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "查询教师失败"))
		return
	}
	if b == false {
		ctx.JSON(lib.FailureResponse(1, "未找到该教工"))
		return
	}
	ctx.JSON(lib.SuccessResponse(teacher, "获取教师信息成功"))
}
