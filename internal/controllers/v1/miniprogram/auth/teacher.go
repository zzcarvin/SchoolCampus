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

type bindRequestTeacher struct {
	Id        int    `json:"id"`
	Cellphone string `json:"cellphone"`
}

func teacher(ctx iris.Context) {
	//获取学号和姓名
	code := ""
	name := ""
	if ctx.URLParamExists("code") != true {
		_, _ = ctx.JSON(lib.NewResponseFail(1, "教工号不存在"))
		return
	}
	if ctx.URLParamExists("name") != true {
		_, _ = ctx.JSON(lib.NewResponseFail(1, "姓名不存在"))
		return
	}
	code = ctx.URLParam("code")
	name = ctx.URLParam("name")

	//先查询，存在就将手机号绑定进user，不存在，直接返回
	teacher := models.Teacher{}
	b, err := lib.Engine.Table("teacher").Where("code=?", code).And("name=?", name).Get(&teacher)
	respTeacher := responseTeacher{Exist: false}
	if err != nil {
		fmt.Printf("%v", err)
		_, _ = ctx.JSON(lib.FailureResponse(respTeacher, "查询教师错误"))
		return
	}
	if b == false {
		_, _ = ctx.JSON(lib.FailureResponse(respTeacher, "工号或姓名不正确"))
		return
	}

	if len(teacher.Cellphone) == 11 {
		_, _ = ctx.JSON(lib.FailureResponse(respTeacher, "该教师已绑定，请勿重复绑定。"))
		return
	}

	_, _ = ctx.JSON(lib.SuccessResponse(teacher, "获取教师信息成功"))
	return
}

//绑定教师
func bindTeacher(ctx iris.Context) {
	println("进行更改：")

	var teacherInfo bindRequestTeacher
	err := ctx.ReadJSON(&teacherInfo)
	if err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	teacher := models.Teacher{
		Cellphone: teacherInfo.Cellphone,
	}
	fmt.Printf("body:%v", teacherInfo)
	resNum, err := lib.Engine.Table("teacher").Where("id=?", teacherInfo.Id).Update(&teacher)
	if err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	println("resNum", resNum)
	_, _ = ctx.JSON(lib.SuccessResponse(teacher, "教师绑定成功"))

}
