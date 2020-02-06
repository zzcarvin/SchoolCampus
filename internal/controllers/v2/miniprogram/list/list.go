package list

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"fmt"
	"github.com/go-xorm/builder"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris"
)

type depRecLis struct {
	Id         int     `json:"id" xorm:'id'`        //Id
	Name       string  `json:"name" xorm:"name"`    //姓名
	Percentage float32 `json:"percentage" xorm:"-"` //百分比
	//Classes    []int               `json:"classes" xorm:"-"`     //班级id
	Times    int `json:"times" xorm:"-"`    //完成次数
	Distance int `json:"distance" xorm:"-"` //完成公里数
	//每个单位的数据
	Progresses []progressData `json:"progresses" xorm:"-"`
}

type weekData struct {
	WeekTimes    int `json:"week_times" xorm:"week_times"`
	WeekDistance int `json:"week_distance" xorm:"week_distance"`
}

type progressData struct {
	PlanId     int     `json:"plan_id"`
	ClassId    int     `json:"class_id"`
	Times      int     `json:"times"`
	Distance   int     `json:"distance"`
	Percentage float32 `json:"percentage"`
}

type weekProgress struct {
	DepartmentId   int           `json:"department_id"`
	DepartmentName string        `json:"department_name"`
	Plans          []models.Plan `json:"plans"`
}

type planSum struct {
	PlanId       int    `json:"plan_id" xorm:"id"`
	PlanName     string `json:"plan_name" xorm:"name"`
	WeekTimes    int    `json:"week_times" xorm:"-"`
	WeekDistance int    `json:"week_distance" xorm:"-"`
}

type xplanIds struct {
	PlanId int `json:"plan_id" xorm:"plan_id"`
}

type ProgressList struct {
	Id         int     `json:"id" xorm:'id'`                 //Id
	Name       string  `json:"name" xorm:"name"`             //名称
	Percentage float32 `json:"percentage" xorm:"percentage"` //百分比
	Times      int     `json:"times" xorm:"times"`           //完成次数
	Distance   int     `json:"distance" xorm:"distance"`     //完成公里数
}

type responseProgressList struct {
	Records []ProgressList `json:"records"`
	Total   int64          `json:"total"`
	LastId  int            `json:"last_id"`
}

type intArr struct {
	Id int ` xorm:"plan_id"`
}

// "student.id","student.name","student.plan_id","plan.total_distance","plan.total_times","plan.min_week_times","plan.min_week_distance"
type stuPro struct {
	StudentId       int       `xorm:"student.id"`
	StudentName     string    `xorm:"student.name"`
	CreateAt        time.Time `xorm:"student.create_at"`
	PlanId          int       `xorm:"student.plan_id"`
	TotalDistance   int       `xorm:"plan.total_distance"`
	TotalTimes      int       `xorm:"plan.total_times"`
	MinWeekTimes    int       `xorm:"plan.min_week_times"`
	MinWeekDistance int       `xorm:"plan.min_week_distance"`
}

type planPro struct {
	PlanId              int `xorm:"plan.id"`
	PlanTotalDistance   int `xorm:"plan.total_distance"`
	PlanTotalTimes      int `xorm:"plan.total_times"`
	PlanMinWeekTimes    int `xorm:"plan.min_week_times"`
	PlanMinWeekDistance int `xorm:"plan.min_week_distance"`
	ProDistance         int `xorm:"plan_progress.distance"`
	ProTimes            int `xorm:"plan_progress.times"`
	ProWeekDistance     int `xorm:"plan_progress.week_distance"`
	ProWeekTimes        int `xorm:"plan_progress.week_times"`
}

//获取一个系，所有班的全部数据
func classesData(ctx iris.Context) {
	//todo 获取院系id
	//departmentId:=0
	//if ctx.URLParamExists("department_id") {
	//	departmentId,err:=ctx.URLParamInt("student_id")
	//	if err!=nil{
	//		ctx.JSON(lib.FailureResponse(lib.NilStruct(),"院系id不存在"))
	//		return
	//	}
	//}

	//获取一个系所有班
	departmentUnion := make([]depRecLis, 0)
	err := lib.Engine.Table("classes").Where("department_id=?", 1).Find(&departmentUnion)
	if err != nil {
		fmt.Printf("查询班级错误：%v", err)
		return
	}

	fmt.Printf("所有的班级：%v", departmentUnion)

	//判断整个系是否用一个计划(1.全校用一个计划，2.全系用一个计划)
	onlyPlanClass := models.PlanClass{}
	resBool, err := lib.Engine.Table("plan_class").
		Join("INNER", "plan", "plan.id=plan_class.plan_id").
		Where("plan.stop=?", 1).And("plan_class.department_id=?", 0).
		Cols("plan_class.id", "plan_class.plan_id", "plan_class.class_id", "plan_class.department_id", "plan_class.gender").Get(&onlyPlanClass)
	if err != nil {
		fmt.Printf("查询计划院系关系错误：%v", err)
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "查询计划错误"))
		return
	}
	if resBool == false {
		//ctx.JSON(lib.FailureResponse(lib.NilStruct(), "查询计划失败"))

	}

	allDepartments := false
	if onlyPlanClass.DepartmentId == 0 {
		allDepartments = true
	}

	resClassBool, err := lib.Engine.Table("plan_class").
		Join("INNER", "plan", "plan.id=plan_class.plan_id").
		Where("plan.stop=?", 1).And("plan_class.department_id=?", 1).And("plan_class.class_id=?", 0).
		Cols("plan_class.id", "plan_class.plan_id", "plan_class.class_id", "plan_class.department_id", "plan_class.gender").Get(&onlyPlanClass)
	if err != nil {
		fmt.Printf("查询计划院系关系错误：%v", err)
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "查询计划错误"))
		return
	}
	if resClassBool == false {
		fmt.Printf("查询计划院系关系失败")
		//ctx.JSON(lib.FailureResponse(lib.NilStruct(), "查询计划失败"))

	}

	allClass := false
	if onlyPlanClass.ClassId == 0 {
		allClass = true
	}
	println("全系是否用一个计划：", allClass)

	for index, value := range departmentUnion {
		//该系的所有正在运行的计划
		planClass := make([]models.PlanClass, 0)

		if allDepartments {
			//全校用一个计划
			println("全校用一个计划")
			err = lib.Engine.Table("plan_class").
				Join("INNER", "plan", "plan.id=plan_class.plan_id").
				Where("plan_class.department_id=?", 0).And("plan.stop=1").
				Cols("plan_class.id", "plan_class.plan_id", "plan_class.class_id", "plan_class.department_id", "plan_class.gender").Find(&planClass)
			if err != nil {
				fmt.Printf("查询计划院系关系错误：%v", err)
				return
			}
		} else if allClass {
			//全系用一个计划
			println("全系用一个计划")
			err = lib.Engine.Table("plan_class").
				Join("INNER", "plan", "plan.id=plan_class.plan_id").
				Where("plan_class.department_id=?", 1).And("plan.stop=1").And("plan_class.class+id=?", 0).
				Cols("plan_class.id", "plan_class.plan_id", "plan_class.class_id", "plan_class.department_id", "plan_class.gender").Find(&planClass)
			if err != nil {
				fmt.Printf("查询计划院系关系错误：%v", err)
				return
			}
		} else {
			//全系用不同的计划
			println("全系用不同的计划")
			err = lib.Engine.Table("plan_class").
				Join("INNER", "plan", "plan.id=plan_class.plan_id").
				Where("plan_class.class_id=?", value.Id).And("plan.stop=1").
				Cols("plan_class.id", "plan_class.plan_id", "plan_class.class_id", "plan_class.department_id", "plan_class.gender").Find(&planClass)
			if err != nil {
				fmt.Printf("查询计划院系关系错误：%v", err)
				return
			}
		}

		println(value.Name + "所有正在运行的计划：")
		fmt.Printf("%v", planClass)
		println("")

		//初始化数据
		for _, plaClass := range planClass {
			departmentUnion[index].Progresses = append(departmentUnion[index].Progresses, progressData{PlanId: plaClass.PlanId, ClassId: plaClass.ClassId})
		}

		//该系所有计划各自的次数和公里数
		for planIndex, valuePro := range departmentUnion[index].Progresses {
			//获取该系的所有周计划进度
			recordTimes := weekData{}
			resData, err := lib.Engine.Table("plan_progress").Where("plan_id=?", valuePro.PlanId).
				And("class_id=?", value.Id).SumsInt(&recordTimes, "week_times", "week_distance")
			if err != nil {
				fmt.Printf("查询计划进度错误：%v", err)
				return
			}
			println("recordTimes:")
			fmt.Printf("recordTimes:%v", resData)

			departmentUnion[index].Progresses[planIndex].Distance = int(resData[0])
			departmentUnion[index].Progresses[planIndex].Times = int(resData[1])

			//累加次数和公里数
			departmentUnion[index].Times += departmentUnion[index].Progresses[planIndex].Times
			departmentUnion[index].Distance += departmentUnion[index].Progresses[planIndex].Distance
		}

	}

	ctx.JSON(lib.SuccessResponse(departmentUnion, "返回全系信息成功"))

}

func departments(ctx iris.Context) {

	//创建查询Session
	query := lib.Engine.Table("department")

	//字段查询
	if ctx.URLParamExists("name") {
		query.And(builder.Like{"name", ctx.URLParam("name")})
	}

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
	size := ctx.URLParamIntDefault("size", 50)
	query.Limit(size, page*size)

	//查询
	var department []models.Department
	err := query.Find(&department)
	if err != nil {
		fmt.Printf("获取院系出错：%v", err)
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "获取院系错误"))
		return
	}

	ctx.JSON(lib.SuccessResponse(department, "获取院系成功"))

}

//获取院系的所有班级
func classes(ctx iris.Context) {

	//创建查询Session
	query := lib.Engine.Table("classes")

	//字段查询
	if ctx.URLParamExists("name") {
		query.And(builder.Like{"name", ctx.URLParam("name")})
	}
	if ctx.URLParamExists("department_id") {
		query.Where("department_id=?", ctx.URLParamInt64Default("department_id", 0))
	}

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
	size := ctx.URLParamIntDefault("size", 0)
	query.Limit(size, page*size)

	//查询
	var classes []models.Classes
	err1 := query.Find(&classes)
	if err1 != nil {
		fmt.Printf("获取班级出错：%v", err1)
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "获取班级错误"))
		return
	}

	ctx.JSON(lib.SuccessResponse(classes, "获取班级成功"))

}

//首页list
func proList(ctx iris.Context) {

	cycle := 0

	//获取院系班级
	departmentId := 0
	if ctx.URLParamExists("department_id") {
		departmentId = ctx.URLParamIntDefault("department_id", 0)
	} else {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "无院系id"))
		return
	}
	classId := 0
	if ctx.URLParamExists("class_id") {
		classId = ctx.URLParamIntDefault("class_id", 0)
	} else {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "无院系id"))
		return
	}

	//获取时间周期
	if ctx.URLParamExists("cycle") {
		parmCycle, err := strconv.Atoi(ctx.URLParam("cycle"))
		if err != nil {
			fmt.Printf("周期获取错误：%d", err)
			ctx.JSON(lib.FailureResponse(lib.NilStruct(), "周期错误"))
			return
		}
		cycle = parmCycle
	} else {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "无周期"))
		return
	}

	//获取当前学年和学期
	termYear := 0
	term := 0
	yearPlans := make([]models.Plan, 0)
	errPlans := lib.Engine.Table("plan").Desc("create_time").Find(&yearPlans)
	if errPlans != nil {
		fmt.Printf("查询所有计划出错：%v", errPlans)
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "该系没有计划"))
		return
	}
	//获取今年的年份和学期
	if len(yearPlans) != 0 {
		termYear = yearPlans[0].TermYear
		term = yearPlans[0].Term
	}

	println("termYear:", termYear, "term:", term)

	//获取单个班级的一周，一学期数据---开始
	if departmentId != 0 && classId != 0 {
		stuQuery := lib.Engine.Table("student")
		//排序
		if ctx.URLParamExists("sort") {
			sort := "student." + ctx.URLParam("sort")
			order := strings.ToLower(ctx.URLParamDefault("order", "asc"))
			switch order {
			case "asc":
				stuQuery.Asc(sort)
				break
			case "desc":
				stuQuery.Desc(sort)
				break
			default:
				ctx.JSON(lib.FailureResponse(lib.NilStruct(), "order参数错误，必须是asc或desc"))
				return
			}
		}

		//分页:用last_id来做限制条件，而且使用last_id就不用page，使用最新的一页
		//size := ctx.URLParamIntDefault("size", 5)
		//lastId := ctx.URLParamIntDefault("last_id", 0)
		//跳过的数量
		//stuQuery.Limit(size, 0)

		//当不是第一页时,这里的lastid应该谨慎使用，
		//if lastId != 0 {
		//	stuQuery.And("student.create_at<(select student.create_at from student where student.id=?)", lastId)
		//}

		//获取该班级所有正在运行计划的学生
		//studens:=[]stuPro{}
		//stuErr:=stuQuery.Join("INNER","plan","student.plan_id=plan.id").
		//	Where("student.department_id=?",departmentId).And("student.class_id=?",classId).And("plan.stop=?",1).
		//	Cols("student.id","student.name","student.plan_id","student.create_at","plan.total_distance","plan.total_times","plan.min_week_times","plan.min_week_distance").Find(&studens)
		//获取包括有计划和没有计划的学生
		students := []models.Student{}
		stuErr := stuQuery.Where("student.department_id=?", departmentId).And("student.class_id=?", classId).Find(&students)
		if stuErr != nil {
			fmt.Printf("查询学生计划出错：%v", stuErr)
			ctx.JSON(lib.FailureResponse(lib.NilStruct(), "该学生没有计划"))
			return
		}
		println("students len:", len(students))
		var progressList []ProgressList
		//当所有的学生没有计划
		//if len(students)==0{
		//	resProgress := responseProgressList{progressList, 0, 0}
		//	ctx.JSON(lib.SuccessResponse(resProgress, "获取学校运动进度成功"))
		//	return
		//}

		for _, value := range students {
			progress := planPro{}

			res, errPro := lib.Engine.Table("plan_progress").
				Join("INNER", "plan", "plan.id=plan_progress.plan_id").
				Where("plan_progress.student_id=?", value.Id).And("plan_progress.plan_id=?", value.PlanId).
				Cols("plan.id", "plan.total_distance", "plan.total_times", "plan.min_week_times", "plan.min_week_distance", "plan_progress.distance", "plan_progress.times", "plan_progress.week_distance", "plan_progress.week_times").
				Get(&progress)
			if errPro != nil {
				fmt.Printf("查询学生计划进度出错：%v", errPro)
				ctx.JSON(lib.FailureResponse(lib.NilStruct(), "该学生没有计划进度"))
				return
			}
			contProgress := ProgressList{
				Id:   value.Id,
				Name: value.Name,
			}
			//没有计划进度，数据为0
			if res == false {
				progressList = append(progressList, contProgress)
			} else {
				println("有计划进度")
				println("times:", contProgress.Times, "distance:", contProgress.Distance, "percentage:", contProgress.Percentage)
				//有计划进度
				if cycle == 1 {
					//周
					contProgress.Times = progress.ProWeekTimes
					contProgress.Distance = progress.ProWeekDistance
					if progress.PlanMinWeekTimes != 0 {
						contProgress.Percentage = float32(progress.ProWeekTimes) / float32(progress.PlanMinWeekTimes)
					} else {
						contProgress.Percentage = 0
					}

				} else {
					//学期
					contProgress.Times = progress.ProTimes
					contProgress.Distance = progress.ProDistance
					if progress.PlanMinWeekTimes != 0 {
						contProgress.Percentage = float32(progress.ProTimes) / float32(progress.PlanTotalTimes)
					} else {
						contProgress.Percentage = 0
					}
				}
				println("times:", contProgress.Times, "distance:", contProgress.Distance, "percentage:", contProgress.Percentage)
				progressList = append(progressList, contProgress)
			}

			fmt.Printf("percentage: %v", contProgress.Percentage)

		}

		retLastId := 0
		if len(progressList) != 0 {
			retLastId = progressList[len(progressList)-1].Id
		}
		println("progressList len:", len(progressList))
		fmt.Printf("progressList:%v,list len:%d,lastId:%d", progressList, int64(len(progressList)), retLastId)
		resProgress := responseProgressList{progressList, int64(len(progressList)), retLastId}
		fmt.Printf("resprogress:%v", resProgress)

		ctx.JSON(lib.SuccessResponse(resProgress, "获取学校运动进度成功"))
		return

	}
	//获取单个班级的一周，一学期数据---结束

	//每周，每学期院系和班级数据
	if cycle == 1 {
		//周
		query := lib.Engine.Table("week_progress")
		//排序
		if ctx.URLParamExists("sort") {
			sort := "week_progress." + ctx.URLParam("sort")
			order := strings.ToLower(ctx.URLParamDefault("order", "asc"))
			switch order {
			case "asc":
				query.Asc(sort)
				break
			case "desc":
				query.Desc(sort)
				break
			default:
				ctx.JSON(lib.FailureResponse(lib.NilStruct(), "order参数错误，必须是asc或desc"))
				return
			}
		}

		//分页:用last_id来做限制条件，而且使用last_id就不用page，使用最新的一页
		size := ctx.URLParamIntDefault("size", 5)
		lastId := ctx.URLParamIntDefault("last_id", 0)
		//跳过的数量
		query.Limit(size, 0)
		println("classId:", classId)
		//获取院系或班级名称,并查询院系和班级
		if departmentId == 0 {
			query.Join("INNER", "department", "department.id=week_progress.department_id")
			query.And("week_progress.class_id=?", 0)
		} else {
			query.Join("INNER", "classes", "classes.id=week_progress.class_id")
			query.And("week_progress.department_id=?", departmentId)
			if classId == 0 {
				query.And("week_progress.class_id!=?", 0)
			} else {
				query.And("week_progress.class_id=?", classId)
			}

		}

		//当不是第一页时
		if lastId != 0 {
			query.And("week_progress.create_at<(select week_progress.create_at from week_progress where week_progress.id=?)", lastId)
		}

		//查询
		var progressList []ProgressList
		toatal, err := query.Where("term_year=?", termYear).And("term=?", term).FindAndCount(&progressList)
		if err != nil {
			fmt.Printf("week,查询运动记录错误：%v", err)
			ctx.JSON(lib.FailureResponse(lib.NilStruct(), "查询运动记录错误"))
			return
		}

		retLastId := 0
		if len(progressList) != 0 {
			retLastId = progressList[len(progressList)-1].Id
		}

		resProgress := responseProgressList{progressList, toatal, retLastId}

		//获取total
		if departmentId == 0 {
			totalRecord := models.WeekProgress{}
			num, err := lib.Engine.Table("week_progress").Join("INNER", "department", "department.id=week_progress.department_id").Where("department.id!=?", 0).
				And("week_progress.class_id=?", 0).And("week_progress.term=?", term).And("week_progress.term_year=?", termYear).Count(totalRecord)
			if err != nil {
				fmt.Printf("获取total错误：%v", err)
				lib.FailureResponse(lib.NilStruct(), "获取total错误")
				return
			}
			resProgress.Total = num
		} else {
			totalRecord := models.WeekProgress{}
			if classId == 0 {
				num, err := lib.Engine.Table("week_progress").Join("INNER", "classes", "classes.id=week_progress.class_id").Where("week_progress.class_id!=?", 0).
					And("week_progress.department_id=?", departmentId).And("week_progress.term=?", term).And("week_progress.term_year=?", termYear).Count(totalRecord)
				if err != nil {
					fmt.Printf("获取total错误：%v", err)
					lib.FailureResponse(lib.NilStruct(), "获取total错误")
					return
				}
				resProgress.Total = num
			} else {
				num, err := lib.Engine.Table("week_progress").Join("INNER", "classes", "classes.id=week_progress.class_id").Where("classes.id=?", classId).And("week_progress.department_id=?", departmentId).
					And("week_progress.class_id=?", classId).And("week_progress.term=?", term).And("week_progress.term_year=?", termYear).Count(totalRecord)
				if err != nil {
					fmt.Printf("获取total错误：%v", err)
					lib.FailureResponse(lib.NilStruct(), "获取total错误")
					return
				}
				println("")
				fmt.Printf("num:%d", num)
				resProgress.Total = num
			}

		}

		ctx.JSON(lib.SuccessResponse(resProgress, "获取学校运动进度成功"))

	} else if cycle == 2 {
		//学期
		query := lib.Engine.Table("term_progress")
		//排序
		if ctx.URLParamExists("sort") {
			sort := "term_progress." + ctx.URLParam("sort")
			order := strings.ToLower(ctx.URLParamDefault("order", "asc"))
			switch order {
			case "asc":
				query.Asc(sort)
				break
			case "desc":
				query.Desc(sort)
				break
			default:
				ctx.JSON(lib.FailureResponse(lib.NilStruct(), "order参数错误，必须是asc或desc"))
				return
			}
		}

		//分页:用last_id来做限制条件，而且使用last_id就不用page，使用最新的一页
		size := ctx.URLParamIntDefault("size", 5)
		lastId := ctx.URLParamIntDefault("last_id", 0)
		//跳过的数量
		query.Limit(size, 0)
		println("classId:", classId)
		//获取院系或班级名称,并查询院系和班级
		if departmentId == 0 {
			query.Join("INNER", "department", "department.id=term_progress.department_id")
			query.And("term_progress.class_id=?", 0)
		} else {
			query.Join("INNER", "classes", "classes.id=term_progress.class_id")
			query.And("term_progress.department_id=?", departmentId)
			if classId == 0 {
				query.And("term_progress.class_id!=?", 0)
			} else {
				query.And("term_progress.class_id=?", classId)
			}

		}

		//当不是第一页时
		if lastId != 0 {
			query.And("term_progress.create_at<(select term_progress.create_at from term_progress where term_progress.id=?)", lastId)
		}

		//查询
		var progressList []ProgressList
		toatal, err := query.Where("term_year=?", termYear).And("term=?", term).FindAndCount(&progressList)
		if err != nil {
			fmt.Printf("查询运动记录错误：%v", err)
			ctx.JSON(lib.FailureResponse(lib.NilStruct(), "查询运动记录错误"))
			return
		}

		retLastId := 0
		if len(progressList) != 0 {
			retLastId = progressList[len(progressList)-1].Id
		}

		resProgress := responseProgressList{progressList, toatal, retLastId}

		//获取total
		if departmentId == 0 {
			totalRecord := models.WeekProgress{}
			num, err := lib.Engine.Table("term_progress").Join("INNER", "department", "department.id=term_progress.department_id").Where("department.id!=?", 0).
				And("term_progress.class_id=?", 0).And("term_progress.term=?", term).And("term_progress.term_year=?", termYear).Count(totalRecord)
			if err != nil {
				fmt.Printf("获取total错误：%v", err)
				lib.FailureResponse(lib.NilStruct(), "获取total错误")
				return
			}
			resProgress.Total = num
		} else {
			totalRecord := models.WeekProgress{}
			if classId == 0 {
				num, err := lib.Engine.Table("term_progress").Join("INNER", "classes", "classes.id=term_progress.department_id").Where("term_progress.class_id!=?", 0).
					And("term_progress.department_id=?", departmentId).And("term_progress.term=?", term).And("term_progress.term_year=?", termYear).Count(totalRecord)
				if err != nil {
					fmt.Printf("获取total错误：%v", err)
					lib.FailureResponse(lib.NilStruct(), "获取total错误")
					return
				}
				resProgress.Total = num
			} else {
				num, err := lib.Engine.Table("term_progress").Join("INNER", "classes", "classes.id=term_progress.department_id").Where("classes.id=?", classId).And("term_progress.department_id=?", departmentId).
					And("term_progress.class_id=?", classId).And("term_progress.term=?", term).And("term_progress.term_year=?", termYear).Count(totalRecord)
				if err != nil {
					fmt.Printf("获取total错误：%v", err)
					lib.FailureResponse(lib.NilStruct(), "获取total错误")
					return
				}
				println("")
				fmt.Printf("num:%d", num)
				resProgress.Total = num
			}

		}

		ctx.JSON(lib.SuccessResponse(resProgress, "获取学校运动进度成功"))
	}

}
