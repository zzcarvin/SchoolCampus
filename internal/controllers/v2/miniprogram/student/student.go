package student

import (
	"Campus/configs"
	"Campus/internal/lib"
	"Campus/internal/models"
	"fmt"
	"github.com/kataras/iris"
	"strings"
	"time"
)

type studentRequest struct {
	Name string `json:"name"`
	Code int    `json:"code"`
}

type studentResponse struct {
}

type responseRecord struct {
	Id           int             `json:"id" xorm:"autoincr id"`
	SchoolId     int             `json:"school_id" xorm:"school_id"`
	PlanId       int             `json:"plan_id" xorm:"plan_id"`
	Name         string          `json:"name" xorm:"name"`
	StudentId    int             `json:"student_id" xorm:"student_id"`
	Type         int             `json:"type" xorm:"type"`
	StartTime    string          `json:"start_time" xorm:"start_time"`
	EndTime      string          `json:"end_time" xorm:"end_time"`
	Distance     int             `json:"distance" xorm:"distance"`
	Duration     int             `json:"duration" xorm:"duration"`
	Calories     float64         `json:"calories" xorm:"calories"`
	Steps        int             `json:"steps" xorm:"steps"`
	Pace         int             `json:"pace" xorm:"pace"`
	FormPace     string          `json:"form_pace" xorm:"form_pace"`
	Points       []models.Points `json:"points" xorm:"points"`
	Frequency    float64         `json:"frequency" xorm:"frequency"`
	Frequencies  []int           `json:"frequencies" xorm:"frequencies"`
	CreateAt     string          `json:"create_at" xorm:"create_at created"`
	XFrequencies []int           `json:"x_frequencies" xorm:"-"`
	XNumber      []int           `json:"x_number" xorm:"-"`
	Status       int             `json:"status" xorm:"status"`
}

type resRecords struct {
	Records []responseRecord `json:"records"`
	Total   int64            `json:"total"`
	LastId  int              `json:"last_id"`
}

//学生运动计划数据
type responseProgress struct {
	Distance float32 `json:"distance" `
	Duration string  `json:"duration" `
	Times    int     `json:"times" `
	Pace     string  `json:"pace"`
}

type StudentAllInfos struct {
	Id             int       `json:"student_id" xorm:"id"`
	Code           string    `json:"student_code" xorm:"code"`
	Name           string    `json:"name" xorm:"name"`
	Gender         int       `json:"gender" xorm:"gender"`
	ClassName      string    `json:"class_name" xorm:"name"`
	DepartmentName string    `json:"department_name" xorm:"name"`
	Cellphone      string    `json:"cellphone" xorm:"cellphone"`
	CreateAt       time.Time `json:"create_at" xorm:"create_at"`
}

type responseStudentInfo struct {
	models.StudentAllInfos
	SchoolName string
}

//weekDistance,weekCount,weekDuration,weekSteps
type weekContent struct {
	Sequence int     `json:"sequence"`
	Distance int     `json:"distance"`
	Count    int     `json:"count"`
	Duration int     `json:"_"`
	Steps    int     `json:"_"`
	Pace     float32 `json:"pace"`
	Status   int     `json:"status"`
}

//学生基本信息 ，学号，姓名，班级，年级，院系

//学生所有跑步记录
func records(ctx iris.Context) {
	//创建查询Session
	query := lib.Engine.Table("plan_record")

	//字段查询
	if ctx.URLParamExists("student_id") {
		query.And("student_id=?", ctx.URLParam("student_id"))
	} else {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "无学生id"))
		return
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
			//ctx.JSON(lib.NewResponseFail(1, "order参数错误，必须是asc或desc"))
			ctx.JSON(lib.FailureResponse(lib.NilStruct(), "order参数错误，必须是asc或desc"))
			return
		}
	}

	//分页:用last_id来做限制条件，而且使用last_id就不用page，使用最新的一页
	size := ctx.URLParamIntDefault("size", 5)
	lastId := ctx.URLParamIntDefault("last_id", 0)
	//跳过的数量
	//passNum:=ctx.URLParamIntDefault("last_id", 0)
	query.Limit(size, 0)

	//获取计划名称
	query.Join("INNER", "plan", "plan.id=plan_record.plan_id")
	//当不是第一页时
	if lastId != 0 {
		query.And("plan_record.create_at<(select plan_record.create_at from plan_record where plan_record.id=?)", lastId)
	}

	//查询
	var planRecord []responseRecord
	toatal, err := query.FindAndCount(&planRecord)
	if err != nil {
		fmt.Printf("查询运动记录错误：%v", err)
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "查询运动记录错误"))
		return
	}

	retLastId := 0
	if len(planRecord) != 0 {
		retLastId = planRecord[len(planRecord)-1].Id
	}

	resRecords := resRecords{planRecord, toatal, retLastId}

	//获取total
	totalRecord := models.PlanRecord{}
	num, err := lib.Engine.Table("plan_record").Join("INNER", "plan", "plan.id=plan_record.plan_id").Where("student_id=?", ctx.URLParam("student_id")).Count(totalRecord)
	resRecords.Total = num

	ctx.JSON(lib.SuccessResponse(resRecords, "获取学生运动记录成功"))
}

//计划运动数据统计
func planProgress(ctx iris.Context) {

	studentId := ctx.Params().GetUint64Default("id", 0)

	//获取该学生的plan_id
	userinfo := models.Student{}
	fmt.Println("request:", ctx.Request())
	res, err := lib.Engine.Table("student").Where("id=?", studentId).Get(&userinfo)
	if err != nil {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "查询学生失败"))
		return
	}
	if res == false {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "未找到该学生"))
		return
	}
	fmt.Printf("学生：%v", userinfo)
	//获取学生计划进度
	planProgress := models.PlanProgress{}
	resPro, err := lib.Engine.Table("plan_progress").Where("student_id=?", studentId).And("plan_id=?", userinfo.PlanId).Get(&planProgress)
	if err != nil {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "计划进度查询失败"))
		return
	}
	if resPro == false {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "未找到该计划进度"))
		return
	}

	//平均配速
	records := make([]models.PlanRecord, 0)
	err = lib.Engine.Table("plan_record").Where("student_id=?", studentId).And("plan_id=?", userinfo.PlanId).And("status=?", 1).Find(records)
	if err != nil {
		fmt.Printf("查询学生运动记录id出错:%v", err)
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "运动记录查询失败"))
		return
	}
	paceSum := 0
	for _, value := range records {
		paceSum += value.Pace
	}
	println("")
	fmt.Printf("paceSum: %d", paceSum)
	//paceAver := float64(paceSum / len(records))
	var paceAver float64
	if len(records) > 0 {
		paceAver = float64(paceSum / len(records))
	} else {
		paceAver = float64(paceSum)
	}

	println("")
	fmt.Printf("len:%d, paceFloat64:%f", len(records), paceAver)
	paceAverageStr := lib.FomPace(paceAver)

	resProgress := responseProgress{
		Distance: float32(planProgress.Distance) / 1000,
		Duration: lib.ConvertTime(planProgress.Duration),
		Times:    planProgress.Times,
		Pace:     paceAverageStr,
	}

	ctx.JSON(lib.SuccessResponse(resProgress, "获取学生跑步进度成功"))
}

//查询学生
func search(ctx iris.Context) {

	code := ""
	name := ""
	if ctx.URLParamExists("code") {
		code = ctx.URLParam("code")
	} else {
		println("需要学号")
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "需要学号"))
		return
	}

	if ctx.URLParamExists("name") {
		name = ctx.URLParam("name")
	} else {
		println("需要姓名")
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "需要姓名"))
		return
	}
	//获取该学生
	userinfo := models.Student{}
	res, err := lib.Engine.Table("student").Where("code=?", code).And("name=?", name).Get(&userinfo)
	if err != nil {
		fmt.Printf("查询学生失败：%V", err)
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "查询学生失败"))
		return
	}
	if res == false {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "未找到该学生"))
		return
	}
	fmt.Printf("学生：%v", userinfo)

	ctx.JSON(lib.SuccessResponse(userinfo, "查询学生成功"))
}

//获取学生整合信息
func get(ctx iris.Context) {

	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	student := models.StudentAllInfos{}
	//根据id查询
	b, err := lib.Engine.Table("student").
		Join("INNER", "classes", "classes.id=student.class_id").
		Join("INNER", "department", "department.id=student.department_id").
		Where("student.id=?", id).
		Cols("student.id", "student.code", "student.name", "classes.name", "department.name", "student.gender", "student.create_at", "student.cellphone", "student.plan_id").
		Get(&student)

	if err != nil {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "查询学生信息错误"))
		return
	}
	if b == false {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "未找到该学生"))
		return
	}
	responsStudentInfo := responseStudentInfo{
		SchoolName:      configs.Conf.School.Name,
		StudentAllInfos: student,
	}
	ctx.JSON(lib.SuccessResponse(responsStudentInfo, "获取学生整合信息成功"))
}

//学生每周计划完成情况
//接口实现思路：
//1.获取当前计划有几周。
//2.获取每周的开始和结束时间，用来累加每周的公里数。
//3.根据周数返回不同数组长度的周进度。
//备注：自己测试，在一周，两周，三周和三周以上是正常的。需要注意：plan.DateBegin.Weekday()当时周日时，结果是0
func everyMonthProgress(ctx iris.Context) {

	//创建查询Session指针
	queryPlan := lib.Engine.Table("plan")
	queryRecord := lib.Engine.Table("plan_record")

	if ctx.URLParamExists("planId") {
		queryPlan.And("id=?", ctx.URLParam("planId"))
	} else {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "缺少planId"))
		return
	}
	if ctx.URLParamExists("studentId") {
		queryRecord.And("student_id=?", ctx.URLParam("studentId"))
	} else {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "缺少studentId"))
		return
	}
	println("获取的planId:", ctx.URLParam("planId"), "studentId:", ctx.URLParam("studentId"))
	//1.获取计划，包括计划的开始时间，结束时间，计划的周数，计划的每周开始时间和结束时间
	var plan models.Plan
	res, err := queryPlan.Get(&plan)
	if err != nil {
		println("")
		fmt.Printf("%v", err)
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "查询计划失败"))
		return
	}
	if res == false {

		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "未查到校园计划"))
		return
	}

	//获取学生信息
	student := models.Student{}
	resStudent, err := lib.Engine.Table("student").Where("id=?", ctx.URLParam("studentId")).Get(&student)
	if err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "查询学生失败"))
		return
	}
	if resStudent == false {
		ctx.JSON(lib.FailureResponse(lib.NilStruct(), "未查到学生信息"))
		return
	}
	timeLayout := "2006-01-02 15:04:05"
	planBegin := plan.DateBegin.Unix()
	//第一周的开始时间和结束时间（时间戳） 结束时间：%v,firstWeekEndDay.Format("2006-01-02 15:04:05"),
	weekNumber := int(plan.DateBegin.Weekday())
	if weekNumber == 0 { //周日，转成7
		weekNumber = 7
	}
	//firstWeekEndDay := planBegin + int64(((7-int(plan.DateBegin.Weekday()))+1)*3600*24)
	firstWeekEndDay := planBegin + int64(7-weekNumber+1)*3600*24
	if plan.DateBegin.Weekday() == 0 { //周日，是0
		firstWeekEndDay = planBegin + 3600*24
	}

	println("计划开始时周几：", plan.DateBegin.Weekday(), int(plan.DateBegin.Weekday()))
	fmt.Printf("第一周开始时间:%v，时间戳：%v,结束时间戳：%v", plan.DateBegin, planBegin, firstWeekEndDay)

	//最后一周的开始时间和结束时间（时间戳）
	planEnd := plan.DateEnd.Unix()
	lastWeekDays := int(plan.DateEnd.Weekday()) - 1
	//if int(plan.DateEnd.Weekday()) == 0 {
	//	lastWeekDays = 6
	//}
	if int(plan.DateEnd.Weekday()) == 0 {
		lastWeekDays = 7
	}
	lastWeekEndDay := planEnd - int64(lastWeekDays*3600*24)
	println("")
	fmt.Printf("最后一周开始时间戳：%v,最后一周结束时间戳:%v", lastWeekEndDay, planEnd)

	//每周的结束时间

	//总周数
	Weeks := (lastWeekEndDay-firstWeekEndDay)/604800 + 2
	fmt.Printf("取整 int:%d", (lastWeekEndDay-firstWeekEndDay)/604800)
	floWeeks := float32(lastWeekEndDay-firstWeekEndDay) / 604800
	fmt.Printf("floweeks: %f", floWeeks)
	if floWeeks <= 1 {
		Weeks = 1
	}
	fmt.Printf("中间周数:%v,取余：%v", Weeks, (lastWeekEndDay-firstWeekEndDay)%604800)

	weekUnixs := make([]int, 0)
	weekUnixs = append(weekUnixs, int(planBegin))

	//分成三种情况，一周，两周，三周和三周以上
	//一周：firstWeekEndDay==lastWeekEndDay

	//两周:secondWeekEndDay==lastWeekEndDay
	println("")
	planWeeks := 0
	println("第一周结束时间戳：", firstWeekEndDay, "form:", lib.UnixToFormTime(firstWeekEndDay), "最后一周结束时间戳：", planEnd, "form:", lib.UnixToFormTime(planEnd))
	if firstWeekEndDay >= planEnd { //一周
		weekUnixs = append(weekUnixs, int(planEnd))
		planWeeks = 1
	} else if (firstWeekEndDay + 604800) >= planEnd { //二周
		weekUnixs = append(weekUnixs, int(firstWeekEndDay))
		weekUnixs = append(weekUnixs, int(planEnd))
		planWeeks = 2
	} else { //三周和三周以上
		weekUnixs = append(weekUnixs, int(firstWeekEndDay))
		planWeeks = 3
		println("开始,长度：", len(weekUnixs))
		for index, value := range weekUnixs {
			println("index:", index, "value:", value, time.Unix(int64(value), 0).Format(timeLayout))
		}

		println("---------")

		println("中间开始第一周第一天：", time.Unix(int64(firstWeekEndDay), 0).Format(timeLayout), "最后一周最后一天", time.Unix(int64(lastWeekEndDay), 0).Format(timeLayout), "周数：", int((lastWeekEndDay-firstWeekEndDay)/604800))

		lastWeekStartDay := lastWeekEndDay + 1
		for i := 0; i < int((lastWeekStartDay-firstWeekEndDay)/604800); i++ {
			weekUnixs = append(weekUnixs, int(firstWeekEndDay)+604800*(i+1))
		}
		println("---------")
		println("中间,长度：", len(weekUnixs))
		for index, value := range weekUnixs {
			println("index:", index, "value:", value, value, time.Unix(int64(value), 0).Format(timeLayout))
		}

		println("---------")
		weekUnixs = append(weekUnixs, int(planEnd))
		println("最后,长度：", len(weekUnixs))

		for index, value := range weekUnixs {
			println("index:", index, "value:", value, value, time.Unix(int64(value), 0).Format(timeLayout))
		}
	}

	println("本计划的周数：", planWeeks)
	for index, value := range weekUnixs {
		println("index:", index, "value:", value, value, time.Unix(int64(value), 0).Format(timeLayout))
	}

	//2.获取该计划的运动记录，并按照开始时间到结束时间排序返回 .Cols("distance","times","duration","steps","create_at")
	//var records []models.PlanRecord
	records := make([]models.PlanRecord, 0)
	errRecords := queryRecord.Where("create_at>?", plan.DateBegin.Format("2006-01-02 15:04:05")).And("create_at<?", plan.DateEnd.Format("2006-01-02 15:04:05")).
		And("student_id=?", ctx.URLParam("studentId")).And("plan_id=?", ctx.URLParam("planId")).And("status=?", 1).
		Asc("create_at").Find(&records)
	if errRecords != nil {
		fmt.Printf("查询运动记录失败:%v", errRecords)
		lib.FailureResponse(lib.NilStruct(), "查询运动记录失败")
		return
	}

	//fmt.Printf("%v", records)

	//3.遍历运动记录，当运动记录在每周计划的时间内，则累加公里数，次数，时间和步数，不在每周计划时间内则每周计划的数组index累加1
	weekContents := make([]weekContent, 0)
	println("")

	println("")
	i3 := 0
	for i := 0; i < len(weekUnixs); i++ {
		newWeekContent := weekContent{}

		println("")
		//fmt.Printf("本周的开始时间戳%v,结束时间戳：%v",time.Unix(int64(weekUnixs[i]), 0).Format(timeLayout),time.Unix(int64(weekUnixs[i+1]), 0).Format(timeLayout))
		for i2 := i3; i2 < len(records); i2++ {
			println("i2:", i2, records[i2].CreateAt.Unix(), "weekUnix i:", i, weekUnixs[i])
			if int(records[i2].CreateAt.Unix()) >= weekUnixs[i] && int(records[i2].CreateAt.Unix()) <= weekUnixs[i+1] {
				//fmt.Printf("当前的运动记录在本周的开始时间和结束时间内，id:%v,创建时间为：%v", records[i2].Id, records[i2].CreateAt)
				newWeekContent.Distance = newWeekContent.Distance + records[i2].Distance
				newWeekContent.Count = newWeekContent.Count + 1
				newWeekContent.Steps = newWeekContent.Steps + records[i2].Steps
				newWeekContent.Duration = newWeekContent.Duration + records[i2].Duration
				i3++
				println("当前record index:", i3)
			} else {
				break
			}
		}

		if newWeekContent.Steps != 0 && newWeekContent.Duration != 0 {
			newWeekContent.Pace = float32(newWeekContent.Steps) / (float32(newWeekContent.Duration) / 60) //步频,修改精度问题
		}
		weekContents = append(weekContents, newWeekContent)

	}

	var weekDistance int //每周运动距离 不支持按照次数情况下的公里数
	weekDistance = plan.MinWeekDistance
	println("当前学生的weekDistance:", weekDistance)
	//4.计算每周的步频 步频的计算公式，总步数/分钟数
	for index, value := range weekContents {
		println("本周distance:", weekContents[index].Distance, "要求每周距离：", weekDistance)
		if weekContents[index].Distance >= weekDistance {
			weekContents[index].Status = 1
		}
		weekContents[index].Sequence = index + 1
		fmt.Printf("第%v,的数据：%v", index, value)
	}

	//根据不同周数返回不同长度数组
	planWeekContents := make([]weekContent, 0)
	println("weekContents length:", len(weekContents))
	println("planWeeks:", planWeeks)
	switch planWeeks {
	case 1:
		planWeekContents = append(planWeekContents, weekContents[0])
		ctx.JSON(lib.SuccessResponse(planWeekContents, "返回学生每周运动数据成功"))
		return
	case 2:
		planWeekContents = append(planWeekContents, weekContents[0])
		planWeekContents = append(planWeekContents, weekContents[1])
		ctx.JSON(lib.SuccessResponse(planWeekContents, "返回学生每周运动数据成功"))
		return
	case 3:
		finWeekContents := weekContents[0 : len(weekContents)-1]
		ctx.JSON(lib.SuccessResponse(finWeekContents, "返回学生每周运动数据成功"))
		return
	}

	//ctx.JSON(lib.NewResponseOK(weekContents))
}
