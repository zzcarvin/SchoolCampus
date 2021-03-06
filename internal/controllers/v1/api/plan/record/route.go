package record

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party) {
	party.Post("/", create)
	party.Delete("/{id:uint64}", remove)
	party.Put("/{id:uint64}", update)
	party.Get("/{id:uint64}", get)

	party.Get("s", search)
	party.Get("t", totaldayrun)

	party.Get("/getrunfailreason/{id:uint64}", getPlanRecordString)

}
