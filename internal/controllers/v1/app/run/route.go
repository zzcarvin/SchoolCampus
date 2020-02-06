package run

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party) {

	//party.Post("/", create)
	//party.Delete("/{id:uint64}", remove)
	//party.Put("/{id:uint64}", update)
	party.Get("/{id:uint64}", get)
	party.Get("/search", search)
	//party.Get("/{phone:uint64}", getByPhone)
	//party.Get("s", search)
	//party.Get("all", getall)
	party.Post("/start",startRun)
	party.Post("/finish",finishRun)

}
