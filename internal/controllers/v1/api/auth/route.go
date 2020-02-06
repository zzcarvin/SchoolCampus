package auth

import (
	"github.com/kataras/iris"
	//"Campus/internal/middleware"

)

func RegisterRoutes(party iris.Party) {

	//party.Post("/", create)
	//party.Delete("/{id:uint64}", remove)
	//party.Put("/{id:uint64}", update)
	//party.Get("/{id:uint64}", get)
	//party.Get("s", search)

	//用学号姓名绑定手机号
	party.Post("/bind", bind)
	party.Post("/exist", exist)
	//party.Put("/bind/{code:string,name:string}", bind)
	//web后台验证用户名密码
	party.Post("/login", login)




}
