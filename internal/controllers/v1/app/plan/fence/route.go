package fence

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party) {

	party.Get("/",get)



}

