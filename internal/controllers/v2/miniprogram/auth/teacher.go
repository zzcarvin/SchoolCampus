package auth

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"fmt"
	"github.com/kataras/iris"
)

type requestExist struct {
	Code string `json:"code" validate:"max=30"`
	Name string `json:"name" validate:"max=15"`
}

type responseTeacher struct {
	Exist bool `json:"exist"`
}

func teacher(ctx iris.Context) {

	//获取教工号和姓名
	code := ""
	name := ""
	if ctx.URLParamExists("code") {
		code = ctx.URLParam("code")
	} else {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "请输入教工号"))
		return
	}

	if ctx.URLParamExists("name") {
		name = ctx.URLParam("name")
	} else {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "请输入姓名"))
		return
	}

	//先查询，存在就将手机号绑定进user，不存在，直接返回
	teacher := models.Teacher{}
	b, err := lib.Engine.Table("teacher").Where("code=?", code).And("name=?", name).Get(&teacher)
	respTeacher := responseTeacher{Exist: false}
	if err != nil {
		fmt.Printf("%v", err)
		_, _ = ctx.JSON(lib.FailureResponse(respTeacher, "查询教师错误"))
		return
	}
	println("teacherId:", teacher.Id)
	if b == false {
		_, _ = ctx.JSON(lib.FailureResponse(respTeacher, "工号或姓名格式不正确"))
		return
	}

	if len(teacher.Cellphone) == 11 {
		_, _ = ctx.JSON(lib.FailureResponse(respTeacher, "该教师已绑定，请勿重复绑定。"))
		return
	}

	_, _ = ctx.JSON(lib.SuccessResponse(teacher, "获取教师信息成功"))
	return
}
