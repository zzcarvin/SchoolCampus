package student

import (
	"Campus/configs"
	"Campus/internal/lib"
	"Campus/internal/models"
	"github.com/go-xorm/builder"
	"github.com/kataras/iris"
	"strings"
	//"encoding/json"
	//"fmt"
	"fmt"
	"time"
)

type responseStudentInfo struct {
	models.StudentAllInfos
	SchoolName string
}

// swagger:parameters  StudentCreateRequest
type StudentCreateRequest struct {
	// in: body
	Body models.Student
}

// 响应结构体
//
// swagger:response    StudentCreateResponse
type StudentCreateResponse struct {
	// in: body
	Body studentresponseMessage
}
type studentresponseMessage struct {
	models.ResponseType
	Data models.Student
}

type studentProgress struct {
	Code                 int    `json:"code" xorm:"code"`
	Name                 string `json:"name" xorm:"name"`
	RecordTimes          int    `json:"week_times" xorm:"weekTimes"`
	RecordDistance       int    `json:"week_distance" xorm:"weekDistance"`
	PlanDistance         int    `json:"plan_distance" xorm:"-"`
	BoyTotalDistance     int    `json:"-" xorm:"boy_total_distance"`
	GirlTotalDistance    int    `json:"-" xorm:"girl_total_distance"`
	PlanProgressDistance int    `json:"plan_progress_distance" xorm:"distance"`
	Times                int    `json:"times" xorm:"times"`
	//Steps                int    `json:"steps" xorm:"steps"`
	//ValidateTimes int `json:"validate_times"`
	Progress float32 `json:"progress" xorm:"-"`
	Gender   int     `json:"-" xorm:"gender"`
}

type recordFeedbacks struct {
	Id             int       `json:"id" xorm:"id"`
	FeedbackId     int       `json:"feedback_id" xorm:"id"`
	Content        string    `json:"content" xorm:"content"`
	Code           int       `json:"code" xorm:"code"`
	Name           string    `json:"name" xorm:"name"`
	RecordDistance int       `json:"distance" xorm:"distance"`
	StartTime      time.Time `json:"start_time" xorm:"start_time"`
	EndTime        time.Time `json:"end_time" xorm:"end_time"` //用开始时间和结束时间来获取时长
	Duration       int       `json:"duration" xorm:"duration"`
	Pace           int       `json:"pace" xorm:"pace"`
	CreateAt       time.Time `json:"plan_record_create_at" xorm:"create_at"`
	Status         int       `json:"status" xorm:"status"`
	FeedbackStatus int       `json:"feedback_status" xorm:"status"`
	FeedBack       string    `json:"feed_back" xorm:"content"`
	//FeedbackCreateAt time.Time `json:"feedback_create_at" xorm:"create_at"`
}

type responseFeedbacks struct {
	Feedbacks     []recordFeedbacks
	MonFinish     int `json:"mon_finish" xorm:"-"`
	MonUnfinished int `json:"mon_unfinished" xorm:"-"`
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

type records struct {
	Distance int       `json:"distance" xorm:"distance"`
	Times    int       `json:"times" xorm:"times"`
	Duration int       `json:"duration" xorm:"duration"`
	Steps    int       `json:"steps" xorm:"steps"`
	CreateAt time.Time `json:"create_at" xorm:"create_at"`
}

func create(ctx iris.Context) {
	// swagger:route POST /api/student/ student StudentCreateRequest
	//
	// 创建学生表
	//
	//     Consumes:
	//     - application/json
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//	   - 200: StudentCreateResponse

	student := models.Student{}

	//解析department
	err := ctx.ReadJSON(&student)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return

	}
	//插入数据
	res, err := lib.Engine.Table("student").Insert(&student)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	lib.NewResponseOK(res)
	ctx.JSON(lib.NewResponseOK(student))

}

// swagger:route DELETE /api/student/:id  student StudentDelete
//
// 删除学生表
//
//     Produces:
//     - application/json
//
//     Responses:
//	   - 200: Response
func remove(ctx iris.Context) {
	id := ctx.Params().GetUint64Default("id", 0)
	student := models.Student{}
	affected, err := lib.Engine.Table("student").ID(id).Delete(&student)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	lib.NewResponseOK(affected)

}

// swagger:parameters  StudentUpdateRequest
type StudentUpdateRequest struct {
	// in: body
	Body models.Student
}

// 响应结构体
//
// swagger:response    StudentUpdateResponse
type StudentUpdateResponse struct {
	// in: body
	Body studentresponseMessage
}

func update(ctx iris.Context) {
	// swagger:route PUT /api/student/:id student StudentUpdateRequest
	//
	// 修改学生表
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//	   - 200:StudentUpdateResponse

	id := ctx.Params().GetUint64Default("id", 0)

	student := models.Student{}

	//解析department
	err := ctx.ReadJSON(&student)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//TODO 验证数据有效性

	//插入数据
	res, err2 := lib.Engine.Table("student").ID(id).Update(student)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	lib.NewResponseOK(res)
	ctx.JSON(lib.NewResponseOK(student))
}

// swagger:route GET /api/student/:id  student StudentGet
//
// 查询学生表
//     Consumes:
//     - application/json
//     Produces:
//     - application/json
//
//     Responses:
//       200: Response
func get(ctx iris.Context) {

	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	student := models.StudentAllInfos{}
	//根据id查询
	b, err := lib.Engine.Table("student").
		Join("INNER", "classes", "classes.id=student.class_id").
		Join("INNER", "department", "department.id=student.department_id").
		Where("student.id=?", id).
		Cols("student.id", "student.code", "student.name", "classes.name", "department.name", "student.gender", "student.create_at", "student.cellphone").
		Get(&student)
	// SELECT `student`.`id`, `student`.`code`, `student`.`name`, `classes`.`name`, `department`.`name`, `student`.`gender`, `student`.`create_at`, `student`.`cellphone`
	// FROM `student` INNER JOIN classes ON classes.id=student.class_id INNER JOIN department ON department.id=student.department_id WHERE (student.id=?) LIMIT 1 []interface {}{0x1}

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该用户"))
		return
	}
	responsStudentInfo := responseStudentInfo{
		SchoolName:      configs.Conf.School.Name,
		StudentAllInfos: student,
	}
	ctx.JSON(lib.NewResponseOK(responsStudentInfo))
}

// swagger:route GET /api/students student StudentSearch
//
// 查询学生表（按字段查询 +s  example http://localhost:8081/api/students?page=0&size=5&name=李一帆&sort=id&order=desc 其他search接口相同）
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: Response

func search(ctx iris.Context) {

	//创建查询Session指针
	query := lib.Engine.Table("student")

	//字段查询
	if ctx.URLParamExists("name") {
		query.And(builder.Like{"name", ctx.URLParam("name")})
	}
	if ctx.URLParamExists("classid") {
		query.And(builder.Like{"class_id", ctx.URLParam("classid")})
	}
	if ctx.URLParamExists("derparmentid") {
		query.And(builder.Like{"derparment_id", ctx.URLParam("derparmentid")})
	}
	if ctx.URLParamExists("code") {
		query.And(builder.Like{"code", ctx.URLParam("code")})
	}
	if ctx.URLParamExists(("userid")) {
		query.And(builder.Like{"uer_id", ctx.URLParam("userid")})
	}
	if ctx.URLParamExists("cellphone") {
		query.And(builder.Like{"cellphone", ctx.URLParam("cellphone")})
	}
	if ctx.URLParamExists("classes_name") {
		classesname := ctx.URLParam("classes_name")
		query.
			Join("INNER", "plan_record", "student.id=plan_record.student_id").
			Join("INNER", "classes", "student.class_id=classes.id").
			Where("classes.name =?", classesname).
			Cols("plan_record.*")
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
	var student []models.Student
	err := query.Find(&student)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(student))
}
func getcount(ctx iris.Context) {
	var classid int64
	classid = ctx.URLParamInt64Default("classid", 0)
	var classcount int64
	count := new(models.Student)
	classcount, err := lib.Engine.Where("class_id=?", classid).Count(count)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	ctx.JSON(lib.NewResponseOK(classcount))
}

//学生计划进度
func studentPlanProgress(ctx iris.Context) {
	//创建查询Session指针
	query := lib.Engine.Table("plan_progress")

	if ctx.URLParamExists("studentCode") {
		query.And(builder.Like{"code", ctx.URLParam("studentCode")})
	}
	if ctx.URLParamExists("studentName") {
		query.And(builder.Like{"student.name", ctx.URLParam("studentName")})
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

	var planProgress []studentProgress

	//根据id查询
	err := query.Select("(SELECT COUNT(`plan_record`.`id`) FROM `plan_record` WHERE `plan_record`.`student_id`=`plan_progress`.`student_id`) AS weekTimes ,"+
		"(SELECT SUM(`plan_record`.`distance`) FROM `plan_record` WHERE `plan_record`.`student_id`=`plan_progress`.`student_id`) AS weekDistance "+
		",plan.boy_total_distance"+",plan.girl_total_distance"+",plan_progress.distance"+",plan_progress.times"+",student.gender"+",student.code"+",student.name").
		Join("INNER", "plan", "plan.id=plan_progress.plan_id").
		Join("INNER", "student", "student.id=plan_progress.student_id").
		Find(&planProgress)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//将整理的数据放入数组中
	for index, _ := range planProgress {
		if planProgress[index].Gender == 1 {
			planProgress[index].PlanDistance = planProgress[index].BoyTotalDistance
		} else {
			planProgress[index].PlanDistance = planProgress[index].GirlTotalDistance
		}
		planProgress[index].Progress = float32(planProgress[index].PlanProgressDistance) / float32(planProgress[index].PlanDistance)
		fmt.Printf("planProgressDistance%v/planDisatance%v=progress%v", float32(planProgress[index].PlanProgressDistance), float32(planProgress[index].PlanDistance), planProgress[index].Progress)
	}

	//fmt.Printf("rest:%v",planProgress)
	ctx.JSON(lib.NewResponseOK(planProgress))
}

//学生异常处理
func studentAbnormal(ctx iris.Context) {
	//1.学生表，运动记录表，记录反馈表---学生异常申请
	//需要在反馈表添加运动记录id字段，反馈表增加字段status表示反馈是否处理比如0，未处理，1，已处理
	//创建查询Session指针
	query := lib.Engine.Table("plan_record")

	var feedbacks []recordFeedbacks

	query.Join("INNER", "student", "student.id=plan_record.student_id").
		Join("INNER", "feedback", "feedback.record_id=plan_record.id").
		//Where("plan_record.status=?",0).
		Cols("student.code",
			"student.name",
			"plan_record.id",
			"feedback.id",
			"plan_record.distance",
			"plan_record.start_time",
			"plan_record.end_time",
			"plan_record.pace",
			"feedback.create_at",
			"plan_record.duration",
			"feedback.content",
			"plan_record.status",
			"feedback.status")

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
	err := query.Find(&feedbacks)

	if err != nil {
		fmt.Printf("%v", err.Error())
		ctx.JSON(lib.NewResponseFail(0, err.Error()))
		return
	}
	//err1 :=query.Join("INNER", "feedback", "plan_record.id=feedback.record_id").Cols("feedback.status", "feedback.content", "feedback.create_at").Find(&feedbacks)
	//if err1 != nil {
	//	fmt.Printf("%v", err1.Error())
	//	ctx.JSON(lib.NewResponseFail(0, err1.Error()))
	//	return
	//}
	//获取本月开始时间
	timeLayout := "2006-01-02 15:04:05"
	now := time.Now().Local().Format(timeLayout)
	nowMonthStart, _ := time.ParseInLocation(timeLayout, string(now[0:7])+"-01 00:00:00", time.Local)
	nowTime, _ := time.ParseInLocation(timeLayout, now, time.Local)

	nowUnix := nowTime.Unix()
	MonStartUnix := nowMonthStart.Unix()
	tm := time.Unix(MonStartUnix, 0)
	fmt.Println(tm.Format("2006-01-02 15:04:05"))
	tm1 := time.Unix(nowUnix, 0)
	fmt.Println(tm1.Format("2006-01-02 15:04:05"))

	fmt.Printf("当前时间戳：%v,本月开始时间戳：%v", nowUnix, MonStartUnix)

	finishFeedbacks := 0
	unfinishFeedbacks := 0
	for index, _ := range feedbacks {
		if feedbacks[index].CreateAt.Unix() >= MonStartUnix && feedbacks[index].CreateAt.Unix() <= nowUnix {
			//println(fmt.Printf("当前时间戳：%v,本月开始时间戳：%v", nowUnix, MonStartUnix))
			fmt.Printf("当前反馈在本月：%v", feedbacks[index])
			if feedbacks[index].FeedbackStatus == 0 {
				unfinishFeedbacks++
				//println("\n未完成计数\n",unfinishFeedbacks)
			} else {
				finishFeedbacks++
				//println("\n完成计数\n",unfinishFeedbacks)
			}
		}
		//	if feedbacks[index].Status == 0 {
		//		unfinishFeedbacks++
		//	} else {
		//		finishFeedbacks++
		//	}

	}

	responseFeedbacks := responseFeedbacks{
		Feedbacks:     feedbacks,
		MonFinish:     finishFeedbacks,
		MonUnfinished: unfinishFeedbacks,
	}

	ctx.JSON(lib.NewResponseOK(responseFeedbacks))

}

//学生每周计划完成情况
func everyMonthProgress(ctx iris.Context) {

	//创建查询Session指针
	queryPlan := lib.Engine.Table("plan")
	queryRecord := lib.Engine.Table("plan_record")

	if ctx.URLParamExists("planId") {
		queryPlan.And("id=?", ctx.URLParam("planId"))
	} else {
		ctx.JSON(lib.NewResponseFail(0, "缺少planId"))
		return
	}
	if ctx.URLParamExists("studentId") {
		queryRecord.And("student_id=?", ctx.URLParam("studentId"))
	} else {
		ctx.JSON(lib.NewResponseFail(0, "缺少studentId"))
		return
	}
	//1.获取计划，包括计划的开始时间，结束时间，计划的周数，计划的每周开始时间和结束时间
	var plan models.Plan
	res, err := queryPlan.Get(&plan)
	if err != nil {
		fmt.Printf("%v", err)
		lib.NewResponseFail(0, err.Error())
		return
	}
	if res == false {
		lib.NewResponseFail(0, "未查到校园计划")
		return
	}

	//获取学生信息
	student := models.Student{}
	resStudent, err := lib.Engine.Table("student").Where("id=?", ctx.URLParam("studentId")).Get(&student)
	if err != nil {
		fmt.Printf("%v", err)
		lib.NewResponseFail(0, err.Error())
		return
	}
	if resStudent == false {
		lib.NewResponseFail(0, "未查到学生信息")
		return
	}
	timeLayout := "2006-01-02 15:04:05"
	planBegin := plan.DateBegin.Unix()
	//第一周的开始时间和结束时间（时间戳） 结束时间：%v,firstWeekEndDay.Format("2006-01-02 15:04:05"),
	firstWeekEndDay := planBegin + int64(((7-int(plan.DateBegin.Weekday()))+1)*3600*24)
	fmt.Printf("第一周开始时间:%v，时间戳：%v,结束时间戳：%v", plan.DateBegin, planBegin, firstWeekEndDay)

	//最后一周的开始时间和结束时间（时间戳）
	planEnd := plan.DateEnd.Unix()
	lastWeekDays := int(plan.DateEnd.Weekday()) - 1
	if int(plan.DateEnd.Weekday()) == 0 {
		lastWeekDays = 6
	}
	lastWeekEndDay := planEnd - int64(lastWeekDays*3600*24)
	println("")
	fmt.Printf("最后一周开始时间戳：%v,最后一周结束时间戳:%v", lastWeekEndDay, planEnd)

	//每周的结束时间

	//总周数
	Weeks := (lastWeekEndDay-firstWeekEndDay)/604800 + 2
	fmt.Printf("中间周数:%v,取余：%v", Weeks, (lastWeekEndDay-firstWeekEndDay)%604800)

	weekUnixs := make([]int, 0)
	weekUnixs = append(weekUnixs, int(planBegin))
	weekUnixs = append(weekUnixs, int(firstWeekEndDay))
	println("开始,长度：", len(weekUnixs))
	for index, value := range weekUnixs {
		println("index:", index, "value:", value, time.Unix(int64(value), 0).Format(timeLayout))
	}
	println("---------")
	for i := 0; i < int((lastWeekEndDay-firstWeekEndDay)/604800); i++ {
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

	//2.获取该计划的运动记录，并按照开始时间到结束时间排序返回 .Cols("distance","times","duration","steps","create_at")
	//var records []models.PlanRecord
	records := make([]models.PlanRecord, 0)
	errRecords := queryRecord.Where("create_at>?", plan.DateBegin.Format("2006-01-02 15:04:05")).And("create_at<?", plan.DateEnd.Format("2006-01-02 15:04:05")).
		And("student_id=?", ctx.URLParam("studentId")).And("plan_id=?", ctx.URLParam("planId")).
		Asc("create_at").Find(&records)
	if errRecords != nil {
		lib.NewResponseFail(0, errRecords.Error())
		return
	}

	fmt.Printf("%v", records)

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
			if int(records[i2].CreateAt.Unix()) >= weekUnixs[i] && int(records[i2].CreateAt.Unix()) <= weekUnixs[i+1] {
				fmt.Printf("当前的运动记录在本周的开始时间和结束时间内，id:%v,创建时间为：%v", records[i2].Id, records[i2].CreateAt)
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
			newWeekContent.Pace = float32(newWeekContent.Steps) / float32(newWeekContent.Duration/60) //步频
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

	ctx.JSON(lib.NewResponseOK(weekContents))
}
