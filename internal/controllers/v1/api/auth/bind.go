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

type requestBind struct {
	Code      string `json:"code" validate:"max=30"`
	Name      string `json:"name" validate:"max=15"`
	Cellphone string `json:"cellphone" validate:"max=11"`
}

// 用姓名学号查询学生是否存在，存在返回学生整合信息，不存在直接返回

func exist(ctx iris.Context) {

	//获取学号和姓名
	var faceMaterial requestExist
	if err := ctx.ReadJSON(&faceMaterial); err != nil {
		ctx.JSON(iris.Map{"phone ReadJSON error": "错误的手机号码"})
		fmt.Println("phone ReadJSON error", err)
		return
	}

	//验证
	valPhone := lib.ValidateRequest(faceMaterial)
	if valPhone == false {
		ctx.JSON(lib.NewResponseFail(1, "手机号或状态码格式错误"))
		return
	}

	fmt.Print(faceMaterial.Code, faceMaterial.Name)
	//先查询，存在就将手机号绑定进user，不存在，直接返回
	student := models.Student{}
	b, err := lib.Engine.Table("student").Where("code=?", faceMaterial.Code).And("name=?", faceMaterial.Name).Get(&student)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该用户"))
		return
	}

	//返回学生的全部信息
	studentAllInfo := models.StudentAllInfos{}
	//根据id查询
	b, err = lib.Engine.Table("student").
		Join("INNER", "classes", "classes.id=student.class_id").
		Join("INNER", "department", "department.id=student.department_id").
		Where("student.id=?", student.Id).
		Cols("student.id", "student.code", "student.name", "classes.name", "department.name", "student.gender", "student.create_at", "student.cellphone").
		Get(&studentAllInfo)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该用户"))
		return
	}
	ctx.JSON(lib.NewResponseOK(studentAllInfo))
	return

}

// swagger:route POST /api/auth/bind  bind bind
//
// 用姓名学号绑定手机号
//     Produces:
//     - application/json
//
//     Responses:
//       200: Response
func bind(ctx iris.Context) {

	//获取学号和姓名
	var faceMaterial requestBind
	if err := ctx.ReadJSON(&faceMaterial); err != nil {
		ctx.JSON(iris.Map{"phone ReadJSON error": "错误的手机号码"})
		fmt.Println("phone ReadJSON error", err)
		return
	}

	//验证
	valPhone := lib.ValidateRequest(faceMaterial)
	if valPhone == false {
		ctx.JSON(lib.NewResponseFail(1, "手机号或状态码格式错误"))
		return
	}
	//插入数据
	student := models.Student{}
	student.Cellphone = faceMaterial.Cellphone
	res, err := lib.Engine.Table("student").Where("code=?", faceMaterial.Code).And("name=?", faceMaterial.Name).Update(student)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	if res == 1 {
		ctx.JSON(lib.NewResponseOK("学生手机号绑定成功"))

		return
	}
	ctx.JSON(lib.NewResponse(1, "学生手机号绑定失败", "学生手机号绑定失败"))
	println("res lines:", res)
	return
}
