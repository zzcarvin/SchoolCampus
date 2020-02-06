package auth

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party) {

	party.Post("/teacher", teacher)

	//绑定教师
	party.Post("/bind_teacher", bindTeacher)

	//party.Delete("/{id:uint64}", remove)
	//party.Put("/{id:uint64}", update)
	//party.Get("/{id:uint64}", get)
	//party.Get("s", search)

}
