package list

import (
	"Campus/internal/lib"
	"fmt"
	"github.com/kataras/iris"
	"strings"
	"time"
)

//学生运动列表
func list(ctx iris.Context) {

	//分页

	//排序
	//获取排序信息
	//创建查询Session
	query := lib.Engine.Table("plan_record")

	//字段查询
	//if ctx.URLParamExists("student_id") {
	//	query.And(builder.Like{"student_id", ctx.URLParam("student_id")})
	//}

	//排序
	if ctx.URLParamExists("sort") {
		sort := ctx.URLParam("sort")
		order := strings.ToLower(ctx.URLParamDefault("order", "asc"))
		switch order {
		case "asc":
			query.Asc(sort)
			break
		case "desc":
			query.Desc(sort)
			break
		default:
			ctx.JSON(lib.NewResponseFail(1, "order参数错误，必须是asc或desc"))
			return
		}
	}

	//分页
	page := ctx.URLParamIntDefault("page", 0)
	size := ctx.URLParamIntDefault("size", 15)
	//周期
	cycle := ctx.URLParamIntDefault("cycle", 1)
	lastId := ctx.URLParamIntDefault("last_id", 0)
	//院系,班级
	department := ctx.URLParamIntDefault("department", 0)
	class := ctx.URLParamIntDefault("class", 0)
	//偏移量
	query.Limit(size, page*size)

	println("", lastId, department, class)

	//获取计划名称

	//查询获取条件的计划进度，获取跑步次数，有效公里数
	//获取开始时间和结束时间
	startTimeStr := ""
	endTimeStr := ""
	if cycle == 1 { //周
		now := time.Now()
		cycle := int(now.Weekday()) //周期天数
		if cycle == 0 {
			cycle = 7
		}
		startDate := now.AddDate(0, 0, -(cycle - 1))
		endDate := now.AddDate(0, 0, 1)
		startTimeStr = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location()).Format("2006-01-02")
		endTimeStr = time.Date(startDate.Year(), startDate.Month(), endDate.Day(), 0, 0, 0, 0, startDate.Location()).Format("2006-01-02")
		fmt.Printf("本周开始时间：%v,本周结束时间：%v", startTimeStr, endTimeStr)
	} else { //学期

	}

	//1.1获取本周或本学期开始和结束时间

	//1.2获取院系，班级所有学生的跑步次数和有效公里数

	//1.3计算完成度，获取该集合，遍历该集合将超过总次数，减为总次数并进行累加。总次数：计划次数*所有学生，或，总里程/单次最高里程（有最高里程），总里程/单次最低里程。注：超过100%，减为100%。

	//计算完成度

}
