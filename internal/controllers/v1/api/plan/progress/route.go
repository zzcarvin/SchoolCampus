package progress

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party) {

	party.Post("/", create)
	party.Delete("/{id:uint64}", remove)
	party.Put("/{id:uint64}", update)
	party.Get("/{id:uint64}", get)
	party.Get("s", search)
	party.Get("t", gethistogram)
	party.Get("b", completiondegree)

	//获取计划进度
	party.Get("/getprogress", getprogress)

	//计划图表
	party.Get("/charts", progressChart)

}
