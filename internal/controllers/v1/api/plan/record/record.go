package record

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"fmt"
	"github.com/gookit/validate"
	"github.com/kataras/iris"
	"strings"
	"time"
)

// swagger:parameters  RecordCreateRequest
type RecordCreateRequest struct {
	// in: body
	Body models.PlanRecord
}

// 响应结构体
//
// swagger:response    RecordCreateResponse
type RecordCreateResponse struct {
	// in: body
	Body RecordresponseMessage
}
type RecordresponseMessage struct {
	models.ResponseType
	Data models.PlanRecord
}

func create(ctx iris.Context) {
	// swagger:route POST /api/plan/Record Record RecordCreateRequest
	//
	// 创建围栏
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: RecordCreateResponse
	planRecord := models.PlanRecord{}
	err := ctx.ReadJSON(&planRecord)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))

	}
	v := validate.Struct(planRecord)
	// v := validate.New(u)

	if v.Validate() { // 验证成功
		// do something ...
	} else {

		fmt.Println(v.Errors)                   // 所有的错误消息
		fmt.Println(v.Errors.One())             // 返回随机一条错误消息
		fmt.Println(v.Errors.Field("[]points")) // 返回该字段的错误消息
		ctx.JSON(lib.NewResponseFail(1, "points字段过长"))
		return
	}

	//ctx.JSON(lib.NewResponseOK(planRecord))

	res, err := lib.Engine.Table("plan_record").Insert(&planRecord)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	lib.NewResponseOK(res)
	ctx.JSON(lib.NewResponseOK(planRecord))
}

// swagger:route DELETE /api/plan/Record Record RecordDelete
//
//	 删除围栏
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: Response
func remove(ctx iris.Context) {
	id := ctx.Params().GetUint64Default("id", 0)
	planRecord := models.PlanRecord{}
	affected, err := lib.Engine.Table("plan_record").ID(id).Delete(&planRecord)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	ctx.JSON(lib.NewResponseOK(affected))
}

// swagger:parameters  RecordUpdateRequest
type RecordUpdateRequest struct {
	// in: body
	Body models.PlanRecord
}

// 响应结构体
//
// swagger:response    RecordUpdateResponse
type RecordUpdateResponse struct {
	// in: body
	Body RecordresponseMessage
}

func update(ctx iris.Context) {
	// swagger:route put /api/plan/Record/:id Record RecordUpdateRequest
	// 修改围栏
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//     Responses:
	//       200: RecordUpdateResponse

	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	planRecord := models.PlanRecord{}

	//解析student
	err := ctx.ReadJSON(&planRecord)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//TODO 验证数据有效性

	//先获取该运动记录，如果之前运动记录的状态与修改的状态一样，直接返回，不做修改
	oldRecord := models.PlanRecord{}
	resRecord, err2 := lib.Engine.Table("plan_record").ID(id).Get(&oldRecord)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if resRecord == false {
		ctx.JSON(lib.NewResponseFail(1, "查询不到该运动记录"))
		return
	}
	println("就运动记录的状态：", oldRecord.Status, "要修改的状态：", planRecord.Status)
	if planRecord.Status == oldRecord.Status {
		ctx.JSON(lib.NewResponseFail(1, "状态未发生变化，不做修改"))
		return
	}

	//插入数据
	res, err2 := lib.Engine.Table("plan_record").ID(id).AllCols().Update(&planRecord)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err2.Error()))
		return
	}

	if res == 0 {
		println("该条运动记录更新失败")
		ctx.JSON(lib.NewResponseFail(1, "该条运动记录更新失败"))
		return
	}

	//将运动记录修改为有效后，将该运动记录的数据添加进计划进度

	//先获取该运动进度
	progress := models.PlanProgress{}
	respro, err := lib.Engine.Table("plan_progress").Where("plan_id=?", planRecord.PlanId).And("student_id=?", planRecord.StudentId).Get(&progress)
	if err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if respro == false {
		println("学生计划进度查询失败")
		ctx.JSON(lib.NewResponseFail(1, "未查询到该运动记录的计划进度"))
		return
	}

	println("原先计划进度里程：", progress.Distance)

	//更新运动计划进度
	planProgress := models.PlanProgress{}
	//无效改有效
	if planRecord.Status == 1 {
		println("无效改有效")
		planProgress.Distance = planRecord.Distance + progress.Distance
		planProgress.Duration = planRecord.Duration + progress.Duration
		planProgress.Calories = planRecord.Calories + progress.Calories
		planProgress.Steps = planRecord.Steps + progress.Steps
		planProgress.Times = progress.Times + 1
	} else { //有效改有效
		println("有效改无效")
		planProgress.Distance = progress.Distance - planRecord.Distance
		planProgress.Duration = progress.Duration - planRecord.Duration
		planProgress.Calories = progress.Calories - planRecord.Calories
		planProgress.Steps = progress.Steps - planRecord.Steps
		planProgress.Times = progress.Times - 1
	}

	fmt.Printf("原先的计划进度：%v", progress)
	println("")
	fmt.Printf("当前的计划进度%v", planProgress)
	resprogress, err := lib.Engine.Table("plan_progress").Where("plan_id=?", planRecord.PlanId).And("student_id=?", planRecord.StudentId).Update(&planProgress)
	if err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if resprogress == 1 {
		println("学生计划进度更新成功")
	}
	if resprogress != 1 {
		println("计划进度更新失败")
		ctx.JSON(lib.NewResponseFail(1, "计划进度更新失败"))
	}

	println("修改后的计划进度里程：", planProgress.Distance)
	ctx.JSON(lib.NewResponseOK(resprogress))
}

// swagger:route GET /api/plan/Record/:id Record RecordGet
//
// 查询围栏
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//     Responses:
//       200: Response
func get(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	planRecord := models.PlanRecord{}
	//根据id查询
	b, err := lib.Engine.Table("plan_record").ID(id).Get(&planRecord)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该记录"))
		return
	}
	ctx.JSON(lib.NewResponseOK(planRecord))
}

// swagger:route GET /api/plan/Records Record RecordSearch
//
// 查询围栏
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//     Responses:
//       200: Response
func search(ctx iris.Context) {
	//创建查询Session指针
	query := lib.Engine.Table("plan_record")

	//字段查询
	//根据名称查询记录()

	if ctx.URLParamExists("start_time") {
		query.Where("start_time >= ? AND end_time <= ?", ctx.URLParam("start_time"), ctx.URLParam("end_time"))

	}
	if ctx.URLParamExists("code") || ctx.URLParamExists("student_name") {
		code := ctx.URLParam("code")
		name := ctx.URLParam("student_name")
		query.Join("INNER", "student", "student.id=plan_record.student_id").
			Where("code=?", code).Or("name =?", name).
			Cols("plan_record.*")
	}
	if ctx.URLParamExists("classes_name") {
		classesname := ctx.URLParam("classes_name")
		query.
			Join("INNER", "student", "student.id=plan_record.student_id").
			Join("INNER", "classes", "student.class_id=classes.id").
			Where("classes.name =?", classesname).
			Cols("plan_record.*")
	}

	if ctx.URLParamExists("student_id") {
		student_id, err := ctx.URLParamInt("student_id")
		if err != nil {
			ctx.JSON(lib.NewResponseFail(1, err.Error()))
			return
		}
		query.Where("student_id = ?", student_id)

	}
	if ctx.URLParamExists("status") {
		status := ctx.URLParam("status")

		query.And("status = ?", status)

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
	var planRecord []models.PlanRecord
	counts, err := query.FindAndCount(&planRecord)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	var pageModel = models.Page{
		List:  planRecord,
		Total: counts,
	}

	ctx.JSON(lib.NewResponseOK(pageModel))
}

//储存折线图所需信息的结构体
type TotalDayRun struct {
	Milestotal int64   `json:"Milestotal"`
	Time1      float64 `json:"time1"`
	Time2      float64 `json:"time2"`
	Time3      float64 `json:"time3"`
	Time4      float64 `json:"time4"`
	Time5      float64 `json:"time5"`
	Time6      float64 `json:"time6"`
	Timetotal  float64 `json:"Timetotal"`
}

//包换两种逻辑，且数据量小，直接展示，不用循环
func totaldayrun(ctx iris.Context) {
	record := models.PlanRecord{}
	TotalDayRun := TotalDayRun{}
	now := time.Now() //获取当前

	// 时间

	year, month, day := now.Date()                                   //截取当前时间的年月日
	today_str1 := fmt.Sprintf("%d-%d-%d 00:00:00", year, month, day) //用截取的时间生成所需的时刻
	today_str2 := fmt.Sprintf("%d-%d-%d 04:00:00", year, month, day)
	today_str3 := fmt.Sprintf("%d-%d-%d 08:00:00", year, month, day)
	today_str4 := fmt.Sprintf("%d-%d-%d 12:00:00", year, month, day)
	today_str5 := fmt.Sprintf("%d-%d-%d 16:00:00", year, month, day)
	today_str6 := fmt.Sprintf("%d-%d-%d 20:00:00", year, month, day)
	today_str7 := fmt.Sprintf("%d-%d-%d 00:00:00", year, month, day+1)

	//
	//zerotime1 :=zerotime
	//zerotime = zerotime +30*60
	//tm1 := time.Unix(zerotime1, 0)
	//time1:=tm1.Format("2006-01-02 15:04:05")
	//tm2:= time.Unix(zerotime, 0)
	//time2:=tm2.Format("2006-01-02 15:04:05")
	//sum,err := query.Where("end_time > ? AND end_time <= ?",time1,time2).Sum(&record,"status")

	//today_time := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	time1, err := lib.Engine.Table("plan_record").Where("start_time > ? AND end_time <= ? AND status = 1", today_str1, today_str2).Sum(&record, "distance") //0-4点的学生跑步总时长
	TotalDayRun.Time1 = time1
	time2, err := lib.Engine.Table("plan_record").Where("start_time > ? AND end_time <= ? AND status = 1", today_str2, today_str3).Sum(&record, "distance")
	TotalDayRun.Time2 = time2
	time3, err := lib.Engine.Table("plan_record").Where("start_time > ? AND end_time <= ? AND status = 1", today_str3, today_str4).Sum(&record, "distance")
	TotalDayRun.Time3 = time3
	time4, err := lib.Engine.Table("plan_record").Where("start_time > ? AND end_time <= ? AND status = 1", today_str4, today_str5).Sum(&record, "distance")
	TotalDayRun.Time4 = time4
	time5, err := lib.Engine.Table("plan_record").Where("start_time > ? AND end_time <= ? AND status = 1", today_str5, today_str6).Sum(&record, "distance")
	TotalDayRun.Time5 = time5
	time6, err := lib.Engine.Table("plan_record").Where("start_time > ? AND end_time <= ? AND status = 1", today_str6, today_str7).Sum(&record, "distance")
	TotalDayRun.Time6 = time6
	Timetotal, err := lib.Engine.Table("plan_record").Where("start_time > ? AND end_time <= ? AND status = 1", today_str1, today_str7).Sum(&record, "times")
	TotalDayRun.Timetotal = Timetotal
	Milestotal, err := lib.Engine.Table("plan_record").Where("start_time > ? AND end_time <= ? AND status = 1", today_str1, today_str7).SumInt(&record, "distance") //一天内总距离
	TotalDayRun.Milestotal = Milestotal
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	fmt.Println("1111111111", TotalDayRun)
	ctx.JSON(lib.NewResponseOK(TotalDayRun))
}

func getPlanRecordString(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	planRecord := models.PlanRecord{}

	//根据id查询
	b, err := lib.Engine.Table("plan_record").Where("id=?", id).Get(&planRecord)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该记录"))
		return
	}
	var PlanRecordString []string
	PlanRecordString = make([]string, 0)
	for _, value := range planRecord.InvalidCode {

		if value == 1 {
			PlanRecordString = append(PlanRecordString, "配速过低")
		}
		if value == 2 {
			PlanRecordString = append(PlanRecordString, "配速过高")
		}
		if value == 3 {
			PlanRecordString = append(PlanRecordString, "步频过低")
		}
		if value == 4 {
			PlanRecordString = append(PlanRecordString, "配速过高")
		}
		if value == 5 {
			PlanRecordString = append(PlanRecordString, "超过24小时")
		}
		if value == 6 {
			PlanRecordString = append(PlanRecordString, "没有经过所有打卡点")
		}
		if value == 7 {
			PlanRecordString = append(PlanRecordString, "未达到最低公里数")
		}
		if value == 8 {
			PlanRecordString = append(PlanRecordString, "人脸认证失败")
		}
		if value == 9 {
			PlanRecordString = append(PlanRecordString, "没有经过蓝牙打卡点")
		}
		if value == 10 {
			PlanRecordString = append(PlanRecordString, "时间太长")
		}
		if value == 11 {
			PlanRecordString = append(PlanRecordString, "时间太短")
		}
		planRecord.PlanRecordString = PlanRecordString

	}

	ctx.JSON(lib.NewResponseOK(planRecord.PlanRecordString))
}
