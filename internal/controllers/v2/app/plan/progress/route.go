package progress

import (
	"github.com/kataras/iris"
)

func RegisterRoutes(party iris.Party) {

	//新计划首页
	party.Get("/{id:uint64}", planProgress)

	//简略学期跑步进度
	party.Post("/detail", detailPlanProgress)

}
