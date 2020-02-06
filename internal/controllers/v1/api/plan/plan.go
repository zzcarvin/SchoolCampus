package plan

import (
	"Campus/internal/controllers/v3/app/run"
	"Campus/internal/lib"
	"Campus/internal/models"
	"bytes"
	"fmt"
	"github.com/go-xorm/builder"
	"github.com/kataras/iris"
	"log"
	"strconv"
	"strings"
	"time"
)

// swagger:parameters  PlanCreateRequest
type PlanCreateRequest struct {
	// in: body
	Body models.Plan
}

type planandtime struct {
	//计划id
	Id int `json:"id" xorm:"not null pk autoincr comment('计划Id') INT(11)"`

	//院系id
	DepartmentId int `json:"department_id" xorm:"not null comment('院系id') index INT(5)"`

	//班级id
	ClassId int `json:"class_id" xorm:"not null comment('班级id') index INT(5)"`

	//计划名称
	Name string `json:"name" xorm:"not null comment('计划名称') VARCHAR(255)"`

	//计划开始日期
	DateBegin time.Time `json:"date_begin" xorm:"not null comment('计划开始日期') DATETIME"`

	//计划结束日期
	DateEnd time.Time `json:"date_end" xorm:"not null comment('计划结束日期') DATETIME"`

	//总里程
	TotalDistance int `json:"total_distance" xorm:"not null comment('总里程') INT(4)"`

	//总次数
	TotalTimes int `json:"total_times" xorm:"not null comment('总次数') INT(4)"`

	//最低周跑次数
	MinWeekTimes int `json:"min_week_times" xorm:"not null comment('最低周跑次数') INT(2)"`

	//最低周跑里程
	MinWeekDistance int `json:"min_week_distance" xorm:"not null comment('最低周跑里程') INT(11)"`

	//最低单次里程
	MinSingleDistance int `json:"min_single_distance" xorm:"not null comment('最低单次里程') INT(11)"`

	//最高周跑次数
	MaxWeekTimes int `json:"max_week_times" xorm:"not null comment('最高周跑次数') INT(11)"`

	//最高周跑里程
	MaxWeekDistance int `json:"max_week_distance" xorm:"not null comment('最高周跑里程') INT(11)"`

	//最高单次里程
	MaxSingleDistance int `json:"max_single_distance" xorm:"not null comment('最高单次里程') INT(11)"`

	//最高日跑次数
	MaxDayTimes int `json:"max_day_times" xorm:"not null comment('最高日跑次数') INT(2)"`

	//必跑日
	MustRunDay string `json:"must_run_day" xorm:"not null comment('必跑日') VARCHAR(20)"`

	//最小配速(单位s)
	MinPace int `json:"min_pace" xorm:"not null comment('最小配速(单位s)') INT(11)"`

	//最大配速(单位s)
	MaxPace int `json:"max_pace" xorm:"not null comment('最大配速(单位s)') INT(11)"`

	//最短时间(单位s)
	MinTimeLong int `json:"min_time_long" xorm:"not null comment('最短时间(单位s)') INT(11)"`

	//最长时间(单位s)
	MaxTimeLong int `json:"max_time_long" xorm:"not null comment('最长时间(单位s)') INT(11)"`

	//创建时间
	CreateTime time.Time `json:"create_time" xorm:"created not null comment('创建时间') DATETIME"`

	//更新时间
	UpdateTime time.Time `json:"update_time" xorm:"not null comment('更新时间') DATETIME"`

	//计划终止状态，1：计划有效中，2：计划已终止
	Stop int `json:"stop" xorm:"not null comment('计划终止状态，1：计划有效中，2：计划已终止') index TINYINT(1)"`

	//年级
	Year int `json:"year" xorm:"not null comment('年级') INT(11)"`

	//学期 1 第一学期 2 第二学期
	Term int `json:"term" xorm:"not null comment('学期 1 第一学期 2 第二学期') TINYINT(1)"`

	//学年
	TermYear int `json:"term_year" xorm:"not null comment('学年') INT(11)"`

	//1为行政班，2为运动班
	PlanType int `json:"plan_type" xorm:"not null comment('1为行政班，2为运动班') INT(11)"`

	PlanTimeFrame [][]string `json:"plantimeframe" xorm:"-"`

	//****************增加性别字段
	Gender int `json:"gender" xorm:"gender"`
}

// 响应结构体
//
// swagger:response    PlanCreateResponse
type PlanCreateResponse struct {
	// in: body
	Body planresponseMessage
}
type planresponseMessage struct {
	// Required: true
	models.ResponseType
	Data models.Plan
}

type resplan struct {
	//计划id
	Id int `json:"id" xorm:"not null pk autoincr comment('计划Id') INT(11)"`

	//院系id
	DepartmentId int `json:"department_id" xorm:"not null comment('院系id') index INT(5)"`

	//班级id
	ClassId int `json:"class_id" xorm:"not null comment('班级id') index INT(5)"`

	//计划名称
	Name string `json:"name" xorm:"not null comment('计划名称') VARCHAR(255)"`

	//计划开始日期
	DateBegin time.Time `json:"date_begin" xorm:"not null comment('计划开始日期') DATETIME"`

	//计划结束日期
	DateEnd time.Time `json:"date_end" xorm:"not null comment('计划结束日期') DATETIME"`

	//总里程
	TotalDistance int `json:"total_distance" xorm:"not null comment('总里程') INT(4)"`

	//总次数
	TotalTimes int `json:"total_times" xorm:"not null comment('总次数') INT(4)"`

	//最低周跑次数
	MinWeekTimes int `json:"min_week_times" xorm:"not null comment('最低周跑次数') INT(2)"`

	//最低周跑里程
	MinWeekDistance int `json:"min_week_distance" xorm:"not null comment('最低周跑里程') INT(11)"`

	//最低单次里程
	MinSingleDistance int `json:"min_single_distance" xorm:"not null comment('最低单次里程') INT(11)"`

	//最高周跑次数
	MaxWeekTimes int `json:"max_week_times" xorm:"not null comment('最高周跑次数') INT(11)"`

	//最高周跑里程
	MaxWeekDistance int `json:"max_week_distance" xorm:"not null comment('最高周跑里程') INT(11)"`

	//最高单次里程
	MaxSingleDistance int `json:"max_single_distance" xorm:"not null comment('最高单次里程') INT(11)"`

	//最高日跑次数
	MaxDayTimes int `json:"max_day_times" xorm:"not null comment('最高日跑次数') INT(2)"`

	//必跑日
	MustRunDay string `json:"must_run_day" xorm:"not null comment('必跑日') VARCHAR(20)"`

	//最小配速(单位s)
	MinPace int `json:"min_pace" xorm:"not null comment('最小配速(单位s)') INT(11)"`

	//最大配速(单位s)
	MaxPace int `json:"max_pace" xorm:"not null comment('最大配速(单位s)') INT(11)"`

	//最短时间(单位s)
	MinTimeLong int `json:"min_time_long" xorm:"not null comment('最短时间(单位s)') INT(11)"`

	//最长时间(单位s)
	MaxTimeLong int `json:"max_time_long" xorm:"not null comment('最长时间(单位s)') INT(11)"`

	//创建时间
	CreateTime time.Time `json:"create_time" xorm:"not null comment('创建时间') DATETIME"`

	//更新时间
	UpdateTime time.Time `json:"update_time" xorm:"not null comment('更新时间') DATETIME"`

	//计划终止状态，1：计划有效中，2：计划已终止
	Stop int `json:"stop" xorm:"not null comment('计划终止状态，1：计划有效中，2：计划已终止') index TINYINT(1)"`

	//年级
	Year int `json:"year" xorm:"not null comment('年级') INT(11)"`

	//学期 1 第一学期 2 第二学期
	Term int `json:"term" xorm:"not null comment('学期 1 第一学期 2 第二学期') TINYINT(1)"`

	//学年
	TermYear int `json:"term_year" xorm:"not null comment('学年') INT(11)"`

	//1为行政班，2为运动班
	PlanType int `json:"plan_type" xorm:"not null comment('1为行政班，2为运动班') INT(11)"`

	PlanTimeFrame []string `json:"plantimeframe" xorm:"-"`

	//获取院系名称
	DepartmentName []string `json:"department_name" xorm:"-"`
	//获取班级名称
	ClassName []string `json:"class_name" xorm:"-"`

	//****************增加性别字段
	Gender int `json:"gender" xorm:"gender"`
}

//func create(ctx iris.Context) {
//	// swagger:route POST /api/plan plan PlanCreateRequest
//	//
//	// 创建跑步计划
//	//     Consumes:
//	//     - application/json
//	//
//	//     Produces:
//	//     - application/json
//	//
//	//     Responses:
//	//       200: PlanCreateResponse
//	plan := planandtime{}
//
//	//解析plan
//	err := ctx.ReadJSON(&plan)
//	if err != nil {
//		ctx.JSON(lib.NewResponseFail(1, err.Error()))
//		return
//	}
//
//	//TODO 验证数据有效性
//
//	//为了修复学生每周计划进度，必须将计划终止日期修改为当天23：59：59
//	finDateEndUnix := plan.DateEnd.Unix() + 24*3600 - 1
//	finDateEnd := time.Unix(int64(finDateEndUnix), 0)
//	plan.DateEnd = finDateEnd
//
//	//插入数据
//	res, err2 := lib.Engine.Table("plan").Insert(&plan)
//	if err2 != nil {
//		ctx.JSON(lib.NewResponseFail(1, err.Error()))
//		return
//	}
//	planid := plan.Id
//
//	planTimeFrames := plan.PlanTimeFrame
//	planTimeFramestruct := []models.PlanTimeFrame{}
//	for i, planTimeFrame := range planTimeFrames {
//		planTimeFramestruct[i].PlanId = planid
//		planTimeFramestruct[i].DurationBegin = planTimeFrame[0]
//		planTimeFramestruct[i].DurationEnd = planTimeFrame[1]
//	}
//	resrframe, errframe := lib.Engine.Table("plan_time_frame").Insert(&planTimeFramestruct)
//	if errframe != nil {
//		ctx.JSON(lib.NewResponseFail(1, err.Error()))
//		return
//	}
//	fmt.Println("\n\n\n\n共插入了", resrframe, "条记录\n\n\n")
//
//	ctx.JSON(lib.NewResponseOK(res))
//}

func create(ctx iris.Context) {
	// swagger:route POST /api/plan plan PlanCreateRequest
	//
	// 创建跑步计划
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: PlanCreateResponse

	plan := planandtime{}
	//解析plan
	err := ctx.ReadJSON(&plan)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//TODO 验证数据有效性

	//为了修复学生每周计划进度，必须将计划终止日期修改为当天23：59：59
	finDateEndUnix := plan.DateEnd.Unix() + 24*3600 - 1
	finDateEnd := time.Unix(int64(finDateEndUnix), 0)
	plan.DateEnd = finDateEnd

	querystudent := lib.Engine.Table("student").Where("1=?", 1)
	querystudent1 := lib.Engine.Table("student").Where("1=?", 1)

	students := []models.Student{}
	studentsid := []models.Student{}

	year := ctx.URLParamIntDefault("year", 0)

	departmentid := ctx.URLParamIntDefault("department_id", 0)
	classid := ctx.URLParamIntDefault("class_id", 0)
	gender := ctx.URLParamIntDefault("gender", 0)
	//updatesql :="UPDATE `student` SET `plan_id` = 100 WHERE `id` in"

	if year > 0 {
		querystudent.And("year=?", year)
		querystudent1.And("year=?", year)
	}
	if year == 0 {
		ctx.JSON(lib.NewResponseFail(1, "年份不存在"))

	}
	if departmentid != 0 {
		//判断是否为运动班，运动班设置指定查找的表student_class
		if plan.PlanType == 2 {
			querystudent.And("student_class.department_id=?", departmentid)
			querystudent1.And("student_class.department_id=?", departmentid)

		} else {

			querystudent.And("department_id=?", departmentid)
			querystudent1.And("department_id=?", departmentid)
		}
	}

	if classid != 0 {
		//判断是否为运动班，运动班设置指定查找的表student_class
		if plan.PlanType == 2 {
			querystudent.And("student_class.class_id=?", classid)
			querystudent1.And("student_class.class_id=?", classid)

		} else {
			querystudent.And("class_id=?", classid)
			querystudent1.And("class_id=?", classid)

		}

	}

	if gender != 0 {
		querystudent.And("gender=?", gender)
		querystudent1.And("gender=?", gender)
	}
	//
	//if plan.PlanType==2 {
	//	errcount := querystudent.Join("INNER","student_class","student.id=student_class.student_id").Cols("student_class.student_id", "plan_id").GroupBy("plan_id").Find(&students)
	//	if errcount != nil {
	//		fmt.Printf("查询需要修改的学生范围内包括几个计划错误：%v", errcount)
	//		ctx.JSON(lib.NewResponseFail(1, errcount.Error()))
	//		return
	//	}
	//	//查询影响范围内有多少学生,查出id和plan_id
	//	iderr := querystudent.Join("INNER","student_class","student.id=student_class.student_id").Cols("student_class.student_id", "plan_id").Find(&students)
	//	if iderr != nil {
	//		fmt.Printf("查询学生计划id错误：%v", errcount)
	//		ctx.JSON(lib.NewResponseFail(1, errcount.Error()))
	//		return
	//	}
	//
	//
	//
	//}

	//判断多计划是以行政班还是体育班为基础
	if plan.PlanType == 2 {
		errcount := querystudent.Cols("student.id", "plan_id").Join("INNER", "student_class", "student.id=student_class.student_id").GroupBy("plan_id").Find(&students)
		if errcount != nil {
			fmt.Printf("查询需要修改的学生范围内包括几个计划错误：%v", errcount)
			ctx.JSON(lib.NewResponseFail(1, errcount.Error()))
			return
		}
		//查询影响范围内有多少学生,查出id和plan_id
		iderr := querystudent1.Cols("student.id", "plan_id").Join("INNER", "student_class", "student.id=student_class.student_id").Find(&studentsid)
		if iderr != nil {
			fmt.Printf("查询学生计划id错误：%v", errcount)
			ctx.JSON(lib.NewResponseFail(1, errcount.Error()))
			return
		}

	} else {

		//查询需要修改的学生范围内包括几个计划\

		errcount := querystudent.Cols("id", "plan_id").GroupBy("plan_id").Find(&students)
		if errcount != nil {
			fmt.Printf("查询需要修改的学生范围内包括几个计划错误：%v", errcount)
			ctx.JSON(lib.NewResponseFail(1, errcount.Error()))
			return
		}
		//查询影响范围内有多少学生,查出id和plan_id
		iderr := querystudent1.Cols("id", "plan_id").Find(&studentsid)
		if iderr != nil {
			fmt.Printf("查询学生计划id错误：%v", errcount)
			ctx.JSON(lib.NewResponseFail(1, errcount.Error()))
			return
		}
	}

	var id []int

	for _, student := range studentsid {
		id = append(id, student.Id)
	}
	fmt.Println("\n\n\n\n", id)
	if len(id) == 0 {
		ctx.JSON(lib.NewResponseFail(1, "创建计划失败，该计划内无学生"))
		return
	}

	//查询计划是否已终止，或存在
	//查询计划是否已终止，或存在
	if len(students) > 0 {
		for i, student := range students {
			if student.PlanId != 0 {
				oldPlan := models.Plan{}
				redPlan, errPlan := lib.Engine.Table("plan").Where("id=?", students[i].PlanId).Where("stop=?", 1).Get(&oldPlan)
				if errPlan != nil {
					fmt.Sprintf("查询计划错误：%v", errPlan)
					return
				}
				if redPlan {
					println("创建计划失败")
					ctx.JSON(lib.NewResponse(1, "不能创建该计划，该计划与学校其他计划冲突", 202))
					return

				}

			}
		}

	}
	//for _, student := range students {
	//
	//	//todo 已终止的计划应该允许覆盖
	//	if  redPlan {
	//		println("创建计划失败")
	//		ctx.JSON(lib.NewResponse(1, "不能创建该计划，该计划与学校其他计划冲突", 202))
	//		return
	//
	//	}
	//}
	//插入数据
	res1, err2 := lib.Engine.Table("plan").Insert(&plan)
	if err2 != nil {
		fmt.Printf("插入计划错误：%v", err2)
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	student := models.Student{
		PlanId: plan.Id,
	}
	println("开始更新，学生plan_id")
	res, errupdate := querystudent.ID(id).Cols("plan_id").Join("INNER", "student_class", "student.id=student_class.student_id").Update(student)
	if errupdate != nil {
		fmt.Printf("更新学生plan_id错误：%v", errupdate)
		ctx.JSON(lib.NewResponseFail(1, errupdate.Error()))
		return
	}

	//b2, err := json.MarshalIndent(plan, "", "   ")
	//if err != nil {
	//	fmt.Println("json err:", err)
	//}
	//fmt.Println("\n\n\nplan结构体\n\n", string(b2))
	//if err != nil {
	//	ctx.JSON(lib.NewResponseFail(1, err.Error()))
	//	return
	//}

	planid := plan.Id

	planTimeFrames := plan.PlanTimeFrame
	n := len(planTimeFrames)

	//planTimeFramestruct := []models.PlanTimeFrame{}
	planTimeFramestruct := make([]models.PlanTimeFrame, n)

	for i, planTimeFrame := range planTimeFrames {

		planTimeFramestruct[i].PlanId = planid
		planTimeFramestruct[i].DurationBegin = planTimeFrame[0]
		planTimeFramestruct[i].DurationEnd = planTimeFrame[1]

	}
	resrframe, errframe := lib.Engine.Table("plan_time_frame").Insert(&planTimeFramestruct)
	if errframe != nil {
		fmt.Printf("添加计划时间段：%v", errframe)
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	fmt.Println("\n\n\n\n共插入了", resrframe, "条时间记录\n\n\n")

	lib.NewResponseOK(res1)
	lib.NewResponseOK(res)
	ctx.JSON(lib.NewResponseOK("计划创建成功"))

}

// swagger:parameters  PlanUpdateRequest
type PlanUpdateRequest struct {
	// in: body
	Body models.Plan
}

// 响应结构体
//
// swagger:response    PlanUpdateResponse
type PlanUpdateResponse struct {
	// in: body
	Body planresponseMessage
}

func remove(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	//plan := models.Plan{}

	//根据id查询
	affected, err := lib.Engine.Exec("update plan set stop=2 where id=?", id)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	ctx.JSON(lib.NewResponseOK(affected))
}

// swagger:route GET /api/plan/:id  plan PlanGet
//
// 查找计划信息
//     Produces:
//     - application/json
//
//     Responses:
//       200: Response
func get(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	plan := resplan{}
	//根据id查询
	b, err := lib.Engine.Table("plan").ID(id).Get(&plan)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该计划"))
		return
	}

	////获取gender
	//students := []models.Student{}
	//errGender := lib.Engine.Table("student").Where("plan_id=?", id).Find(&students)
	//if errGender != nil {
	//	fmt.Printf("%v", errGender)
	//	ctx.JSON(lib.NewResponseFail(1, "查询学生错误"))
	//	return
	//}
	//
	//genderMan := false
	//genderMen := false
	//for _, value := range students {
	//	if value.Gender == 1 {
	//		genderMan = true
	//	}
	//	if value.Gender == 2 {
	//		genderMen = true
	//	}
	//
	//}
	//
	//genderType := 0
	//if genderMan {
	//	genderType = 1
	//}
	//if genderMen {
	//	genderType = 2
	//}
	//if genderMen && genderMan {
	//	genderType = 0
	//}
	//plan.Gender = genderType

	//获取跑步规定时间段
	timeFrame := []models.PlanTimeFrame{}
	errTimeFrame := lib.Engine.Table("plan_time_frame").Where("plan_id=?", id).Find(&timeFrame)
	if errTimeFrame != nil {

		ctx.JSON(lib.NewResponseFail(1, "查询学生错误"))
		return
	}

	timeFrameStr := make([]string, 0)
	if len(timeFrame) > 0 {
		for _, value := range timeFrame {
			timeFrameStr = append(timeFrameStr, value.DurationBegin)
			timeFrameStr = append(timeFrameStr, value.DurationEnd)
		}
	}

	plan.PlanTimeFrame = timeFrameStr

	departments := make([]string, 0)
	classes := make([]string, 0)
	//获取院系名称
	if plan.DepartmentId == 0 {
		departments = append(departments, "")
	} else {
		department := models.Department{}
		res, err := lib.Engine.Table("department").Where("id=?", plan.DepartmentId).Get(&department)
		if err != nil {
			fmt.Printf("%v", err)
			ctx.JSON(lib.NewResponseFail(1, "查询院系错误"))
			return
		}

		if res == false {
			println("查询院系失败，该院系id不存在。")
		}
		departments = append(departments, department.Name)
	}

	//获取班级名称
	if plan.ClassId == 0 {
		classes = append(classes, "")
	} else {
		class := models.Classes{}
		res, err := lib.Engine.Table("classes").Where("id=?", plan.ClassId).Get(&class)
		if err != nil {
			fmt.Printf("%v", err)
			ctx.JSON(lib.NewResponseFail(1, "查询班级错误"))
			return
		}

		if res == false {
			println("查询院系失败，该院系id不存在。")
		}
		classes = append(classes, class.Name)
	}
	plan.DepartmentName = departments
	plan.ClassName = classes

	ctx.JSON(lib.NewResponseOK(plan))
}

// swagger:route GET /api/plans plan PlanSearch
//
// 查找多条计划信息
//     Produces:
//     - application/json
//
//     Responses:
//       200: Response
func search(ctx iris.Context) {

	//创建查询Session
	query := lib.Engine.Table("plan")

	//字段查询
	if ctx.URLParamExists("name") {
		query.And(builder.Like{"name", ctx.URLParam("name")})
	}
	if ctx.URLParamExists("stop") {
		stop, errstop := ctx.URLParamInt("stop")
		if errstop != errstop {
			ctx.JSON(lib.NewResponseFail(1, errstop.Error()))
		}

		if stop == 1 {
			query.And("stop=?", stop)

		}
		if stop == 2 {

			query.And("stop=?", stop)
		}

	}

	////已终止计划不返回
	//query.And("stop=?", 1)

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
	var plan []resplan
	err := query.Find(&plan)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//获取所有院系
	departments := []models.Department{}
	errDep := lib.Engine.Table("department").Find(&departments)
	if errDep != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, "查询院系错误"))
		return
	}
	//获取所有班级
	classes := []models.Classes{}
	errCla := lib.Engine.Table("classes").Find(&classes)
	if errCla != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, "查询班级错误"))
		return
	}

	//获取计划的时间段和性别
	for index, value := range plan {

		//获取跑步规定时间段
		timeFrame := []models.PlanTimeFrame{}
		errTimeFrame := lib.Engine.Table("plan_time_frame").Where("plan_id=?", value.Id).Find(&timeFrame)
		if errTimeFrame != nil {
			fmt.Printf("查询计划时间段错误：%v", errTimeFrame)
			ctx.JSON(lib.NewResponseFail(1, "查询学生错误"))
			return
		}
		timeFrameStr := make([]string, 0)
		if len(timeFrame) > 0 {
			for _, value := range timeFrame {
				timeFrameStr = append(timeFrameStr, value.DurationBegin)
				timeFrameStr = append(timeFrameStr, value.DurationEnd)
			}
		}

		plan[index].PlanTimeFrame = timeFrameStr

		//添加院系
		for _, valueDep := range departments {
			if value.DepartmentId == valueDep.Id {
				plan[index].DepartmentName = append(plan[index].DepartmentName, valueDep.Name)
			}
		}
		//添加班级
		for _, valueCla := range classes {
			if value.ClassId == valueCla.Id {
				plan[index].ClassName = append(plan[index].ClassName, valueCla.Name)
			}
		}
	}

	ctx.JSON(lib.NewResponseOK(plan))
}
func findclasstype(ctx iris.Context) {
	query := lib.Engine.Table("plan").Where("stop=1")
	plan := models.Plan{}
	plansum, err := query.SumInt(plan, "stop")
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	if plansum == 0 {
		ctx.JSON(lib.NewResponse(0, "该学校目前无计划运行", 0))
		return
	} else {
		b, err1 := query.Limit(1, 0).Cols("plan_type").Get(&plan)
		if err1 != nil {
			ctx.JSON(lib.NewResponseFail(1, err.Error()))
			return
		}
		if b == false {
			ctx.JSON(lib.NewResponseFail(1, "无正在运行的计划"))
			return
		}

		ctx.JSON(lib.NewResponse(1, "该学校有正在运行的计划", plan.PlanType))
		return
	}
}

func newUpdate(ctx iris.Context) {

	//1.获取旧计划和新计划
	plan := planandtime{}

	//解析plan
	err := ctx.ReadJSON(&plan)

	//为了修复学生每周计划进度，必须将计划终止日期修改为当天23：59：59
	dateEndStr := plan.DateEnd.Format("2006-01-02")
	timeEnd, _ := time.ParseInLocation("2006-01-02 15:04:05", dateEndStr+" 23:59:59", time.Local)
	plan.DateEnd = timeEnd

	oldstudent := lib.Engine.Table("student").Where("1=1")
	newstudent := lib.Engine.Table("student").Where("1=1")
	newstudent1 := lib.Engine.Table("student").Where("1=1")

	oldplanid := ctx.URLParamIntDefault("old_plan_id", 0)

	year := ctx.URLParamIntDefault("old_year", 0)

	departmentid := ctx.URLParamIntDefault("old_department_id", 0)
	classid := ctx.URLParamIntDefault("old_class_id", 0)
	gender := ctx.URLParamIntDefault("old_gender", 0)
	//updatesql :="UPDATE `student` SET `plan_id` = 100 WHERE `id` in"

	newyear := ctx.URLParamIntDefault("new_year", 0)
	newdepartmentid := ctx.URLParamIntDefault("new_department_id", 0)
	newclassid := ctx.URLParamIntDefault("new_class_id", 0)
	newgender := ctx.URLParamIntDefault("new_gender", 0)

	//2.获取新计划的覆盖的学生的计划
	if oldplanid == 0 {
		fmt.Printf("找不到计划id:")
		ctx.JSON(lib.NewResponseFail(1, "需要修改的计划不存在，未传入原计划id"))
		return
	}
	if year > 0 {
		oldstudent.And("student.year=?", year)
	}
	if year == 0 {
		ctx.JSON(lib.NewResponseFail(1, "年份不存在"))

	}

	if departmentid != 0 {
		if plan.PlanType == 2 {
			oldstudent.And("student_class.department_id=?", departmentid)
			oldstudent.And("student_class.department_id=?", departmentid)
		} else {

			oldstudent.And("student.department_id=?", departmentid)
		}
	}
	if classid != 0 {
		if plan.PlanType == 2 {
			oldstudent.And("student_class.class_id=?", classid)
		} else {
			oldstudent.And("student.class_id=?", classid)
		}
	}
	if gender != 0 {
		oldstudent.And("student.gender=?", gender)
	}

	//获取新计划学生范围
	if newyear > 0 {
		newstudent.And("student.year=?", newyear)
		newstudent1.And("student.year=?", newyear)
	} else if newyear == 0 {
		ctx.JSON(lib.NewResponseFail(1, "年份不存在"))

	}

	if newdepartmentid != 0 {
		if plan.PlanType == 2 {
			newstudent.And("student_class.department_id=?", newdepartmentid)
			newstudent1.And("student_class.department_id=?", newdepartmentid)

		} else {
			newstudent.And("student.department_id=?", newdepartmentid)
			newstudent1.And("student.department_id=?", newdepartmentid)
		}
	}

	if newclassid != 0 {
		if plan.PlanType == 2 {
			newstudent.And("student_class.class_id=?", newclassid)
			newstudent1.And("student_class.class_id=?", newclassid)

		} else {
			newstudent.And("student.class_id=?", newclassid)
			newstudent1.And("student.class_id=?", newclassid)
		}
	}

	if newgender != 0 {
		newstudent.And("student.gender=?", newgender)
		newstudent1.And("student.gender=?", newgender)
	}

	studentstag2 := []models.Student{}
	//新学生的id Join("INNER", "classes", "classes.id=student.class_id").
	if plan.PlanType == 2 {
		newerrcount := newstudent1.Join("INNER", "plan", "plan.id=student.plan_id").Join("INNER", "student_class", "student.id=student_class.student_id").
			And("plan.stop=1").And("student.plan_id!=?", oldplanid).Cols("student.id", "student.plan_id").Find(&studentstag2)
		if newerrcount != nil {
			fmt.Printf("查询学生计划错误：%v", newerrcount)
			ctx.JSON(lib.NewResponseFail(1, newerrcount.Error()))
			return
		}

	} else {
		newerrcount := newstudent1.Join("INNER", "plan", "plan.id=student.plan_id").
			And("plan.stop=1").And("student.plan_id!=?", oldplanid).Cols("student.id", "student.plan_id").Find(&studentstag2)
		if newerrcount != nil {
			fmt.Printf("查询学生计划错误：%v", newerrcount)
			ctx.JSON(lib.NewResponseFail(1, newerrcount.Error()))
			return
		}
	}
	//当新计划范围下的学生已有计划正在运行且不是旧计划，返回，不允许修改
	if len(studentstag2) != 0 {
		println("新计划覆盖的学生有正在运行的计划！")
		ctx.JSON(lib.NewResponseFail(1, "新计划覆盖的学生有正在运行的计划"))
		return
	} else {
		fmt.Println("更新计划表之前，查询新计划内是否存在学生")

		//***************************更新计划表之前，查询新计划内是否存在学生**************

		studentstag1 := []models.Student{}
		//新学生的id
		newerrcount := newstudent.Cols("student.id", "plan_id").Join("INNER", "student_class", "student.id=student_class.student_id").Find(&studentstag1)
		if newerrcount != nil {
			ctx.JSON(lib.NewResponseFail(1, newerrcount.Error()))
			return
		}
		var newid []int
		for _, student := range studentstag1 {
			newid = append(newid, student.Id)
		}

		if len(newid) == 0 {
			ctx.JSON(lib.NewResponseFail(1, "创建计划失败，该计划内无学生"))
			return
		}

		//***************************更新计划表之前，查询新计划内是否存在学生  完成**************

		fmt.Println("\n\n\n1.开始更新计划表\n\n\n")
		//更新计划表
		res1, err2 := lib.Engine.Table("plan").AllCols().ID(oldplanid).Update(&plan)
		if err2 != nil {
			fmt.Printf("插入计划错误：%v", err2)
			ctx.JSON(lib.NewResponseFail(1, err.Error()))
			return
		}
		fmt.Printf("计划修改的条数：%d", res1)

		students := []models.Student{}
		fmt.Println("\n\n\n2.根据原计划查出影响的学生id\n\n\n")
		//根据原计划查出影响的学生id
		errcount := oldstudent.Cols("student.id").Join("INNER", "student_class", "student.id=student_class.student_id").Find(&students)
		if errcount != nil {
			fmt.Printf("查询旧计划受影响的学生错误：%v", errcount)
			ctx.JSON(lib.NewResponseFail(1, errcount.Error()))
			return
		}

		var id []int
		for _, student := range students {
			id = append(id, student.Id)
		}
		//把原学生的planid更新为0（原来为大范围）
		student := models.Student{
			PlanId: 0,
		}

		//1.置空旧计划影响的学生表plan_id
		println("开始置空旧学生的计划：")
		oldres, olderrupdate := oldstudent.ID(id).Cols("plan_id").Join("INNER", "student_class", "student.id=student_class.student_id").Update(student)
		if olderrupdate != nil {
			fmt.Printf("置空旧学生计划错误：%v", olderrupdate)
			ctx.JSON(lib.NewResponseFail(1, olderrupdate.Error()))
			return
		}

		if oldres == 0 {
			println("更新旧学生的planId影响的数量：0")
			//ctx.JSON(lib.NewResponseFail(1, "更新旧学生的planId失败"))
			//return
		}

		//2.把原计划储存的时间段删除
		println("开始删除旧计划时间段：")
		plantimeframe := models.PlanTimeFrame{}
		olddeleteres, olddeleteerr := lib.Engine.Table("plan_time_frame").Where("plan_id=?", oldplanid).Delete(plantimeframe)
		if olddeleteerr != nil {
			ctx.JSON(lib.NewResponseFail(1, olddeleteerr.Error()))
			return
		}

		if olddeleteres == 0 {
			println("删除旧计划的时间段失败，原时间段不存在")
			//ctx.JSON(lib.NewResponseFail(1, "删除旧计划的时间段失败"))
			//return
		}

		//3.更新学生表plan_id
		//studentstag1 := []models.Student{}
		////新学生的id
		//newerrcount := newstudent.Cols("id", "plan_id").Find(&studentstag1)
		//if newerrcount != nil {
		//	ctx.JSON(lib.NewResponseFail(1, newerrcount.Error()))
		//	return
		//}
		//var newid []int
		//for _, student := range studentstag1 {
		//	newid = append(newid, student.Id)
		//}

		newstudentstruct := models.Student{
			PlanId: plan.Id,
		}

		//更新影响的学生表的Plan_id
		println("开始更新，新的学生plan_id")

		res, errupdate := newstudent.Join("INNER", "student_class", "student.id=student_class.student_id").ID(newid).Update(newstudentstruct)
		if errupdate != nil {
			ctx.JSON(lib.NewResponseFail(1, errupdate.Error()))
			return
		}

		println("受影响的学生：")
		fmt.Printf("\n\n\n**********更新新计划受影响的学生%d\n\n\n", res)

		planid := plan.Id

		planTimeFrames := plan.PlanTimeFrame
		n := len(planTimeFrames)

		//4.插入新时间段
		planTimeFramestruct := make([]models.PlanTimeFrame, n)

		for i, planTimeFrame := range planTimeFrames {

			planTimeFramestruct[i].PlanId = planid
			planTimeFramestruct[i].DurationBegin = planTimeFrame[0]
			planTimeFramestruct[i].DurationEnd = planTimeFrame[1]

		}
		println("插入新的时间段")
		resrframe, errframe := lib.Engine.Table("plan_time_frame").Insert(&planTimeFramestruct)
		if errframe != nil {
			fmt.Printf("插入新时间段错误：%v", errframe)
			ctx.JSON(lib.NewResponseFail(1, err.Error()))
			return
		}

		if resrframe == 0 {
			println("插入新时间段失败")
			ctx.JSON(lib.NewResponseFail(1, "插入新时间段失败"))
			return
		}

		ctx.JSON(lib.NewResponseOK("修改计划成功"))
		return

	}

	//2.1. 如果newStuPlanIds!=oldPlan&&newStuPlanIds is running {return  }

	//3. 允许修改：

	//3.1置空原计划student表覆盖的

}

type planinformationcardstruct struct {
	//总次数
	TotalTimes int `json:"total_times" xorm:"not null comment('总次数') INT(4)"`
	//最低周跑次数
	MinWeekTimes int `json:"min_week_times" xorm:"not null comment('最低周跑次数') INT(2)"`
	//最低单次里程
	MinSingleDistance int `json:"min_single_distance" xorm:"not null comment('最低单次里程') INT(11)"`
	//必跑日
	MustRunDay string `json:"must_run_day" xorm:"not null comment('必跑日') VARCHAR(20)"`
	//总里程
	TotalDistance int `json:"total_distance" xorm:"not null comment('总里程') INT(4)"`
	//最小配速(单位s)
	MinPace int `json:"min_pace" xorm:"not null comment('最小配速(单位s)') INT(11)"`
}

func planinformationcard(ctx iris.Context) {

	//plan := models.Plan{}
	planinformationcard := planinformationcardstruct{}
	id := ctx.URLParamIntDefault("plan_id", 0)
	fmt.Println("\n\n\n\n\n\n", id)

	res, err := lib.Engine.Table("plan").Where("id=?", id).
		Cols("total_times", "min_week_times", "min_single_distance", "total_distance", "must_run_day", "min_pace").Get(&planinformationcard)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
	}
	lib.NewResponseOK(res)

	ctx.JSON(lib.NewResponseOK(planinformationcard))

}

type PlanCard struct {
	Todaydistance int64 `json:"todaydistance"`
	Distance      int64 `json:"distance"`
	Runstudent    int64 `json:"runstudent"`
}

func blackboard(ctx iris.Context) {

	id := ctx.URLParamIntDefault("plan_id", 0)
	planrecord := models.PlanRecord{}

	//b,err :=lib.Engine.Table("plan").Cols("gender","year","class_id","department_id").Get(&plan)

	//i,err :=lib.Engine.Table("plan_record").
	//	Where("gender=?",plan.Gender).
	//	And("year",plan.Year).And("class_id",plan.ClassId).
	//	And("department_id",plan.DepartmentId).
	//	SumsInt(&planrecord,"distance")
	//
	//
	//i,err :=lib.Engine.Table("plan_record").
	//	Where("gender=?",plan.Gender).
	//	And("year",plan.Year).And("class_id",plan.ClassId).
	//计划内今日跑量
	//
	//todaydistance, err := lib.Engine.Table("plan_record").Where("plan_id=?", id).And("to_days(create_at) = to_days(now())").SumInt(&planrecord, "distance")
	//if err != nil {
	//	ctx.JSON(lib.NewResponseFail(1, err.Error()))
	//}
	//计划内总跑量
	distance, err := lib.Engine.Table("plan_record").Where("plan_id=?", id).SumInt(&planrecord, "distance")

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
	}
	//本日跑量
	Todaydistance, err := lib.Engine.Table("plan_record").Where("plan_id=?", id).And("to_days(create_at) = to_days(now())").SumInt(&planrecord, "distance")
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
	}
	var V1string string
	//总跑步人数
	sql := "select count(*) from ( SELECT count(*) FROM `plan_record`"
	if id != 0 {
		sql = fmt.Sprintf("%s WHERE (plan_id=%d)  AND (to_days(create_at) = to_days(now()) ) GROUP BY student_id ) as s", sql, id)

	}

	fmt.Println("sql:", sql)

	runstudentmap, err := lib.Engine.Query(sql)
	for k, v := range runstudentmap {
		fmt.Printf("打印1%d%s", k, v)

		for _, v1 := range v {
			V1string = string(v1)
			fmt.Println("学生本日跑步总计", V1string)
		}
	}

	i64, err := strconv.ParseInt(V1string, 10, 64)
	if err == nil {
		fmt.Printf("i64: %v\n", i64)
	}

	planCard := PlanCard{}
	planCard.Distance = distance
	planCard.Runstudent = i64

	planCard.Todaydistance = Todaydistance

	ctx.JSON(lib.NewResponseOK(planCard))

}

//func allblackboard(ctx iris.Context){
//	type PlanCard struct {
//		Todaydistance int64  `json:"todaydistance"`
//		Distance int64    `json:"distance"`
//		Runstudent int64   `json:"runstudent"`
//
//
//	}
//	type PlanId struct {
//
//	}
//	lib.Engine.Table("plan").Cols("id")
//
//
//	var id []int
//
//	for _, student := range studentsid {
//		id = append(id, student.Id)
//	}
//	planrecord :=models.PlanRecord{}
//
//	//b,err :=lib.Engine.Table("plan").Cols("gender","year","class_id","department_id").Get(&plan)
//
//
//
//	//i,err :=lib.Engine.Table("plan_record").
//	//	Where("gender=?",plan.Gender).
//	//	And("year",plan.Year).And("class_id",plan.ClassId).
//	//	And("department_id",plan.DepartmentId).
//	//	SumsInt(&planrecord,"distance")
//	//
//	//
//	//i,err :=lib.Engine.Table("plan_record").
//	//	Where("gender=?",plan.Gender).
//	//	And("year",plan.Year).And("class_id",plan.ClassId).
//	//计划内今日跑量
//	todaydistance,err :=lib.Engine.Table("plan_record").Where("plan_id=?",id).And("create_at=to_days(now())").SumInt(&planrecord,"distance")
//	if err != nil {
//		ctx.JSON(lib.NewResponseFail(1, err.Error()))
//	}
//	//计划内总跑量
//	distance,err :=lib.Engine.Table("plan_record").Where("plan_id=?",id).SumInt(&planrecord,"distance")
//
//	if err != nil {
//		ctx.JSON(lib.NewResponseFail(1, err.Error()))
//	}
//	//总跑步人数
//	runstudent,err :=lib.Engine.Table("plan_record").Where("plan_id=?",id).GroupBy("student_id").Count(&planrecord)
//	if err != nil {
//		ctx.JSON(lib.NewResponseFail(1, err.Error()))
//	}
//
//	planCard :=PlanCard{}
//	planCard.Distance=distance
//	planCard.Runstudent=runstudent
//	planCard.Todaydistance=todaydistance
//
//
//	ctx.JSON(lib.NewResponseOK(planCard))
//
//
//}
func addprogress(ctx iris.Context) {
	lib.MainLogger.Info("-------------批量更新学生进度成功开始-------------")
	planid1 := ctx.URLParamIntDefault("plan_id", 0)
	studentId := ctx.URLParamIntDefault("student_id", 0)

	plan1 := models.Plan{}
	student2 := []models.Student{}

	//查询计划
	_, err := lib.Engine.Table("plan").Where("id=?", planid1).Get(&plan1)
	if err != nil {
		fmt.Printf("查询学生计划错误：%v", err)
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	if studentId == 0 {
		iderr := lib.Engine.Table("student").Cols("id").Where("plan_id=?", planid1).Find(&student2)
		if iderr != nil {
			fmt.Printf("查询学生计划id错误：%v", iderr)
			ctx.JSON(lib.NewResponseFail(1, iderr.Error()))
			return
		}
		if len(student2) == 0 {
			fmt.Println("共有0个学生")
			return
		}

		for _, student := range student2 {
			run.UpdateStudentProgress(student.Id, plan1)
		}
	} else {
		run.UpdateStudentProgress(studentId, plan1)
	}
	lib.MainLogger.Info("-------------批量更新学生进度成功结束-------------")

}

func TestRecord(ctx iris.Context) {
	recordId := ctx.URLParamIntDefault("record_id", 0)

	key := "finishRunJob"
	lib.MainLogger.Info("模拟跑完，record_id:" + string(recordId))
	err := run.PutJob(recordId, key)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

}

type Doughnutchart struct {
	BoyCount       int `json:"boy_count"`
	GirlCount      int `json:"girl_count"`
	BoySum         int `json:"boy_sum"`
	GirlSum        int `json:"girl_sum"`
	Boycompletion  int `json:"boy_completion"`
	Girlcompletion int `json:"girl_completion"`
}

func doughnutchart(ctx iris.Context) {

	var boyCount int
	var girlCount int
	var boycompletion int
	var girlcompletion int

	zerotime := zerotime()
	zerotime1 := zerotime
	zerotime = zerotime + 24*60*60
	tm1 := time.Unix(zerotime1, 0).Format("2006-01-02 15:04:05")
	tm2 := time.Unix(zerotime, 0).Format("2006-01-02 15:04:05")
	println(tm2)
	println(tm1)

	plan := models.Plan{}
	planid := ctx.URLParamIntDefault("plan_id", 0)
	student := models.Student{}
	//planRecord :=models.PlanRecord{}
	_, err := lib.Engine.Table("plan").Where("id=?", planid).Get(&plan)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return

	}

	//判断本周有多少男生跑步次数合格
	sql := "SELECT count(*) FROM `plan_record` WHERE (status=1) AND(gender=1)"

	sql = fmt.Sprintf("%s AND (plan_id=%d)  AND (to_days(create_at) = to_days(now()) ) GROUP BY student_id", sql, planid)

	sqlmap, err := lib.Engine.Query(sql)

	println("\n\n\n*******111111111*********\n\n\n", plan.MinWeekTimes)

	for _, v := range sqlmap {
		for _, v1 := range v {
			v1, err := strconv.Atoi(string(v1))
			if err != nil {
				println("map值转换为int出错")
			}
			if plan.MinWeekTimes < v1 {
				boyCount++
			}

		}

	}

	fmt.Printf("共有%d名男生完成跑步", boyCount)
	//******************结束***********************

	//判断本周有多少女生跑步次数合格
	sql = "SELECT count(*) FROM `plan_record` WHERE (status=1) AND(gender=2)"

	sql = fmt.Sprintf("%s AND (plan_id=%d)  AND (to_days(create_at) = to_days(now()) ) GROUP BY student_id", sql, planid)

	sqlmap, err = lib.Engine.Query(sql)
	println("\n\n\n*******2222222*********\n\n\n", plan.MinWeekTimes)

	for _, v := range sqlmap {
		for _, v1 := range v {
			v1, err := strconv.Atoi(string(v1))
			if err != nil {
				println("map值转换为int出错")
			}
			if plan.MinWeekTimes < v1 {
				girlCount++
			}

		}

	}

	fmt.Printf("共有%d名女生完成跑步", girlCount)

	//******************结束***********************

	//查询男孩总数
	girlSum, err := lib.Engine.Table("student").Where("gender=2").And("plan_id=?", planid).Count(&student)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return

	}

	//查询女孩总数
	boySum, err := lib.Engine.Table("student").Where("gender=1").And("plan_id=?", planid).Count(&student)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return

	}

	//男生完成人数
	if boySum != 0 {

		boycompletion = (boyCount * 100) / int(boySum)

	}

	//女生完成人数
	if girlSum != 0 {
		girlSumInt := int(girlSum)
		girlcompletion = (girlCount * 100) / girlSumInt

	}

	doughnutChart := Doughnutchart{
		BoyCount:       boyCount,
		GirlCount:      girlCount,
		BoySum:         int(boySum),
		GirlSum:        int(girlSum),
		Boycompletion:  boycompletion,
		Girlcompletion: girlcompletion,
	}

	ctx.JSON(lib.NewResponseOK(doughnutChart))

}
func zerotime() (zerotime int64) {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	//fmt.Println(t.Format(time.UnixDate))
	//Unix返回早八点的时间戳，减去8个小时
	timestamp := t.UTC().Unix() - 8*3600
	//fmt.Println("timestamp:", timestamp)
	return timestamp

}

// 个人页图表

func gethistogram(ctx iris.Context) {

	var (
		gender       int
		departmentid int
		month        int
		year         int
		studentid    int
	)

	sumChan := make(chan Sorts, 49)

	//totalfrequency := make([]int, 0)

	zerotime := zerotime()

	totalfrequencystruct := []Sorts{}

	count := 1
	counts := 1

	if ctx.URLParamExists("department_id") {
		departmentid = int(ctx.URLParamInt64Default("departmentid_id", 0))
		//query :=lib.Engine.Table("plan_record").Where("departmentid_id",departmentid)
	}

	if ctx.URLParamExists("gender") {
		//gender := ctx.Params().GetUint64Default("gender",0)
		gender = int(ctx.URLParamInt64Default("gender", 0))

		//fmt.Println("性别", gender)
	}
	if ctx.URLParamExists("month") {
		month = int(ctx.URLParamInt64Default("month", 0))

	}
	if ctx.URLParamExists("year") {
		year = int(ctx.URLParamInt64Default("year", 0))

	}
	if ctx.URLParamExists("student_id") {
		studentid = int(ctx.URLParamInt64Default("student_id", 0))

	}

	//输入到管道缓存
	for i := 0; i < 12; i++ {

		zerotime1 := zerotime
		zerotime = zerotime + 60*60*2
		tm1 := time.Unix(zerotime1, 0).Format("2006-01-02 15:04:05")
		tm2 := time.Unix(zerotime, 0).Format("2006-01-02 15:04:05")
		counts1 := counts
		counts++

		go func() {

			j, sum, err := findRecord(counts1, gender, departmentid, tm1, tm2, month, year, studentid)

			//log.Println("sum:", sum)
			s := Sorts{}
			s.Num = j
			s.Sum = sum
			if err != nil {
				log.Println("err:", err)
			}

			sumChan <- s
			return //这里一定要写return  否则协程不会自动结束，导致数据库连接不会释放
		}()
	}

	//totalfrequency = make([]int, 0)
	for value := range sumChan {
		if count == 12 {
			log.Println("1111111周", value)
			totalfrequencystruct = append(totalfrequencystruct, value)
			close(sumChan)
			break

		}

		log.Println(value)
		totalfrequencystruct = append(totalfrequencystruct, value)
		count++
	}
	fmt.Println(totalfrequencystruct)

	ctx.JSON(lib.NewResponseOK(totalfrequencystruct))
	return
}

func findRecord(counts int, gender int, departmentid int, tm1, tm2 string, month, year int, studentid int) (j int, sum int, err error) {
	//log.Println(index, tm1, tm2, gender)
	record := new(models.PlanRecord)
	session := lib.Engine.NewSession()
	defer session.Close()

	if gender != 0 {
		session.And("gender = ?", gender)
	}
	if departmentid != 0 {
		session.And("department_id = ?", departmentid)
	}

	newtm1 := tm1
	newtm2 := tm2
	//给定月份，返回本月每天固定时段的跑步次数总和
	if month == 1 {
		now := time.Now()
		currentYear, currentMonth, _ := now.Date()
		currentLocation := now.Location()

		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1).Day()
		//时分秒
		var buffer3 bytes.Buffer
		var buffer4 bytes.Buffer
		hour1 := tm1[11:len(tm1)]
		buffer3.WriteString(hour1)
		newtm1 = buffer3.String()
		hour2 := tm2[11:len(tm2)]
		buffer4.WriteString(hour2)
		newtm2 = buffer4.String()
		//月份
		var buffer bytes.Buffer
		monthtm1 := tm1[0:8]
		//ltm1 := tm1[10:len(tm1)]
		buffer.WriteString(monthtm1)
		buffer.WriteString("01")
		//buffer.WriteString(ltm1)

		var buffer1 bytes.Buffer
		monthtm2 := tm2[0:8]
		//ltm2 := tm2[10:len(tm2)]
		buffer1.WriteString(monthtm2)
		buffer1.WriteString(strconv.Itoa(lastOfMonth))
		//buffer1.WriteString(ltm2)
		tm1 = buffer.String()
		tm2 = buffer1.String()

		//session.Exec("SELECT COALESCE(sum(`status`),0) FROM `plan_record` WHERE (DATE_FORMAT(end_time,'%H:%i:%S') > '15:00:00'  and DATE_FORMAT(end_time,'%H:%i:%S') <= '18:00:00')")
		monthsum, err1 := session.Table("plan_record").
			And("student_id=?", studentid).
			And("DATE_FORMAT(end_time,'%H:%i:%S') > ?", newtm1).
			And("DATE_FORMAT(end_time,'%H:%i:%S')<= ?", newtm2).And("end_time > ?", tm1).
			And("end_time <= ?", tm2).
			Sum(record, "status")
		if err1 != nil {
			return 0, 0, err1
		}
		sum = int(monthsum)
		//log.Printf("%v,sum:%v.\n", index, sum)
		j = counts
		session.Close()
		return

	}
	if year == 1 {
		//时分秒
		var buffer3 bytes.Buffer
		var buffer4 bytes.Buffer
		hour1 := tm1[11:len(tm1)]
		buffer3.WriteString(hour1)
		newtm1 = buffer3.String()
		hour2 := tm2[11:len(tm2)]
		buffer4.WriteString(hour2)
		newtm2 = buffer4.String()

		//年份
		var buffer5 bytes.Buffer
		ftm1 := tm1[0:5]
		buffer5.WriteString(ftm1)
		buffer5.WriteString("01-01")
		var buffer6 bytes.Buffer
		ftm2 := tm2[0:5]
		buffer6.WriteString(ftm2)
		buffer6.WriteString("12-31")
		tm1 = buffer5.String()
		tm2 = buffer6.String()

		//session.Exec("SELECT COALESCE(sum(`status`),0) FROM `plan_record` WHERE (DATE_FORMAT(end_time,'%H:%i:%S') > '15:00:00'  and DATE_FORMAT(end_time,'%H:%i:%S') <= '18:00:00')")
		yearsum, err2 := session.Table("plan_record").
			And("DATE_FORMAT(end_time,'%H:%i:%S') > ?", newtm1).
			And("student_id=?", studentid).
			And("end_time > ?", tm1).
			And("end_time <= ?", tm2).
			And("DATE_FORMAT(end_time,'%H:%i:%S')<= ?", newtm2).
			Sum(record, "status")
		if err2 != nil {
			return 0, 0, err2
		}
		sum = int(yearsum)
		//log.Printf("%v,sum:%v.\n", index, sum)
		j = counts
		session.Close()
		return

	}

	sum2, err := session.Table("plan_record").
		And("student_id=?", studentid).
		And("end_time > ?", newtm1).
		And("end_time<= ?", newtm2).
		Sum(record, "status")
	if err != nil {
		return 0, 0, err
	}
	sum = int(sum2)
	//log.Printf("%v,sum:%v.\n", index, sum)
	j = counts
	return

}

type Sorts struct {
	Num int //自己加的管道tag标记
	Sum int //每次根据条件查出的跑步次数
}

type ChartArr struct {
	Num  [7]int64  `json:"num"`
	Data [7]string `json:"data"`
}

//运动趋势
func movementtrend(ctx iris.Context) {
	type Chart struct {
		Num  int64  `json:"num"`
		Data string `json:"data"`
	}

	ChartArr1 := ChartArr{}

	var ChartString [7]Chart

	id := ctx.URLParamIntDefault("plan_id", 0)

	var i int64
	i = 1
	J := 0
	for i = 1; i < 8; i++ {

		zerotime := zerotime()
		zerotime1 := zerotime - (i-1)*24*60*60
		zerotime = zerotime - i*24*60*60
		tm1 := time.Unix(zerotime1, 0).Format("2006-01-02 15:04:05")
		tm2 := time.Unix(zerotime, 0).Format("2006-01-02 15:04:05")
		println(i, tm2)
		println(tm1)

		var V1string string
		//总跑步人数
		sql := "select count(*) from ( SELECT count(*) FROM `plan_record`"
		if id != 0 {
			sql = fmt.Sprintf("%s WHERE (plan_id=%d)   AND (end_time > %q) AND (end_time<= %q) GROUP BY student_id ) as s", sql, id, tm2, tm1)

		}

		fmt.Println("sql:", sql)

		runstudentmap, err := lib.Engine.Query(sql)
		for k, v := range runstudentmap {
			fmt.Printf("打印1%d%s", k, v)

			for _, v1 := range v {
				V1string = string(v1)
				fmt.Println("学生本日跑步总计", V1string)
			}
		}

		i64, err := strconv.ParseInt(V1string, 10, 64)
		if err == nil {
			fmt.Printf("i64: %v\n", i64)
		}

		println("*************", tm1)

		tmString1 := tm1[5:7]
		tmString2 := tm1[8:10]
		var buffer bytes.Buffer
		buffer.WriteString(tmString1)
		buffer.WriteString("/")
		buffer.WriteString(tmString2)
		tmString := buffer.String()

		println(tmString)

		ChartString[J].Data = tmString
		ChartString[J].Num = i64
		ChartArr1.Num[J] = i64
		ChartArr1.Data[J] = tmString
		J++

	}

	ctx.JSON(lib.NewResponseOK(ChartArr1))
}
