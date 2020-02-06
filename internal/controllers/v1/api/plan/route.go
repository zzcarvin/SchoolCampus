package plan

import (
	"Campus/internal/controllers/v1/api/plan/fence"
	"Campus/internal/controllers/v1/api/plan/line"
	"Campus/internal/controllers/v1/api/plan/points"
	"Campus/internal/controllers/v1/api/plan/progress"
	"Campus/internal/controllers/v1/api/plan/record"
	"github.com/kataras/iris"

	"Campus/internal/controllers/v1/api/plan/route"
)

//路由
func RegisterRoutes(party iris.Party) {

	party.Post("/", create)
	party.Delete("/{id:uint64}", remove)
	//party.Put("/", update)
	party.Put("/", newUpdate)
	party.Get("/{id:uint64}", get)
	party.Get("s", search)
	party.Get("/class", findclasstype)
	//计划信息卡获取
	party.Get("/planinformationcard", planinformationcard)
	party.Get("/planblackboard", blackboard)
	//party.Get("/allplanblackboard",allblackboard)
	party.Get("/addprogress", addprogress)
	party.Get("/test_record", TestRecord)

	//学生个人页图表
	party.Get("/gethistogram", gethistogram)
	//运动趋势接口
	party.Get("/movementtrend", movementtrend)

	//图表接口

	party.Get("/doughnutchart", doughnutchart)

	fence.RegisterRoutes(party.Party("/fence"))
	line.RegisterRoutes(party.Party("/line"))
	points.RegisterRoutes(party.Party("/points"))
	progress.RegisterRoutes(party.Party("/progress"))
	route.RegisterRoutes(party.Party("/route"))
	record.RegisterRoutes(party.Party("/record"))

}
