package student

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"bytes"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/kataras/iris"
	"path"
	"path/filepath"

	"io"
	"log"
	"os"
	"strconv"

	"Campus/configs"
	"github.com/go-xorm/builder"
	"strings"
	//"encoding/json"
	//"fmt"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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
	InvalidCode   []int  `json:"invalid_code" xorm:"invalid_code"`
	Reply_message string `json:"reply_message" xorm:"reply_message"`
	Check_status  int    `json:"check_status" xorm:"check_status"`
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
	planrecord := models.PlanRecord{}
	//
	//affected, err := lib.Engine.Table("student").ID(id).Delete(&student)
	//if err != nil {
	//	ctx.JSON(lib.NewResponseFail(1, err.Error()))
	//	return
	//}

	session := lib.Engine.NewSession()
	defer session.Close()

	err1 := session.Begin()
	if err1 != nil {
		ctx.JSON(lib.NewResponseFail(1, "事务开启失败"))
		println("事务开启失败")
		return
	}
	affected, err2 := session.Table("student").ID(id).Delete(&student)
	if err2 != nil {
		session.Rollback()
		ctx.JSON(lib.NewResponseFail(1, "删除学生失败"))
		return
	}
	recordaffected, err3 := session.Table("plan_record").Where("student_id=?", id).Delete(&planrecord)
	if err3 != nil {
		session.Rollback()
		ctx.JSON(lib.NewResponseFail(1, "删除计划记录失败"))
		return
	}

	err4 := session.Commit()
	if err4 != nil {
		panic(err4.Error())
	}
	ctx.JSON(lib.NewResponseOK(affected))
	ctx.JSON(lib.NewResponseOK(recordaffected))

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
	//querystudent := lib.Engine.Table("student")
	//
	//students := []models.Student{}
	//student := models.Student{
	//	PlanId: 1000,
	//}
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
	//****************************************************************
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
	//**********************************************************************
	//year := ctx.URLParamIntDefault("year", 0)
	//
	//departmentid := ctx.URLParamIntDefault("department_id", 0)
	//classid := ctx.URLParamIntDefault("class_id", 0)
	//gender := ctx.URLParamIntDefault("gender", 0)
	////updatesql :="UPDATE `student` SET `plan_id` = 100 WHERE `id` in"
	//
	//if year != 0 {
	//	querystudent.Where("year=?", year)
	//} else {
	//	ctx.JSON(lib.NewResponseFail(1, "年份不存在"))
	//
	//}
	//
	//if departmentid != 0 {
	//	querystudent.And("department_id=?", departmentid)
	//}
	//
	//if classid != 0 {
	//	querystudent.And("class_id=?", classid)
	//}
	//
	//if gender != 0 {
	//	querystudent.And("gender=?", year)
	//}
	//errcount := querystudent.Cols("id").Find(&students)
	//if errcount != nil {
	//	ctx.JSON(lib.NewResponseFail(1, errcount.Error()))
	//	return
	//}
	//var id []int
	//for _, student := range students {
	//	id = append(id, student.Id)
	//}
	//fmt.Println("\n\n\n\n", id)
	//res, errupdate := querystudent.ID(id).Cols("plan_id").Update(student)
	//if errupdate != nil {
	//	ctx.JSON(lib.NewResponseFail(1, errupdate.Error()))
	//	return
	//}
	//
	//lib.NewResponseOK(res)
	//ctx.JSON(lib.NewResponseOK(res))

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
	if ctx.URLParamExists("departmentid") {
		query.And(builder.Like{"department_id", ctx.URLParam("departmentid")})
	}
	if ctx.URLParamExists("code") {
		query.And(builder.Like{"code", ctx.URLParam("code")})
	}
	if ctx.URLParamExists(("userid")) {
		query.And(builder.Like{"user_id", ctx.URLParam("userid")})
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

//分两次查询学生计划进度 这是第一次查出数据后 放入的结构体
type progress1 struct {
	Code                 int    `json:"code" xorm:"code"`
	Name                 string `json:"name" xorm:"name"`
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
type progress2 struct {
	RecordTimes    int `json:"week_times" xorm:"weekTimes"`
	RecordDistance int `json:"week_distance" xorm:"weekDistance"`
}

//学生计划进度
func studentPlanProgress(ctx iris.Context) {
	//创建查询Session指针
	query := lib.Engine.Table("plan_progress")
	//query1 :=lib.Engine.Table("plan_progress")
	//用不到，是全部查询，不需要指定学号，姓名，前端没改，这里也暂时不改
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

	now := time.Now()

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	weekend := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset+7)
	fmt.Println(weekStart)
	fmt.Println(weekend)

	var planProgress []studentProgress
	//progress2 :=progress2{}
	//progress1 :=progress1{}
	//query.Cols("planRecord.Distance").Where("planRecord.CreateAt < ?",weekStart).Where("planRecord.CreateAt>?",weekend)
	//
	//query.Join("INNER","student","student.id=plan_progress.student_id").
	//	  Join("INNER", "plan", "plan.id=plan_progress.plan_id").
	//	  Cols("plan.boy_total_distance","plan.girl_total_distance", "plan_progress.distance ",
	//	  	"plan_progress.times","student.gender","student.code","student.name").Select("distance as week_times",).
	//	  Where("planRecord.CreateAt < ?",weekStart).Where("planRecord.CreateAt>?",weekend)

	//b,err :=query.Join("INNER", "plan", "plan.id=plan_progress.plan_id").
	//	Join("INNER", "student", "student.id=plan_progress.student_id").
	//	Cols("plan.boy_total_distance", "plan.girl_total_distance", "plan_progress.distance", "plan_progress.times", "student.gender", "student.code", "student.name").Get(progress1)
	//if err != nil {
	//	ctx.JSON(lib.NewResponseFail(1, err.Error()))
	//	return
	//}
	//if b == false {
	//	ctx.JSON(lib.NewResponseFail(1, "未找到该progress1包含的计划进度字段"))
	//	return
	//}
	//query1.

	//根据id查询
	err := query.Select("(SELECT sum(`plan_record`.`status`) FROM `plan_record` WHERE `plan_record`.`student_id`=`plan_progress`.`student_id`   AND  YEARWEEK( DATE_FORMAT(  `plan_record`.`create_at`, '%Y-%m-%d' ),1 ) = YEARWEEK( NOW(),1 )) AS weekTimes ,(SELECT SUM(`plan_record`.`distance`) FROM `plan_record` WHERE `plan_record`.`student_id`=`plan_progress`.`student_id` AND (status = 1) AND YEARWEEK( DATE_FORMAT(  `plan_record`.`create_at`, '%Y-%m-%d' ),1 ) = YEARWEEK( NOW(),1 ) ) AS weekDistance ,plan.boy_total_distance,plan.girl_total_distance,plan_progress.distance,plan_progress.times,student.gender,student.code,student.name ").And("plan.stop=?", 1).
		Join("INNER", "plan", "plan.id=plan_progress.plan_id").
		Join("INNER", "student", "student.id=plan_progress.student_id").Find(&planProgress)

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

	fmt.Printf("rest:%v", planProgress)
	b, err := ctx.JSON(lib.NewResponseOK(&planProgress))
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == 0 {
		ctx.JSON(lib.NewResponseFail(1, "未找到该progress1包含的计划进度字段"))
		return
	}
	//fmt.Println("****************************",planProgress)
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
			"feedback.status",
			"plan_record.invalid_code",
			"feedback.Reply_message",
			"feedback.Check_status",
		)
	//根据处理状态返回
	if ctx.URLParamExists("feedback_status") {

		query.Where("feedback.status = ?", ctx.URLParam("feedback_status"))

	}
	//排序
	if ctx.URLParamExists("feedback_status") {
		query.And(builder.Like{"feedback.status", ctx.URLParam("status")})
	}

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
		ctx.JSON(lib.NewResponseFail(0, "缺少planId"))
		return
	}
	if ctx.URLParamExists("studentId") {
		queryRecord.And("student_id=?", ctx.URLParam("studentId"))
	} else {
		ctx.JSON(lib.NewResponseFail(0, "缺少studentId"))
		return
	}
	println("获取的planId:", ctx.URLParam("planId"), "studentId:", ctx.URLParam("studentId"))
	//1.获取计划，包括计划的开始时间，结束时间，计划的周数，计划的每周开始时间和结束时间
	var plan models.Plan
	res, err := queryPlan.Get(&plan)
	if err != nil {
		println("")
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(0, err.Error()))
		return
	}
	if res == false {

		ctx.JSON(lib.NewResponseFail(0, "未查到校园计划"))
		return
	}

	//获取学生信息
	student := models.Student{}
	resStudent, err := lib.Engine.Table("student").Where("id=?", ctx.URLParam("studentId")).Get(&student)
	if err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(0, err.Error()))
		return
	}
	if resStudent == false {
		ctx.JSON(lib.NewResponseFail(0, "未查到学生信息"))
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
	println("第一周结束时间戳：", firstWeekEndDay, "form:", UnixToFormTime(firstWeekEndDay), "最后一周结束时间戳：", planEnd, "form:", UnixToFormTime(planEnd))
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
		lib.NewResponseFail(0, errRecords.Error())
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
		planWeekContents := weekContents[0]
		ctx.JSON(lib.NewResponseOK(planWeekContents))
		return
	case 2:
		planWeekContents = append(planWeekContents, weekContents[0])
		planWeekContents = append(planWeekContents, weekContents[1])
		ctx.JSON(lib.NewResponseOK(planWeekContents))
		return
	case 3:
		finWeekContents := weekContents[0 : len(weekContents)-1]
		ctx.JSON(lib.NewResponseOK(finWeekContents))
		return
	}

	//ctx.JSON(lib.NewResponseOK(weekContents))
}

func zerotime() (zerotime int64) {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	fmt.Println(t.Format(time.UnixDate))
	//Unix返回早八点的时间戳，减去8个小时
	timestamp := t.UTC().Unix() - 8*3600
	fmt.Println("timestamp:", timestamp)
	return timestamp

}

type Sorts struct {
	Num int //自己加的管道tag标记
	Sum int //每次根据条件查出的跑步次数
}

//跑步页面学生个人信息柱状图代码
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
		println("1111111111", month)

	}
	if ctx.URLParamExists("year") {
		year = int(ctx.URLParamInt64Default("year", 0))

	}
	if ctx.URLParamExists("student_id") {
		studentid = int(ctx.URLParamInt64Default("student_id", 0))

	}

	//输入到管道缓存
	for i := 0; i < 48; i++ {

		zerotime1 := zerotime
		zerotime = zerotime + 30*60
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
		if count == 48 {
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
func execll(ctx iris.Context) {
	student := models.Student{}
	b, _ := lib.Engine.Table("department").Where("name=?", "物理系").Cols("department.id").Get(student)
	if b == false {
		department := "未找到" + "物理系" + "院系的id"
		ctx.JSON(lib.NewResponseFail(1, department))
		return
	}

}

//func execl(ctx iris.Context) {
//
//	println("yunxing")
//	// Get the file from the request.
//	file, info, err1 := ctx.FormFile("file")
//	if err1 != nil {
//		ctx.StatusCode(iris.StatusInternalServerError)
//		ctx.JSON(lib.NewResponseFail(0, "文件上传失败"))
//		return
//	}
//	ctx.JSON(lib.NewResponseOK("文件上传成功"))
//	println("yunxing2")
//	defer file.Close()
//	fname := info.Filename
//
//	//创建一个具有相同名称的文件
//	//假设你有一个名为'uploads'的文件夹
//	out, err2 := os.OpenFile("./execl/"+fname,
//		os.O_WRONLY|os.O_CREATE, 0666)
//	println("yunxing3")
//	if err2 != nil {
//		ctx.StatusCode(iris.StatusInternalServerError)
//
//		return
//	}
//	defer out.Close()
//	io.Copy(out, file)
//
//	dir := string("./execl/" + fname)
//	println("\n打印打印\n", dir)
//
//	f, err3 := excelize.OpenFile(dir)
//	if err3 != nil {
//		fmt.Println(err3)
//		return
//	}
//	//// 获取工作表中指定单元格的值
//	//cell := f.GetCellValue("Sheet1", "B2")
//
//	//fmt.Println(cell)
//	//获取 Sheet1 上所有单元格
//	rows := f.GetRows("Sheet1")
//	//rows := f.GetRows("Sheet1","A")
//	err5 := XormIni1()
//	if err5 != nil {
//		fmt.Println(err5)
//		return
//	}
//	ctx.JSON("新的数据库引擎初始化成功")
//
//	//defer lib.Engine.Close()
//
//	for i, row := range rows {
//		fmt.Println("\n \n", row)
//		if i != 0 {
//
//			go func(row []string) {
//				department := models.Department{}
//				classes := models.Classes{}
//				student := models.Student{}
//				fmt.Print(row, "\t")
//				student.Code = row[0]
//				student.Name = row[1]
//				//判断表格里的性别，并转化为int 男1女2
//				if row[2] == "男" {
//					var gender int = 1
//					student.Gender = gender
//				} else {
//					var gender int = 2
//					student.Gender = gender
//				}
//				println(row[3])
//
//				b, err := Engine1.Table("department").Where("name=?", row[3]).Get(&department)
//				if b == false {
//					department := "未找到" + row[3] + "院系的id"
//					ctx.JSON(lib.NewResponseFail(1, department))
//					return
//				}
//				if err != nil {
//					ctx.JSON(lib.NewResponseFail(1, err.Error()))
//					return
//				}
//
//				student.DepartmentId = department.Id
//				b2, _ := Engine1.Table("classes").Where("name= ?", row[4]).Get(&classes)
//				if b2 == false {
//					classes := "未找到" + row[3] + "班级的id"
//					ctx.JSON(lib.NewResponseFail(1, classes))
//					return
//				}
//				student.ClassId = classes.Id
//				res, _ := Engine1.Table("student").Insert(&student)
//				if res == 0 {
//					classes := "学号为" + row[0] + "的学生信息插入失败," + "或者该学生已存在"
//					ctx.JSON(lib.NewResponseFail(1, classes))
//					return
//
//				}
//			}(row)
//		}
//
//	}
//	err4 := XormClose1()
//	if err4 != nil {
//		//TODO log error
//		fmt.Printf("%v", err1.Error())
//		fmt.Printf("\n\n\n\n引擎关闭成功\n\n\n\n\n\n\n\n")
//	}
//
//}
//
////******************加一个新的数据库引擎*******************************************************************
//var Engine1 *xorm.Engine
//
//func XormIni1() error {
//	if Engine1 != nil {
//		return fmt.Errorf("Xorm已经初始化")
//	}
//
//	//获取配置
//	cfg := configs.Conf.Database
//	println("mysql连接信息:", cfg.Conn)
//	var err error
//	Engine1, err = xorm.NewEngine(cfg.Driver, cfg.Conn)
//
//	if err != nil {
//		fmt.Printf("xorm初始化失败：%v", err.Error())
//		return err
//	}
//
//	//TODO 配置化
//	Engine1.SetMaxIdleConns(1000)
//	Engine1.SetMaxOpenConns(1000)
//
//	//打印调试信息
//	Engine1.ShowSQL(true)
//	Engine1.Logger().SetLevel(core.LOG_DEBUG)
//
//	//err = Engine.Sync2(new(models.PlanProgress), new(models.Plan), new(models.PlanRecord),new(models.PlanRoute))
//	//if err != nil {
//	//	print(err)
//	//	return err
//	//}
//	//这个只执行一次，平时不要执行，这个要是每次都执行，相当于每次都重新创建数据库表，只要有不同的地方，就会改数据库，不推荐这样同步数据库，如果要改数据库结构，直接去数据库改
//	return err
//}
//
//func XormClose1() error {
//	err := Engine1.Close()
//	if err != nil {
//		//TODO log error
//		fmt.Printf("%v", err.Error())
//	}
//	return err
//}

func execl(ctx iris.Context) {

	type StudentInfo struct {
		ID           string
		realName     string
		gender       int
		departmentId int
		classId      int
	}

	departments := make(map[string]int)
	classes := make(map[string]int)
	sqlStr := "insert into student(department_id,class_id,name,gender,code) values"

	println("yunxing")
	// Get the file from the request.
	file, info, err := ctx.FormFile("file")
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(lib.NewResponseFail(0, "文件上传失败"))
		return
	}

	defer file.Close()
	fname := info.Filename
	pwd, _ := filepath.Abs(`.`)
	println("\n*********************\n\n当前运行环境目录为\n\n***********************", pwd)
	execl := path.Join(pwd, "execl")
	println("\n*********************\n\n当前execl文件夹目录为\n\n***********************", execl)
	execlDir := path.Join(pwd, "execl", fname)
	println("\n*********************\n\n当前execl表格目录为\n\n***********************", execlDir)

	exist, err := PathExists(execl)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return
	}

	if exist {
		fmt.Printf("has dir![%v]\n", execl)
	} else {
		fmt.Printf("no dir![%v]\n", execl)
		// 创建文件夹
		err := os.Mkdir(execl, 0666)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			fmt.Printf("mkdir success!\n")
		}
	}

	//创建一个具有相同名称的文件
	out, errread := os.OpenFile(execlDir,
		os.O_WRONLY|os.O_CREATE, 0666)
	println("yunxing3")
	if errread != nil {

		ctx.JSON(lib.NewResponseFail(1, "文件打开失败"))
		fmt.Println("\n\n\n\n\n\n\n\n文件打开失败的原因\n\n\n\n\n\n\n", errread)

		return
	}
	defer os.RemoveAll(execlDir)

	defer out.Close()

	io.Copy(out, file)

	println("\n打印打印\n", execlDir)

	f, err3 := excelize.OpenFile(execlDir)
	if err3 != nil {
		fmt.Println(err3)
		return
	}
	//// 获取工作表中指定单元格的值
	//cell := f.GetCellValue("Sheet1", "B2")

	//fmt.Println(cell)
	//获取 Sheet1 上所有单元格
	rows := f.GetRows("Sheet1")
	//rows := f.GetRows("Sheet1","A")

	// Get value from cell by given worksheet name and axis.
	//cell, err := f.GetCellValue("Sheet1", "B2")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(cell)

	//db, err := sql.Open("mysql", "xtb_admin:1232@x232C1xz@tcp(rm-uf6p47w3r7nt50s6c0o.mysql.rds.aliyuncs.com:3306)/college_demo")
	db := lib.Engine.DB()

	//查班级
	classSql := `SELECT id,name FROM classes`
	rr, errclass := db.Query(classSql)

	if errclass != nil {
		fmt.Println(errclass)
		return
	}
	for rr.Next() {
		var classId int
		var className string
		_ = rr.Scan(&classId, &className)
		classes[className] = classId
	}

	//查院系
	depSql := `SELECT id,name FROM department`
	rr, errdepartment := db.Query(depSql)

	if errdepartment != nil {
		fmt.Println(errdepartment)
		return
	}
	for rr.Next() {
		var depId int
		var depName string
		_ = rr.Scan(&depId, &depName)
		departments[depName] = depId
	}

	studentNum := 0
	//Get all the rows in the Sheet1.
	//rows := f.GetRows("Sheet1")
	for _, row := range rows {

		studentNum++

		if studentNum == 1 {
			continue
		}

		var info = StudentInfo{
			ID:       row[0],
			realName: row[1],
		}
		if row[2] == "男" {
			info.gender = 1
		} else {
			info.gender = 2
		}
		info.departmentId = departments[row[3]]
		info.classId = classes[row[4]]

		sqlStr += fmt.Sprintf("(%d,%d,'%s',%d,'%s'),", info.departmentId, info.classId, info.realName, info.gender, info.ID)
	}
	strings.TrimRight(sqlStr, ",")
	fmt.Println(sqlStr)
	sqlStr = strings.TrimRight(sqlStr, ",")
	_, err5 := lib.Engine.Query(sqlStr)

	if err5 != nil {
		fmt.Println(err5)
		ctx.JSON(lib.NewResponseFail(1, "表格插入失败"))
		return
	}

	fmt.Printf("添加学生sql %s", sqlStr)

	fmt.Printf("学生总数 %v", studentNum)
	num := strconv.Itoa(studentNum - 1)
	str := "插入学生成功,共插入" + num + "条学生"
	ctx.JSON(lib.NewResponseOK(str))

	return
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func UnixToFormTime(timeStamp int64) string {

	t := time.Unix(int64(timeStamp), 0)
	//返回string
	dateStr := t.Format("2006/01/02 15:04:05")
	return dateStr
}
