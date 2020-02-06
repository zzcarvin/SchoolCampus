package list

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party) {

	//party.Get("/", departmentsData)

	party.Get("/", classesData)

	//获取所有院系
	party.Get("/departments", departments)
	//获取院系的所有班级
	party.Get("/classes", classes)
	//获取所有体育班级

	//测试获取所有院系每周跑步数据
	//party.Get("/test_week_data", cronProgress)

	//测试获取所有院系每学期跑步数据
	//party.Get("/test_term_data", cronTermProgress)

	//获取list
	party.Get("/progress", proList)

}
