package plan

import (
	"Campus/internal/controllers/v2/app/plan/progress"
	"Campus/internal/controllers/v2/app/plan/record"
	"github.com/kataras/iris"
)

func RegisterRoutes(party iris.Party) {

	party.Post("/", create)
	party.Delete("/{id:uint64}", remove)
	party.Put("/{id:uint64}", update)
	party.Get("/{id:uint64}", get)
	party.Get("s", search)

	//fence.RegisterRoutes(party.Party("/fence"))
	//line.RegisterRoutes(party.Party("/line"))
	//points.RegisterRoutes(party.Party("/points"))
	//progress.RegisterRoutes(party.Party("/progress"))

	record.RegisterRoutes(party.Party("/record"))

	progress.RegisterRoutes(party.Party("/progress"))

}
