package feedback
import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party) {

	party.Post("/", create)
	party.Get("/", get)
	party.Get("s",search)

	//party.Get("/{id:uint64}", get)
	//party.Get("/search", search)
}