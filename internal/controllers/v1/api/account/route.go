package account

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party) {
	//账号验证
	party.Get("/", account)
	//修改密码
	party.Put("/modifypassword", modifypassword)

}
