package auth

import (
	"github.com/kataras/iris"
)

func RegisterRoutes(party iris.Party) {

	//party.Post("/", create)
	//party.Delete("/{id:uint64}", remove)
	//party.Put("/{id:uint64}", update)
	//party.Get("/{id:uint64}", get)
	//party.Get("s", search)

	//party.Use(cors.AllowAll())

	//用学号姓名绑定手机号
	party.Post("/bind", bind)
	party.Post("/exist", exist)

	//party.Put("/bind/{code:string,name:string}", bind)
	//该接口已废弃
	//party.Post("/exist", login)

}
