package teacher

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party) {

	//教师导入

	party.Post("/", create)
	party.Delete("/{id:uint64}", remove)
	party.Put("/{id:uint64}", update)
	party.Get("/{id:uint64}", get)

	party.Get("s", search)

	party.Post("/excel", excel)

	//批量更新教师信息，包括教师姓名，院系，性别，管理的体育班
	party.Post("/many_excel", manyUpdateExcel)

	//party.Post("/update_excel", UpdateExcel)

}
