package fence

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"fmt"
	"github.com/kataras/iris"
)

func get(ctx iris.Context){

	//获取所有围栏
	fenceInfo := make([]models.PlanFence, 0)
	err := lib.Engine.Table("plan_fence").Find(&fenceInfo)
	if err != nil {
		fmt.Printf("查询围栏信息出错：%v",err)
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if len(fenceInfo) == 0 {
		ctx.JSON(lib.NewResponse(0, "没有电子围栏", nil))
		return
	}

	ctx.JSON(lib.NewResponseOK(fenceInfo))
}
