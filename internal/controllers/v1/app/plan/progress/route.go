package progress

import (
"github.com/kataras/iris"
)

func RegisterRoutes(party iris.Party) {

	//party.Post("/", create)
	//party.Delete("/{id:uint64}", remove)
	//party.Put("/{id:uint64}", update)
	//party.Get("/{id:uint64}", get)
	//party.Get("s", search)

	//fence.RegisterRoutes(party.Party("/fence"))
	//line.RegisterRoutes(party.Party("/line"))
	//points.RegisterRoutes(party.Party("/points"))
	//progress.RegisterRoutes(party.Party("/progress"))

	//party.Get("/{id:uint64}", get)
	//新计划首页
	party.Get("/{id:uint64}",planProgress)

	//简略计划页
	party.Post("/simple",planSimple)

	//简略学期跑步进度
	party.Post("/detail",detailPlanProgress)

}
