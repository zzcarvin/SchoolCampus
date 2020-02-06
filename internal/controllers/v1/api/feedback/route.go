package feedback

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party) {

	party.Post("/", create)
	party.Get("/", get)
	party.Get("s",search)
	party.Put("/{id:uint64}", update)
	party.Put("record_id/{id:uint64}", statusupdate)
	party.Put("feedback_id/{id:uint64}", feedbackStatus)



	//party.Get("/{id:uint64}", get)
	//party.Get("/search", search)
}
