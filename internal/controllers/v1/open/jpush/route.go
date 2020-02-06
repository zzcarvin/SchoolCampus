package jpush

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party){
	party.Post("/notice",notice)
}