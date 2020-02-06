package teacher

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party) {

	//获取教师信息
	party.Get("/{id:uint64}", get)

	//party.Post("/student", student)

	//基本信息 dev.api.xiaotibang.com/app/student/:id

	//头像

	//party.Delete("/{id:uint64}", remove)
	//party.Put("/{id:uint64}", update)
	//party.Get("/{id:uint64}", get)
	//party.Get("s", search)

}
