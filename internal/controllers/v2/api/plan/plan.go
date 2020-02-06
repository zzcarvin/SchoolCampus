package plan

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"encoding/json"
	"fmt"
	"github.com/go-xorm/builder"
	"github.com/kataras/iris"
	"strings"
	"time"
)

// swagger:parameters  PlanCreateRequest
type PlanCreateRequest struct {
	// in: body
	Body models.Plan
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

	//用于储存时间段
	PlanTimeFrame [][]string `json:"plantimeframe" xorm:"-"`
}

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

	querystudent := lib.Engine.Table("student")
	students := []models.Student{}

	year := ctx.URLParamIntDefault("year", 0)

	departmentid := ctx.URLParamIntDefault("department_id", 0)
	classid := ctx.URLParamIntDefault("class_id", 0)
	gender := ctx.URLParamIntDefault("gender", 0)
	//updatesql :="UPDATE `student` SET `plan_id` = 100 WHERE `id` in"

	if year != 0 {
		querystudent.Where("year=?", year)
	} else {
		ctx.JSON(lib.NewResponseFail(1, "年份不存在"))

	}

	if departmentid != 0 {
		querystudent.And("department_id=?", departmentid)
	}

	if classid != 0 {
		querystudent.And("class_id=?", classid)
	}

	if gender != 0 {
		querystudent.And("gender=?", gender)
	}
	errcount := querystudent.Cols("id", "plan_id").Find(&students)
	if errcount != nil {
		ctx.JSON(lib.NewResponseFail(1, errcount.Error()))
		return
	}

	var id []int
	for _, student := range students {
		id = append(id, student.Id)
	}
	fmt.Println("\n\n\n\n", id)
	for _, student := range students {
		if student.PlanId != 0 {
			ctx.JSON(lib.NewResponseFail(1, "不能创建该计划，该计划与学校其他计划冲突"))
			return

		}
	}
	//插入数据
	res1, err2 := lib.Engine.Table("plan").Insert(&plan)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	student := models.Student{
		PlanId: plan.Id,
	}

	res, errupdate := querystudent.ID(id).Cols("plan_id").Update(student)
	if errupdate != nil {
		ctx.JSON(lib.NewResponseFail(1, errupdate.Error()))
		return
	}

	planid := plan.Id

	planTimeFrames := plan.PlanTimeFrame
	planTimeFramestruct := []models.PlanTimeFrame{}
	for i, planTimeFrame := range planTimeFrames {
		planTimeFramestruct[i].PlanId = planid
		planTimeFramestruct[i].DurationBegin = planTimeFrame[0]
		planTimeFramestruct[i].DurationEnd = planTimeFrame[1]
	}
	resrframe, errframe := lib.Engine.Table("plan_time_frame").Insert(&planTimeFramestruct)
	if errframe != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	fmt.Println("\n\n\n\n共插入了", resrframe, "条记录\n\n\n")

	lib.NewResponseOK(res1)
	lib.NewResponseOK(res)
	ctx.JSON(lib.NewResponseOK(res))

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

func update(ctx iris.Context) {
	// swagger:route PUT /api/plan/:id  plan PlanUpdateRequest
	//
	// 修改计划
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: PlanUpdateResponse
	//取URL参数 id

	id := ctx.Params().GetUint64Default("id", 0)

	plan := models.Plan{}

	//解析plan
	err := ctx.ReadJSON(&plan)
	b2, err := json.MarshalIndent(plan, "", "   ")
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Println("\n\n\nplan结构体\n\n", string(b2))
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//TODO 验证数据有效性

	//为了修复学生每周计划进度，必须将计划终止日期修改为当天23：59：59
	finDateEndUnix := plan.DateEnd.Unix() + 24*3600 - 1
	finDateEnd := time.Unix(int64(finDateEndUnix), 0)
	plan.DateEnd = finDateEnd

	//插入数据
	res, err2 := lib.Engine.Table("plan").ID(id).AllCols().Update(plan)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(res))
}

// swagger:route DELETE /api/plan/:id  plan PlanDelete
//
// 删除计划
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: Response
func remove(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	plan := models.Plan{}

	//根据id查询
	affected, err := lib.Engine.Table("plan").ID(id).Delete(&plan)
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

	plan := models.Plan{}
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
	size := ctx.URLParamIntDefault("size", 0)
	query.Limit(size, page*size)

	//查询
	var plan []models.Plan
	err := query.Find(&plan)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
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
