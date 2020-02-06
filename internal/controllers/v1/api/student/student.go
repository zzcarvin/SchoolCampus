package student

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"bytes"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/kataras/iris"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"

	"Campus/configs"
	"github.com/go-xorm/builder"
	"strings"
	//"encoding/json"
	//"fmt"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
	"time"
)

type SportClassInfo struct {
	SportDepartmentId int `json:"sport_department_id" xorm:"department_id"`
	SportClassId      int `json:"sport_class_id" xorm:"class_id"`
}

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

	//---添加学生计划开始
	//假设该学生所在班级有正在运行的计划，将计划更新进该学生的plan_id
	//1.获取正在运行的计划
	var planList []models.Plan
	planType := 0
	errPlans := lib.Engine.Table("plan").Where("stop=?", 1).Find(&planList)
	if errPlans != nil {
		fmt.Printf("查询计划错误：%v", errPlans)
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if len(planList) != 0 {
		for _, value := range planList {
			//planList[value.Id] = value
			if value.Stop == 1 {
				planType = value.PlanType
			}
		}
	}
	//1.1判断计划类型,遍历计划
	for _, value := range planList {
		//行政班级
		if planType == 1 {
			//全部年级，或年级相同
			if value.Year == -1 || student.Year == value.TermYear {
				//全校或院系相同
				if value.DepartmentId == 0 || value.DepartmentId == student.DepartmentId {
					if value.DepartmentId == 0 {
						student.PlanId = value.Id
					} else if value.DepartmentId == student.DepartmentId {
						//全系
						if value.ClassId == 0 || student.ClassId == value.ClassId {
							if value.Gender == 0 || value.Gender == student.Gender {
								//全班
								if value.Gender == 0 {
									student.PlanId = value.Id
								} else if value.Gender == student.Gender {
									student.PlanId = value.Id
								}
							}
						}
					}
				}
			}
		} else {
			//体育班，添加学生没有体育班字段，暂时不管，只更新行政班计划，如果要更新体育班计划，使用many_update的思路。

		}

	}
	//-------添加学生计划结束

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
	id := ctx.Params().GetUint64Default("id", 0)
	println("id:", id)
	student := models.Student{}

	//解析department
	err := ctx.ReadJSON(&student)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//补充获取每个学生的计划id----开始

	//获取所有计划
	var planList []models.Plan
	planType := 0
	lib.Engine.Table("plan").Find(&planList)
	if len(planList) != 0 {
		for _, value := range planList {
			if value.Stop == 1 {
				planType = value.PlanType
			}
		}
	}

	println("plan_type:", planType)

	var studentList []models.Student
	lib.Engine.Table("student").Find(&studentList)
	if len(studentList) != 0 {
	}

	//获取student_class
	var studentClass []models.StudentClass
	lib.Engine.Table("student_class").Find(&studentClass)
	if len(studentClass) != 0 {
	}

	if planType != 0 {
		if planType == 1 {
			//行政班
			stuClassId := student.ClassId
			stuGender := student.Gender
			student.PlanId = getPlanIdByClaGen(stuClassId, stuGender, studentList)
		} else if planType == 2 {
			//体育班
			if student.SportClassId != 0 {
				student.PlanId = getPlanIdBySportCla(student.SportClassId, studentClass, studentList)
			} else {
				student.PlanId = 0
			}

		}
	} else {
		student.PlanId = 0
	}

	println("plan_id:", student.PlanId)

	//补充获取每个学生的计划id----结束

	//插入数据
	session := lib.Engine.NewSession()
	defer session.Close()
	err1 := session.Begin()
	if err1 != nil {
		ctx.JSON(lib.NewResponseFail(1, "事务开启失败"))
		println("事务开启失败")
		return
	}

	_, err2 := lib.Engine.Table("student").ID(id).Update(student)
	if err2 != nil {
		fmt.Printf("err2:%v", err2)
		session.Rollback()
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if student.SportClassId != 0 && student.SportDepartmentId != 0 {
		hasClass, errClass := lib.Engine.Table("classes").Where("id=?", student.SportClassId).And("department_id=?", student.SportDepartmentId).Exist()
		if errClass != nil {
			session.Rollback()
			ctx.JSON(lib.NewResponseFail(1, err.Error()))
			return
		}
		if !hasClass {
			session.Rollback()
			ctx.JSON(lib.NewResponseFail(1, "找不到体育班"))
			return
		}

		hasId := 0
		has, err := lib.Engine.Table("student_class").Where("student_id=?", id).Cols("id").Get(&hasId)
		if has {
			_, err3 := lib.Engine.Exec("update student_class set class_id=?,department_id=? where student_id=?", student.SportClassId, student.SportDepartmentId, id)
			if err3 != nil {
				session.Rollback()
				ctx.JSON(lib.NewResponseFail(1, err.Error()))
				return
			}
		} else {
			_, err3 := lib.Engine.Exec("insert into student_class(class_id,department_id,student_id) values(?,?,?)", student.SportClassId, student.SportDepartmentId, id)
			if err3 != nil {
				session.Rollback()
				ctx.JSON(lib.NewResponseFail(1, err.Error()))
				return
			}
		}

	}
	session.Commit()
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
		Cols("student.id", "student.class_id", "student.department_id", "student.code", "student.name", "classes.name", "department.name", "student.gender", "student.create_at", "student.cellphone", "student.year").
		Get(&student)
	// SELECT `student`.`id`, `student`.`code`, `student`.`name`, `classes`.`name`, `department`.`name`, `student`.`gender`, `student`.`create_at`, `student`.`cellphone`
	// FROM `student` INNER JOIN classes ON classes.id=student.class_id INNER JOIN department ON department.id=student.department_id WHERE (student.id=?) LIMIT 1 []interface {}{0x1}

	sportClassInfo := SportClassInfo{}
	_, err = lib.Engine.Table("student_class").
		Where("student_id=?", id).
		Cols("department_id,class_id").
		Get(&sportClassInfo)
	student.SportDepartmentId = sportClassInfo.SportDepartmentId
	student.SportClassId = sportClassInfo.SportClassId

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
type pageStudentInfo struct {
	Id           int    `json:"id" xorm:"autoincr id pk" `
	Year         int    `json:"year" xorm:"year"`
	DepartmentId int    `json:"department_id" xorm:"department_id"`
	ClassId      int    `json:"class_id" xorm:"class_id"`
	Name         string `json:"name" xorm:"name"`
	Gender       int    `json:"gender" xorm:"gender"`
	Code         string `json:"code" xorm:"code"`
	Department   string `json:"department" xorm:"department_name"`
	Class        string `json:"class" xorm:"class_name"`
	SportClass   string `json:"sport_class"`
	Cellphone    string `json:"cellphone" xorm:"cellphone"`
	Device       string `json:"device" xorm:"device"`
	Face         string `json:"face" xorm:"face"`
}

func search(ctx iris.Context) {
	cfg := configs.Conf.Database
	//创建查询Session指针
	query := lib.Engine.Table("student").Select("student.*,classes.name as class_name,department.name as department_name,login_device as device,face,"+
		"(select name from classes where id=(select class_id from student_class where student_id=student.id)) as sport_class").
		Join("inner", "department", "department.id=student.department_id").
		Join("inner", "classes", "classes.id=student.class_id")

	//字段查询
	if ctx.URLParamExists("school_id") {
		schoolId := ctx.URLParam("school_id")
		query.Join("left", "(select login_device,face,school_user_id from "+cfg.CoreDBName+".user where "+cfg.CoreDBName+".user.school_id="+schoolId+" and "+cfg.CoreDBName+".user.school_user_type=1) as core_user", "core_user.school_user_id=student.id")
	} else {
		ctx.JSON(lib.NewResponseFail(1, "学校编码不能为空"))
		return
	}

	if ctx.URLParamExists("name") {
		query.And(builder.Like{"student.name", ctx.URLParam("name")})
	}
	if ctx.URLParamExists("classid") {
		query.And(builder.Like{"student.class_id", ctx.URLParam("classid")})
	}
	if ctx.URLParamExists("departmentid") {
		query.And(builder.Like{"student.department_id", ctx.URLParam("departmentid")})
	}
	if ctx.URLParamExists("code") {
		query.And(builder.Like{"student.code", ctx.URLParam("code")})
	}
	if ctx.URLParamExists(("userid")) {
		query.And(builder.Like{"student.user_id", ctx.URLParam("userid")})
	}
	if ctx.URLParamExists("cellphone") {
		query.And(builder.Like{"student.cellphone", ctx.URLParam("cellphone")})
	}
	if ctx.URLParamExists("classes_name") {
		classesname := ctx.URLParam("classes_name")
		query.Where("classes.name =?", classesname)
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
	if page <= 0 {
		page = 1
	}
	query.Limit(size, (page-1)*size)

	//查询
	var dataList []pageStudentInfo
	counts, err := query.FindAndCount(&dataList)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	var pageModel = models.Page{
		List:  dataList,
		Total: counts,
	}

	ctx.JSON(lib.NewResponseOK(pageModel))
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
	size := ctx.URLParamIntDefault("size", 0)
	query.Limit(size, page*size)

	now := time.Now()

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

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
		//fmt.Printf("planProgressDistance%v/planDisatance%v=progress%v", float32(planProgress[index].PlanProgressDistance), float32(planProgress[index].PlanDistance), planProgress[index].Progress)
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

func EveryMonthProgress(finishRun models.RequestFinishRun1) (a []weekContent) {

	//创建查询Session指针
	queryPlan := lib.Engine.Table("plan")
	queryRecord := lib.Engine.Table("plan_record")

	if finishRun.PlanId != 0 {
		queryPlan.And("id=?", finishRun.PlanId)
	} else {
		fmt.Println(lib.NewResponseFail(0, "缺少planId"))
		return
	}
	if finishRun.StudentId != 0 {
		queryRecord.And("student_id=?", finishRun.StudentId)
	} else {
		fmt.Println(lib.NewResponseFail(0, "缺少studentId"))
		return
	}
	println("获取的planId:", finishRun.PlanId, "studentId:", finishRun.StudentId)
	//1.获取计划，包括计划的开始时间，结束时间，计划的周数，计划的每周开始时间和结束时间
	var plan models.Plan
	res, err := queryPlan.Get(&plan)
	if err != nil {
		println("")
		fmt.Printf("%v", err)
		fmt.Println(lib.NewResponseFail(0, err.Error()))
		return
	}
	if res == false {

		fmt.Println(lib.NewResponseFail(0, "未查到校园计划"))
		return
	}

	//获取学生信息
	student := models.Student{}
	resStudent, err := lib.Engine.Table("student").Where("id=?", finishRun.StudentId).Get(&student)
	if err != nil {
		fmt.Printf("%v", err)
		//ctx.JSON(lib.NewResponseFail(0, err.Error()))
		fmt.Println(lib.NewResponseFail(0, err.Error()))
		return
	}
	if resStudent == false {
		println("未查到学生信息")
		//ctx.JSON(lib.NewResponseFail(0, "未查到学生信息"))
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
		And("student_id=?", finishRun.StudentId).And("plan_id=?", finishRun.PlanId).And("status=?", 1).
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
		planWeekContents = weekContents[0:1]

		return planWeekContents
	case 2:
		planWeekContents = append(planWeekContents, weekContents[0])
		planWeekContents = append(planWeekContents, weekContents[1])
		//ctx.JSON(lib.NewResponseOK(planWeekContents))
		return planWeekContents
	case 3:
		planWeekContents = weekContents[0 : len(weekContents)-1]
		//ctx.JSON(lib.NewResponseOK(planWeekContents))
		return planWeekContents
	}
	return planWeekContents

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

type StudentInfo struct {
	ID           string
	realName     string
	gender       int
	departmentId int
	classId      int
	Year         int
}

type ClassInfo struct {
	Id   int
	Name string
}

type DepartmentInfo struct {
	Id   int
	Name string
	//departmentType int
}

func excel(ctx iris.Context) {
	// Get the file from the request.
	_, info, err := ctx.FormFile("file")
	ext := path.Ext(info.Filename)
	fileName := uuid.NewV4().String() + ext
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(lib.NewResponseFail(0, "文件上传失败"))
		return
	}
	lib.SaveUploadedFile(info, "./excel", fileName)

	f, err := excelize.OpenFile("./excel/" + fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove("./excel/" + fileName)

	students := make(map[string]models.Student)
	classes := make(map[string]int)
	departments := make(map[string]int)
	//查班级
	var classInfoList []ClassInfo
	lib.Engine.Table("classes").Cols("id", "name").Find(&classInfoList)
	if len(classInfoList) != 0 {
		for _, value := range classInfoList {
			classes[value.Name] = value.Id
		}
	}
	//查院系
	var departmentList []DepartmentInfo
	lib.Engine.Table("department").Cols("id", "name").Find(&departmentList)
	if len(departmentList) != 0 {
		for _, value := range departmentList {
			departments[value.Name] = value.Id
		}
	}

	//查学生
	var studentList []models.Student
	lib.Engine.Table("student").Find(&studentList)
	if len(studentList) != 0 {
		for _, value := range studentList {
			students[value.Code] = value
		}
	}

	//假设该学生所在班级有正在运行的计划，将计划更新进该学生的plan_id
	//1.获取正在运行的计划
	var planList []models.Plan
	planType := 0
	errPlans := lib.Engine.Table("plan").Where("stop=?", 1).Find(&planList)
	if errPlans != nil {
		fmt.Printf("查询计划错误：%v", errPlans)
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if len(planList) != 0 {
		for _, value := range planList {
			//planList[value.Id] = value
			if value.Stop == 1 {
				planType = value.PlanType
			}
		}
	}

	updateNum := 0
	createAt := time.Now().Format("2006-01-02 15:04:05")
	sqlStr := "insert into student(department_id,class_id,name,gender,code,create_at,year,plan_id) values"
	rows := f.GetRows("Sheet1")

	IdList := make(map[string]int)
	rowIndex := make(map[string]int) //列名对应的索引，这样导入的时候不用担心列名顺序的问题
	for _, row := range rows {
		updateNum++
		if updateNum == 1 {
			//第一行根据列名导入
			var columns = []string{"学号", "姓名", "性别", "院系", "班级", "入学年份"}
			rowIndex = lib.GetRowIndex(row, columns)
			if len(rowIndex) != len(columns) {
				ctx.JSON(lib.NewResponseFail(1, "请确认所有列名是否正确"))
				return
			}
			continue
		}
		//学号已存在数据库直接返回
		if students[row[rowIndex["学号"]]].Code == row[rowIndex["学号"]] {
			fmt.Println("重复学号：" + row[rowIndex["学号"]])
			ctx.JSON(lib.NewResponseFail(1, "学号："+row[rowIndex["学号"]]+"已存在"))
			return
		}
		//excel中数据，学号重复
		if IdList[row[rowIndex["学号"]]] == 1 {
			fmt.Println("重复学号：" + row[rowIndex["学号"]])
			ctx.JSON(lib.NewResponseFail(1, "学号："+row[rowIndex["学号"]]+"重复"))
			return
		}

		var info = StudentInfo{
			ID:       row[rowIndex["学号"]],
			realName: row[rowIndex["姓名"]],
		}
		if row[rowIndex["性别"]] == "男" {
			info.gender = 1
		} else {
			info.gender = 2
		}
		info.departmentId = departments[row[rowIndex["院系"]]]
		info.classId = classes[row[rowIndex["班级"]]]
		info.Year = lib.Str2int(row[rowIndex["入学年份"]])
		IdList[info.ID] = 1
		stuPlanId := 0

		//1.1判断计划类型,遍历计划
		for _, value := range planList {
			//行政班级
			if planType == 1 {
				//全部年级，或年级相同
				if value.Year == -1 || info.Year == value.TermYear {
					//全校或院系相同
					if value.DepartmentId == 0 || value.DepartmentId == info.departmentId {
						if value.DepartmentId == 0 {
							stuPlanId = value.Id
						} else if value.DepartmentId == info.departmentId {
							//全系
							if value.ClassId == 0 || info.classId == value.ClassId {
								if value.Gender == 0 || value.Gender == info.gender {
									//全班
									if value.Gender == 0 {
										stuPlanId = value.Id
									} else if value.Gender == info.gender {
										stuPlanId = value.Id
									}
								}
							}
						}
					}
				}
			} else {
				//体育班，添加学生没有体育班字段，暂时不管，只更新行政班计划，如果要更新体育班计划，使用many_update的思路。

			}

		}
		//-------添加学生计划结束

		sqlStr += fmt.Sprintf("(%d,%d,'%s',%d,'%s','%s',%d,%d),", info.departmentId, info.classId, info.realName, info.gender, info.ID, createAt, info.Year, stuPlanId)
	}

	//TODO: 如果和库里比对，有重复的学号要提示
	sqlStr = strings.TrimRight(sqlStr, ",")
	fmt.Printf("添加学生sql %s", sqlStr)
	num := strconv.Itoa(updateNum - 1)
	str := "导入学生成功,共添加" + num + "条学生"
	fmt.Println(str)
	_, err = lib.Engine.Query(sqlStr)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(lib.NewResponseFail(1, "导入失败"))
		return
	}
	ctx.JSON(lib.NewResponseOK(str))
	return
}

type UpdateModel struct {
	Code         string
	DepartmentId int    `json:"department_id"`
	ClassId      int    `json:"class_id"`
	Name         string `json:"name"`
	Gender       int    `json:"gender"`
}

//批量修改学生
type ManyUpdateStu struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Gender       string `json:"gender"`
	ClassId      string `json:"class_id"`
	DepartmentId string `json:"department_id"`
	//SportClassId string  `json:"sport_class_id"`
	PlanId string `json:"plan_id"`
}

//只传入体育班
type OnlySportClassUpdateStu struct {
	Code   string `json:"code"`
	PlanId string `json:"plan_id"`
}

//更新导入
func UpdateExcel(ctx iris.Context) {

}

func UpdateBatch(tableName string, models []UpdateModel) {
	var columns []string
	var sqlStr string
	var whereIn string
	columns = lib.GetFieldName(models[0])
	referenceColumn := columns[0]

	if tableName != "" {
		sqlStr += "UPDATE " + tableName + " SET "
		for _, column := range columns {
			if column != referenceColumn {
				sqlStr += column + " = CASE "
				for _, vv := range models {
					t := reflect.ValueOf(vv)
					referenceValue := lib.GetReferenceValue(t, referenceColumn)
					sqlStr += fmt.Sprintf("WHEN %s = '%s' THEN '%s' ", strings.ToLower(referenceColumn), referenceValue, t.FieldByName(column))
				}
				sqlStr += "ELSE " + column + " END, "
			}
		}
		for _, vv := range models {
			t := reflect.ValueOf(vv)
			whereIn += fmt.Sprintf("'%s', ", lib.GetReferenceValue(t, referenceColumn))
		}
		sqlStr = fmt.Sprintf("%s WHERE %s IN(%s)", strings.TrimRight(sqlStr, ", "), referenceColumn, strings.TrimRight(whereIn, ", "))

	}

	fmt.Println(sqlStr)
}

//学生批量更新
//1.获取excel
//2.获取基本信息，学生表，院系表，班级表
//3.将excel表数据整理成golang数据
//4.更新student表
//5.删除旧的student_class表数据,插入新的student_class数据
//备注：学生班级不正确，返回不更新，只更新体育班级存在的学生的体育班级表
//导入excel，格式：学号	姓名	性别	院系	班级	体育班 。使用sheet1
//补充功能：添加学生的计划id

//2019.12.10突然要改逻辑
//补充逻辑：
//1.院系，班级，性别可以不填，不填时表示不更新这些信息。
func manyUpdateExcel(ctx iris.Context) {
	//1.获取excel
	// Get the file from the request.
	_, info, err := ctx.FormFile("file")
	ext := path.Ext(info.Filename)
	fileName := uuid.NewV4().String() + ext
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(lib.NewResponseFail(0, "文件上传失败"))
		return
	}
	lib.SaveUploadedFile(info, "./excel", fileName)

	f, err := excelize.OpenFile("./excel/" + fileName)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(lib.NewResponseFail(0, "文件上传失败"))
		return
	}
	defer os.Remove("./excel/" + fileName)
	//2.获取基本信息，学生表，院系表，班级表
	students := make(map[string]models.Student)
	classes := make(map[string]models.Classes)

	departments := make(map[string]int)
	//查学生
	var studentList []models.Student
	lib.Engine.Table("student").Find(&studentList)
	if len(studentList) != 0 {
		for _, value := range studentList {
			students[value.Code] = value
		}
	}
	//查班级
	var classList []models.Classes
	lib.Engine.Table("classes").Find(&classList)
	if len(classList) != 0 {
		for _, value := range classList {
			classes[value.Name] = value
		}
	}
	//查院系
	var departmentList []DepartmentInfo
	lib.Engine.Table("department").Cols("id", "name").Find(&departmentList)
	if len(departmentList) != 0 {
		for _, value := range departmentList {
			departments[value.Name] = value.Id
		}
	}

	//获取所有计划
	var planList []models.Plan
	planType := 0
	lib.Engine.Table("plan").Find(&planList)
	if len(planList) != 0 {
		for _, value := range planList {
			//planList[value.Id] = value
			if value.Stop == 1 {
				planType = value.PlanType
			}
		}
	}

	println("plan_type:", planType)

	//获取student_class
	var studentClass []models.StudentClass
	lib.Engine.Table("student_class").Find(&studentClass)
	if len(studentClass) != 0 {
		//for _, value := range studentClass {
		//	studentClass[value.Id] = value
		//}
	}

	//3.将excel表数据整理成golang数据
	updateNum := 0
	createAt := time.Now().Format("2006-01-02 15:04:05")
	rows := f.GetRows("Sheet1")

	//2.将excel表数据整理成golang数据
	//IdList := make(map[string]int)
	rowIndex := make(map[string]int) //列名对应的索引，这样导入的时候不用担心列名顺序的问题
	println("row len:", len(rows))
	if len(rows) == 0 {
		fmt.Println(err)
		ctx.JSON(lib.NewResponseFail(0, "文件上传失败,文件内容为空"))
		return
	}

	//如果当前没有计划
	println("当前没有计划")
	if planType == 0 {
		if rows[1][rowIndex["体育班"]] != "" {
			planType = 2
		} else {
			planType = 1
		}
	}
	println("最后要更新的计划类型：", planType)

	//根据行政班级是否为空选择不同结构体

	studentsUpdate := make([]ManyUpdateStu, len(rows))
	studentsOnlySportUpdate := make([]OnlySportClassUpdateStu, len(rows))
	//删除语句
	delStuClaWhere := ""

	//插入语句
	sqlStuClaInsertStr := "insert into student_class(student_id,department_id,class_id) values"

	if len(rows[1][4]) != 0 {
		//原版全部更新
		for indexRow, row := range rows {
			updateNum++
			if updateNum == 1 {
				//第一行根据列名导入
				var columns = []string{"学号", "姓名", "性别", "院系", "班级", "体育班"}
				rowIndex = lib.GetRowIndex(row, columns)
				continue
			}
			studentsUpdate[indexRow] = ManyUpdateStu{
				Code: row[rowIndex["学号"]],
				Name: row[rowIndex["姓名"]],
			}
			if row[rowIndex["性别"]] == "男" {
				studentsUpdate[indexRow].Gender = "1"
			} else {
				studentsUpdate[indexRow].Gender = "2"
			}

			//获取每个学生的计划id
			if planType != 0 {
				if planType == 1 {
					//行政班
					stuClassId, _ := strconv.Atoi(studentsUpdate[indexRow].ClassId)
					stuGender, _ := strconv.Atoi(studentsUpdate[indexRow].Gender)
					studentsUpdate[indexRow].PlanId = strconv.Itoa(getPlanIdByClaGen(stuClassId, stuGender, studentList))
				} else if planType == 2 {
					//体育班
					if row[rowIndex["体育班"]] != "" {
						studentsUpdate[indexRow].PlanId = strconv.Itoa(getPlanIdBySportCla(classes[row[rowIndex["体育班"]]].Id, studentClass, studentList))
						println("planId:", studentsUpdate[indexRow].PlanId)
					} else {
						studentsUpdate[indexRow].PlanId = "0"
					}

				}
			} else {
				studentsUpdate[indexRow].PlanId = "0"
			}

			studentsUpdate[indexRow].DepartmentId = strconv.Itoa(departments[row[rowIndex["院系"]]])
			studentsUpdate[indexRow].ClassId = strconv.Itoa(classes[row[rowIndex["班级"]]].Id)
			fmt.Printf("(%s,'%s',%s,%s,'%s','%s'),", studentsUpdate[indexRow].Code, studentsUpdate[indexRow].Name, studentsUpdate[indexRow].Gender, studentsUpdate[indexRow].DepartmentId, createAt, studentsUpdate[indexRow].PlanId)
			//教师id
			//teachId := updateTeachers[teachers[indexRow].Name]

			//班级id
			//row[rowIndex["班级"]]
			//去掉班级验证
			//println("class_id:", classes[row[rowIndex["班级"]]].Id, "class_name:", row[rowIndex["班级"]])
			//if classes[row[rowIndex["班级"]]].Id == 0 {
			//	println("文件上传失败，" + row[rowIndex["班级"]] + "不存在，请先插入班级")
			//	ctx.JSON(lib.NewResponseFail(0, "文件上传失败，班级："+row[rowIndex["班级"]]+"不存在，请先插入班级"))
			//	return
			//}

			//获取需要更新体育班的学号
			println("体育班：", row[rowIndex["体育班"]])
			if row[rowIndex["体育班"]] != "" {
				delStuClaWhere += fmt.Sprintf("'%s', ", row[rowIndex["学号"]])
			}

			//体育班级id,存在体育班时进行插入学生体育班级对应表
			if classes[row[rowIndex["体育班"]]].Id != 0 {

				stuId := students[row[rowIndex["学号"]]].Id
				sportClaId := classes[row[rowIndex["体育班"]]].Id
				sportDepId := classes[row[rowIndex["体育班"]]].DepartmentId
				println("学生id:", stuId, "体育班id:", sportClaId, "体育班院系：", sportDepId)
				sqlStuClaInsertStr += fmt.Sprintf("(%d,%d,%d),", stuId, sportDepId, sportClaId)
			}

		}
	} else {

		//只更新体育班
		for indexRow, row := range rows {
			updateNum++
			if updateNum == 1 {
				//第一行根据列名导入
				var columns = []string{"学号", "姓名", "性别", "院系", "班级", "体育院系", "体育班"}
				rowIndex = lib.GetRowIndex(row, columns)
				continue
			}
			studentsOnlySportUpdate[indexRow] = OnlySportClassUpdateStu{
				Code: row[rowIndex["学号"]],
			}

			//获取每个学生的计划id
			if planType != 0 {
				//TODO 缺少学号等字段验证，缺少这些字段返回，不允许更新
				println("体育班id：", classes[row[rowIndex["体育班"]]].Id)
				//体育班
				if row[rowIndex["体育班"]] != "" {
					//studentsOnlySportUpdate[indexRow].PlanId = strconv.Itoa(getPlanIdBySportCla(classes[row[rowIndex["体育班"]]].Id, studentClass, studentList))
					//println("planId:", studentsOnlySportUpdate[indexRow].PlanId)
					//重新获取体育班id
					stuDepartmentId := departments[row[rowIndex["体育院系"]]]
					for _, value := range classList {
						if value.Name == row[rowIndex["体育班"]] && value.DepartmentId == stuDepartmentId {
							//如果没有计划正在执行，应该返回

							studentsOnlySportUpdate[indexRow].PlanId = strconv.Itoa(getPlanIdBySportCla(value.Id, studentClass, studentList))
						}
					}
				} else {
					studentsOnlySportUpdate[indexRow].PlanId = "0"
				}

			} else {
				studentsOnlySportUpdate[indexRow].PlanId = "0"
			}

			//获取需要更新体育班的学号
			println("体育班：", row[rowIndex["体育班"]])
			if row[rowIndex["体育班"]] != "" {
				delStuClaWhere += fmt.Sprintf("'%s', ", row[rowIndex["学号"]])
			}

			//体育班级id,存在体育班时进行插入学生体育班级对应表
			if classes[row[rowIndex["体育班"]]].Id != 0 {

				stuId := students[row[rowIndex["学号"]]].Id
				sportClaId := classes[row[rowIndex["体育班"]]].Id
				sportDepId := classes[row[rowIndex["体育班"]]].DepartmentId
				println("学生id:", stuId, "体育班id:", sportClaId, "体育班院系：", sportDepId)
				sqlStuClaInsertStr += fmt.Sprintf("(%d,%d,%d),", stuId, sportDepId, sportClaId)
			}

		}
	}

	//println(sqlTeaClaInsertStr)
	//4.更新student表
	sqlUpTeacherStr := ""
	if len(rows[1][4]) != 0 {
		studentUpdatesInfo := studentsUpdate[1:len(studentsUpdate)]
		sqlUpTeacherStr = lib.UpdateBatch("student", studentUpdatesInfo)
	} else {
		//只更新体育班
		studentsOnlySportUpdateInfo := studentsOnlySportUpdate[1:len(studentsOnlySportUpdate)]
		sqlUpTeacherStr = lib.UpdateBatch("student", studentsOnlySportUpdateInfo)
	}
	//studentUpdatesInfo := studentsUpdate[1:len(studentsUpdate)]
	//fmt.Printf("%v", studentUpdatesInfo)
	//sqlUpTeacherStr := lib.UpdateBatch("student", studentUpdatesInfo)
	println(sqlUpTeacherStr)
	//5.更新student_class表
	//5.1删除旧的student_class表数据
	println("获取需要更新体育班的学号")

	println(delStuClaWhere)
	delSqlStr := fmt.Sprintf("DELETE student_class FROM student_class INNER JOIN student ON student.id=student_class.student_id WHERE student.code IN(%s)", strings.TrimRight(delStuClaWhere, ", "))
	println(delSqlStr)

	//5.2插入新的student_class数据
	println("插入student_class")
	println(sqlStuClaInsertStr)
	insertStuCla := strings.TrimRight(sqlStuClaInsertStr, ", ")

	// 创建 Session 对象
	sess := lib.Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err = sess.Begin(); err != nil {
		fmt.Println(err)
		ctx.JSON(lib.NewResponseFail(1, "更新失败"))
		return
	}

	if _, err = sess.Query(sqlUpTeacherStr); err != nil {
		sess.Rollback()
		fmt.Println(err)
		ctx.JSON(lib.NewResponseFail(1, "更新失败"))
		return
	} else if _, err = sess.Query(delSqlStr); err != nil {
		sess.Rollback()
		fmt.Println(err)
		ctx.JSON(lib.NewResponseFail(1, "更新失败"))
		return
	} else if _, err = sess.Query(insertStuCla); err != nil {
		sess.Rollback()
		fmt.Println(err)
		ctx.JSON(lib.NewResponseFail(1, "更新失败"))
		return
	}

	// 完成事务
	sess.Commit()

	ctx.JSON(lib.NewResponseOK("修改成功"))
	return

}

func UnixToFormTime(timeStamp int64) string {

	t := time.Unix(int64(timeStamp), 0)
	//返回string
	dateStr := t.Format("2006/01/02 15:04:05")
	return dateStr
}

func completenumber(ctx iris.Context) {

	tomorrowtmstring := daytimesring(1)
	tmstring := daytimesring(0)
	year := ctx.URLParamIntDefault("year", 0)
	month := ctx.URLParamIntDefault("month", 0)
	day := ctx.URLParamIntDefault("day", 0)
	if year == 1 {
		return

	}
	if month == 1 {

	}
	for i := 0; i < 26; i++ {

		if day == 1 {
			numsum, err := lib.Engine.Table("plan_progress").Cols("status").Where("create_time BETWEEN '?' AND '?'", tomorrowtmstring, tmstring).SumInt("plan_record", "status")
			if err != nil {
				ctx.JSON(lib.NewResponseFail(1, err.Error()))
				return
				lib.NewResponseOK(numsum)
			}

		}

	}

}

//func zerotime() (zerotime int64) {
//	timeStr := time.Now().Format("2006-01-02")
//	t, _ := time.Parse("2006-01-02", timeStr)
//	fmt.Println(t.Format(time.UnixDate))
//	//Unix返回早八点的时间戳，减去8个小时
//	timestamp := t.UTC().Unix() - 8*3600
//	fmt.Println("timestamp:", timestamp)
//	return timestamp
//
//}

func daytimesring(i int) (tmstring string) {

	year := time.Now().Format("2006") //年
	month := time.Now().Format("01")  //月
	day := time.Now().Day() + i       //日

	tm2, _ := time.Parse("01/02/2006", month+"/"+strconv.Itoa(day)+"/"+year)
	tmstring = tm2.Format("2006-01-02")

	fmt.Println("---------------", tm2)

	return tmstring

}

func getPlan(ctx iris.Context) {
	plan := models.Plan{}

	//取URL参数 id
	studentId := ctx.Params().GetUint64Default("id", 0)
	if studentId == 0 {
		ctx.JSON(lib.NewResponseFail(1, "参数错误"))
		return
	}

	b, err := lib.Engine.Table("plan").
		Join("INNER", "student", "plan.id=student.plan_id").Where("student.id=?", studentId).Get(&plan)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该计划"))
		return
	}
	var timeFrame []models.PlanTimeFrame
	err = lib.Engine.Table("plan_time_frame").Where("plan_id=?", plan.Id).Find(&timeFrame)
	if err != nil {
		fmt.Printf("%v", err.Error())
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	plan.TimeFrame = timeFrame

	ctx.JSON(lib.NewResponseOK(plan))

}

//行政班 通过 class_id，gender来查询plan_id
func getPlanIdByClaGen(classId int, gender int, students []models.Student) int {
	for _, value := range students {
		if value.ClassId == classId && value.Gender == gender {
			return value.PlanId
		}
	}
	return 0
}

//体育班获取plan_id
func getPlanIdBySportCla(sportClassId int, studentClass []models.StudentClass, students []models.Student) int {
	for _, value := range studentClass {
		if value.ClassId == sportClassId {
			println("studentId:", value.StudentId)
			planId := getPlanIdByStudentId(value.StudentId, students)
			println("value.ClassId:", value.ClassId)
			return planId
		}
	}
	return 0
}

//只通过体育班获取id获取计划id
//func getPlanIdOnlySportClaId(sportClassId int) int{
//
//}

//通过student_id获取plan_id
func getPlanIdByStudentId(studentId int, students []models.Student) int {
	println("student_id:", studentId)
	for _, value := range students {

		if value.Id == studentId {
			println("student-id:", studentId)
			return value.PlanId
		}
	}
	return 0
}
