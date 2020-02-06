package plan

import (
	"github.com/kataras/iris"
)

func RegisterRoutes(party iris.Party) {

	party.Post("/", create)
	party.Delete("/{id:uint64}", remove)
	party.Put("/{id:uint64}", update)
	party.Get("/{id:uint64}", get)
	party.Get("s", search)
	party.Get("/class", findclasstype)

	//fence.RegisterRoutes(party.Party("/fence"))
	//line.RegisterRoutes(party.Party("/line"))
	//points.RegisterRoutes(party.Party("/points"))
	//progress. RegisterRoutes(party.Party("/progress"))
	//route.RegisterRoutes(party.Party("/route"))
	//record.	RegisterRoutes(party.Party("/record"))

}
