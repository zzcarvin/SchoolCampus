package mara

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"fmt"
	"github.com/kataras/iris"
)

func ApplyMara(ctx iris.Context) {
	applyMara := models.ApplyMara{}
	err := ctx.ReadForm(&applyMara)
	if err != nil {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), err.Error()))
		return
	}
	valType := lib.ValidateRequest(applyMara)
	if valType == false {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "请确保学号，姓名，班级，手机号都正确填写"))
		return
	}
	//判断学号是否存在
	has, err := lib.Engine.Table("apply_mara").Where("code = ?", applyMara.Code).Exist()
	if err != nil {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "服务器错误"))
		return
	}
	if has {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "你已经报过名了，请不要重复报名"))
		return
	}

	//判断超过800
	counts, err := lib.Engine.Table("apply_mara").Count()
	fmt.Printf("counts:%d", counts)
	if counts >= 800 {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "您来晚了，报名人数到上限了"))
		return
	}

	//入库
	_, err = lib.Engine.Table("apply_mara").Insert(&applyMara)

	if err != nil {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), err.Error()))
		return
	}

	ctx.JSON(lib.SuccessResponse(lib.NilStruct(), "报名成功"))

}

func GetStatus(ctx iris.Context) {

}
