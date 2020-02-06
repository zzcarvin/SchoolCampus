package progress

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"time"
)

// 成功返回，返回的大多是必要数据
//swagger:response  responseIndexPlanProgress
type responseIndexPlanProgress struct {
	//返回的结构体
	//in:body
	Body trueBody
}

//首页返回体
type trueBody struct {
	// Required: true
	models.APPResponseType
	Data responPlanProgress
}

type responPlanProgress struct {
	PlanId                    int     `json:"plan_id" xorm:"id"`
	Name                      string  `json:"name" xorm:"name"`    //计划名称
	Date                      string  `json:"date" xorm:"date"`    //计划日期
	DateBegin                 string  `json:"-" xorm:"date_begin"` //计划开始日期
	DateEnd                   string  `json:"-" xorm:"date_end"`
	Types                     int     `json:"types" xorm:"types"`
	BoyTotalDistance          int     `json:"-" xorm:"boy_total_distance"` //计划总里程
	BoySingleMindistance      int     `json:"-" xorm:"boy_single_mindistance"`
	GirlTotalDistance         int     `json:"-" xorm:"girl_total_distance"`
	GirlSingleMindistance     int     `json:"-" xorm:"girl_single_mindistance"`
	BoyWeekDistance           int     `json:"-" xorm:"boy_week_distance"`
	BoyWeekSingleMaxdistance  int     `json:"-" xorm:"boy_week_single_maxdistance"`
	BoyWeekTimes              int     `json:"-" xorm:"comment('男生每周跑步次数') INT(11)"`
	GirlWeekDistance          int     `json:"-" xorm:"girl_week_distance"`
	GirlWeekSingleMaxdistance int     `json:"-" xorm:"girl_week_single_maxdistance"`
	GirlWeekTimes             int     `json:"-" xorm:"comment('女生每周跑步次数') INT(11)"` //女生每周跑步次数
	BoyTotalTimes             int     `json:"-" xorm:"boy_total_times"`
	BoyTimesMindistance       int     `json:"-" xorm:"boy_times_mindistance"`
	GirlTotalTimes            int     `json:"-" xorm:"girl_total_times"`
	GirlTimesMindistance      int     `json:"-" xorm:"girl_times_mindistance"`
	Gender                    int     `json:"gender" xorm:"gender"`
	Duration                  string  `json:"duration" xorm:"duration"` //计划时间段
	DurationBegin             string  `json:"-" xorm:"duration_begin"`
	DurationEnd               string  `json:"-" xorm:"duration_end"`
	StrideFrequency           float32 `json:"-" xorm:"stride_frequency"` //计划步频
	Pace                      int     `json:"-" xorm:"pace"`             //计划配速
	PaceMin                   int     `json:"-" xorm:"pace_min"`
	StrideFrequencyMin        int     `json:"-" xorm:"stride_frequency_min"`
	PlanProgressId            int     `json:"plan_progress_id" xorm:"id"`
	TotalTimes                int     `json:"total_times" xorm:"-"`
	Distance                  int     `json:"progress_distance" xorm:"distance"` //进度总公里数
	ProgressDuration          int     `json:"-" xorm:"duration"`                 //进度总时间
	Times                     int     `json:"times" xorm:"times"`                //进度总次数
	ProgressTimes             int     `json:"progress_times" xorm:"-"`           //当前进度
	Calories                  float32 `json:"-" xorm:"calories"`                 //进度总卡路里,修改
	Steps                     int     `json:"-" xorm:"steps"`                    //进度总步数
	Weeks                     int     `json:"weeks"`
	TargetDistance            string  `json:"target_distance" xorm:"-"` //目标公里数

	//***************************增加男女生配速*************************

	//最小配速
	BoyPaceMin int `json:"boy_pace_min" xorm:"comment('最小配速') INT(11)"`

	//最小配速
	GirlPaceMin int `json:"girl_pace_min" xorm:"comment('最小配速') INT(11)"`

	//*****************************************************************

	PlanCardA1 string `json:"plan_card_a_1" xorm:"-"`
	PlanCardA2 string `json:"plan_card_a_2" xorm:"-"`
	PlanCardB1 string `json:"plan_card_b_1" xorm:"-"`
	PlanCardB2 string `json:"plan_card_b_2" xorm:"-"`
	PlanCardC1 string `json:"plan_card_c_1" xorm:"-"`
	PlanCardC2 string `json:"plan_card_c_2" xorm:"-"`
	PlanCardD1 string `json:"plan_card_d_1" xorm:"-"`
	PlanCardD2 string `json:"plan_card_d_2" xorm:"-"`

	//Duration时间段数组
	//DurationArr []string `json:"duration_arr" xorm:"-"`
}

//请求体
// swagger:parameters  IndexSimplePlan
type swaggerRequestPlan struct {
	// in: body
	Body requestPlan
}
type requestPlan struct {
	PlanId int `json:"plan_id"`
	Gender int `json:"gender"`
}

//返回首页二级页面计划字段
//swagger:response  responseIndexSimplePlan
type shellResponsePlan struct {
	// 简单计划需要的字段
	Body swaggerResponsePlan
}

type swaggerResponsePlan struct {
	// Required: true
	models.APPResponseType
	Data responsePlan
}

type responsePlan struct {
	PlanName     string `json:"plan_name" xorm:"-"`
	PlanDate     string `json:"plan_date" xorm:"-"`
	PlanDuration string `json:"plan_duration" xorm:"-"`
	PlanCardA1   string `json:"plan_card_a_1" xorm:"-"`
	PlanCardA2   string `json:"plan_card_a_2" xorm:"-"`
	PlanCardB1   string `json:"plan_card_b_1" xorm:"-"`
	PlanCardB2   string `json:"plan_card_b_2" xorm:"-"`
	PlanCardC1   string `json:"plan_card_c_1" xorm:"-"`
	PlanCardC2   string `json:"plan_card_c_2" xorm:"-"`
	PlanCardD1   string `json:"plan_card_d_1" xorm:"-"`
	PlanCardD2   string `json:"plan_card_d_2" xorm:"-"`
}

//返回二级页面计划进度字段
//swagger:response responseIndexSimplePlanProgress
type swaggerResponsePlanProgress struct {
	//计划进度返回字段
	Body swaggerResponsePlanProgressShell
}

type swaggerResponsePlanProgressShell struct {
	// Required: true
	models.APPResponseType
	Data responseSimpleProgress
}

type simplePlanProgress struct {
	ProgressDistance      string `json:"distance"`
	Target                string `json:"target"`
	TargetUnit            string `json:"target_unit"`
	StepsValue            string `json:"steps_value"`
	StepsKey              string `json:"steps_key"`
	PaceValue             string `json:"pace_value"`
	PaceKey               string `json:"pace_key"`
	CaloriesValue         string `json:"calories_value"`
	CaloriesKey           string `json:"calories_key"`
	DurationValue         string `json:"duration_value"`
	DurationKey           string `json:"duration_key"`
	ScaleTotalTimes       int    `json:"scale_total_times"`
	ScaleProgressTimes    int    `json:"scale_progress_times"`
	ScaleProgressDistance int    `json:"-"`
}

// swagger:parameters  IndexSimplePlanProgress
type swaggerRequestIndexSimplePlanProgress struct {
	// in: body
	Body requestPlanProgress
}

//简略计划进度请求body
type requestPlanProgress struct {
	PlanId         int `json:"plan_id"`
	PlanProgressId int `json:"plan_progress_id"`
	Gender         int `json:"gender"`
	StudentId      int `json:"student_id"`
}

type responseSimpleProgress struct {
	TermProgress simplePlanProgress
	WeekProgress simplePlanProgress
	WeekExist    bool
}

type nullWeek struct {
}

type responseOnlyTerm struct {
	TermProgress simplePlanProgress
	WeekProgress nullWeek
}

//首页计划

// swagger:route GET  /app/plan/progress APP首页计划和计划进度 IndexPlanProgress
//
// 获取首页计划和计划类型,request请求参数，query params key: id, 参数类型：int类型
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: responseIndexPlanProgress
func planProgress(ctx iris.Context) {

	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	//获取该学生
	var student models.Student
	resPlan, err := lib.Engine.Table("student").Where("id=?", id).Get(&student)
	if err != nil {
		//fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	println("")
	println("after select plan, plan id:", student.PlanId)
	println("")
	if resPlan == false {
		ctx.JSON(lib.NewResponseOK("当前无有效计划"))
		return
	}

	//添加判断，如果学生的计划已终止，直接返回，不展示计划进度
	stuPlanStatus := models.Plan{}
	resStuPlan, err := lib.Engine.Table("plan").Where("id=?", student.PlanId).And("stop=?", 1).Get(&stuPlanStatus)
	if err != nil {
		fmt.Printf("err:%v", err)
		ctx.JSON(lib.NewResponseOK("当前无有效计划"))
		return
	}
	//计划已终止直接返回
	if resStuPlan == false || stuPlanStatus.Stop == 2 {
		ctx.JSON(lib.NewResponseOK("当前无有效计划"))
		return
	}

	fmt.Printf("stuPlan:%v", stuPlanStatus)

	//如果该学生没有计划进度，新建计划进度，
	var progress1 models.PlanProgress
	res, err := lib.Engine.Table("plan_progress").
		Where("plan_progress.student_id=?", id).
		Where("plan_progress.plan_id=?", student.PlanId).
		Desc("plan_progress.create_at").
		Get(&progress1)

	if err != nil {
		//fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	println("")
	println("after select plan_progress, plan_progress id:", progress1.Id, "res:", res)
	println("")
	//没有计划进度，新建计划进度
	if res == false {

		//新建计划进度
		newPlanAndProgress := models.PlanProgress{
			StudentId: int(id),
			PlanId:    student.PlanId,
		}
		numbs, err := lib.Engine.Table("plan_progress").Insert(&newPlanAndProgress)
		if err != nil {
			ctx.JSON(lib.NewResponseFail(1, err.Error()))
			return
		}
		if numbs == 0 {
			ctx.JSON(lib.NewResponseFail(0, "新建计划进度失败"))
			return
		}
	}
	planAndProgress := responPlanProgress{}
	//根据id查询
	b, err := lib.Engine.Table("plan_progress").
		Join("INNER", "plan", "plan.id=plan_progress.plan_id").
		Join("INNER", "student", "student.id=plan_progress.student_id").
		Where("plan_progress.student_id=?", id).
		Where("plan_progress.plan_id=?", student.PlanId).
		Cols("plan.id", "plan.name", "plan.date_begin", "plan.date_end", "plan.total_distance", "plan.min_single_distance",
			"plan.min_week_distance", "plan.max_week_distance", "plan.total_times", "plan.min_pace", "plan.max_pace",
			"plan_progress.id", "plan_progress.distance", "plan_progress.duration", "plan_progress.times", "plan_progress.calories", "plan_progress.steps", "student.gender").
		Get(&planAndProgress)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "该学生信息不完整，无法返回正确信息"))
		return
	}

	//获取计划
	var stuPlan models.Plan
	resPlan, planErr := lib.Engine.Table("plan").Where("id=?", student.PlanId).Get(&stuPlan)
	if planErr != nil {
		//fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//重新获取计划卡片----开始
	//计划返回的参数
	//conf := configs.Conf

	//minDistance := 0 //单次最低

	//minPace := conf.Limit.MinPace
	//minFrequency := conf.Limit.MinFrequency //最低步频

	//运动次数
	planCardA1 := ""
	//计划最低次数
	planCardA2 := "次"
	//周里程
	planCardB1 := ""
	planCardB2 := "周里程"
	//周最低次数
	planCardC1 := ""
	planCardC2 := "周最低"
	//单次最低
	planCardD1 := ""
	planCardD2 := "单次最低"

	//重新获取最低配速
	//minPace=stuPlan.MinPace

	//获取计划中的计量标准和计划中的进度
	planCardA1 = fmt.Sprintf("%d/%d", planAndProgress.Times, stuPlan.TotalTimes) //总跑步次数
	planCardB1 = fmt.Sprintf("%.2fkm", float32(stuPlan.MinWeekDistance)/1000)
	planCardC1 = fmt.Sprintf("%d次", stuPlan.MinWeekTimes)
	planCardD1 = fmt.Sprintf("%.2fkm", float32(stuPlan.MinSingleDistance)/1000)
	//重新获取计划卡片----结束

	//重新获取第几周
	planAndProgress.Weeks = lib.GetWeekSequence(planAndProgress.PlanId, student.Id, time.Now())
	println("第", planAndProgress.Weeks)

	//println("sumDays:",sumDays,"weeks:", planAndProgress.Weeks)
	//跑步区域，因为本来就在学校里，感觉不如去掉

	//添加新卡片的变动数据
	planAndProgress.PlanCardA1 = planCardA1
	planAndProgress.PlanCardA2 = planCardA2
	planAndProgress.PlanCardB1 = planCardB1
	planAndProgress.PlanCardB2 = planCardB2
	planAndProgress.PlanCardC1 = planCardC1
	planAndProgress.PlanCardC2 = planCardC2
	planAndProgress.PlanCardD1 = planCardD1
	planAndProgress.PlanCardD2 = planCardD2

	//添加duration

	//跑步时间段
	//时间段判断
	var timeFrame []models.PlanTimeFrame
	err = lib.Engine.Table("plan_time_frame").Where("plan_id=?", stuPlan.Id).Find(&timeFrame)
	if err != nil {
		fmt.Printf("%v", err.Error())
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//now := time.Now().Format(timeLayout)
	//nowTime, _ := time.ParseInLocation(timeLayout, now, time.Local)
	//nowUnix := nowTime.Unix()
	//firstDuration:=""
	durationStr := ""
	//durationArr:=make([]string, len(timeFrame))
	if len(timeFrame) > 0 {
		//最近时间段 start
		//获取最近的时间段，在时间段内直接显示，不在时间段找下一个最近的时间段
		nowTime := time.Now()
		nowStr := fmt.Sprintf("%d:%d", nowTime.Hour(), nowTime.Minute())
		for _, value := range timeFrame {
			if nowStr >= value.DurationBegin[0:5] && nowStr <= value.DurationEnd[0:5] {
				durationStr = value.DurationBegin[0:5] + "-" + value.DurationEnd[0:5]
				break
			}
		}

		if durationStr == "" {
			tempDurStart := ""
			tempDurEnd := ""

			minDurStart := ""
			minDurEnd := ""
			for _, value := range timeFrame {
				if minDurStart != "" {
					if value.DurationBegin[0:5] < minDurStart {
						minDurStart = value.DurationBegin[0:5]
						minDurEnd = value.DurationEnd[0:5]
					}
				} else {
					minDurStart = value.DurationBegin[0:5]
					minDurEnd = value.DurationEnd[0:5]
				}

				if nowStr < value.DurationBegin[0:5] {
					if tempDurStart != "" && tempDurEnd != "" {
						if value.DurationBegin[0:5] < tempDurStart {
							tempDurStart = value.DurationBegin[0:5]
							tempDurEnd = value.DurationEnd[0:5]
						}
					} else {
						tempDurStart = value.DurationBegin[0:5]
						tempDurEnd = value.DurationEnd[0:5]
					}
				}
			}
			if tempDurStart == "" {
				tempDurStart = minDurStart
				tempDurEnd = minDurEnd
			}
			durationStr = tempDurStart + "-" + tempDurEnd
		}
		//最近时间段 end

	}
	//duration := planAndProgress.DurationBegin[11:16] + "-" + planAndProgress.DurationEnd[11:16]

	//返回新的时间段，在时间段内，返回时间段，不在规定时间段返回最近的时间段
	//now := time.Now().Format(timeLayout)
	//nowTime, _ := time.ParseInLocation(timeLayout, now, time.Local)
	//nowUnix := nowTime.Unix()//当前时间戳
	//
	//for _, value := range timeFrame {
	//	durationStr = value.DurationBegin[0:5] + "-" + value.DurationEnd[0:5] + " "
	//	planDurationStart := now[0:11] + value.DurationBegin
	//	planDurationEnd := now[0:11] + value.DurationEnd
	//
	//
	//
	//}

	planAndProgress.Duration = durationStr
	//planAndProgress.DurationArr=durationArr
	planAndProgress.TotalTimes = stuPlan.TotalTimes
	planAndProgress.ProgressTimes = planAndProgress.Times

	//保持float,单次最低公里数
	planAndProgress.TargetDistance = fmt.Sprintf("%.2f", float32(stuPlan.MinSingleDistance)/1000)

	//获取最近的时间段

	ctx.JSON(lib.NewResponseOK(planAndProgress))

}

//简略计划
func planSimple(ctx iris.Context) {

	//获取计划id和性别
	reqPlan := requestPlan{}
	if err := ctx.ReadJSON(&reqPlan); err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		fmt.Println("plan_id gender ReadJSON error", err)
		return
	}

	plan := models.Plan{}
	resu, err := lib.Engine.Table("plan").Where("id=?", reqPlan.PlanId).Get(&plan)
	if err != nil {
		fmt.Printf("%v", err.Error())
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if resu == false {
		println("没有找到计划")
		ctx.JSON(lib.NewResponseFail(1, "没有运动计划"))
		return
	}

	simpleCard := responsePlan{PlanName: plan.Name}

	//获取变化数据

	//总里程
	simpleCard.PlanCardA1 = fmt.Sprintf("%.2fkm", float32(plan.TotalDistance)/1000)
	simpleCard.PlanCardA2 = "总里程"

	//总次数
	simpleCard.PlanCardB1 = fmt.Sprintf("%d次", plan.TotalTimes)
	simpleCard.PlanCardB2 = "总次数"

	//单次最低里程
	simpleCard.PlanCardD1 = fmt.Sprintf("%.2fkm/次", float64(float64(plan.MinSingleDistance)/1000))
	simpleCard.PlanCardD2 = "单次最低"

	//日最高次数
	simpleCard.PlanCardC1 = fmt.Sprintf("%d次", plan.MaxDayTimes)
	simpleCard.PlanCardC2 = "日最高"

	//每周计划（每周计划的公里数，每周计划的次数乘单次最低公里数）
	//if plan.MinWeekTimes != 0 {
	//	//simpleCard.PlanCardD1 = fmt.Sprintf("%.2fkm", float64(float64(plan.BoyWeekTimes*plan.BoySingleMindistance)/1000))
	//	//修改为每周次数
	//	simpleCard.PlanCardD1 = fmt.Sprintf("%d次", plan.MinWeekTimes)
	//	simpleCard.PlanCardD2 = "每周计划"
	//}

	//获取固定数据
	simpleCard.PlanDate = plan.DateBegin.Format("2006.01.02") + "-" + plan.DateEnd.Format("2006.01.02")

	var timeFrame []models.PlanTimeFrame
	err = lib.Engine.Table("plan_time_frame").Where("plan_id=?", reqPlan.PlanId).Find(&timeFrame)
	if err != nil {
		fmt.Printf("%v", err.Error())
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	var bufferDuration bytes.Buffer
	if len(timeFrame) > 0 {

		for _, value := range timeFrame {
			bufferDuration.WriteString(value.DurationBegin[0:5])
			bufferDuration.WriteString("-")
			bufferDuration.WriteString(value.DurationEnd[0:5])
			bufferDuration.WriteString("")
		}
	}
	simpleCard.PlanDuration = bufferDuration.String()
	ctx.JSON(lib.NewResponseOK(simpleCard))

}

//计划进度详情
//简略计划
// swagger:route POST  /app/plan/detail APP首页二级页面计划进度详情 IndexSimplePlanProgress
//
// 获取首页二级页的简略计划进度
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: responseIndexSimplePlanProgress
func detailPlanProgress(ctx iris.Context) {

	reqPlan := requestPlanProgress{}
	if err := ctx.ReadJSON(&reqPlan); err != nil {
		ctx.JSON(lib.NewResponseFail(1, "计划id或性别格式错误"))
		fmt.Println("plan_id gender ReadJSON error", err)
		return
	}

	var planProgress models.PlanProgress

	res, err := lib.Engine.Table("plan_progress").Where("id=?", reqPlan.PlanProgressId).Get(&planProgress)
	if err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(0, err.Error()))
		return
	}
	if res == false {
		ctx.JSON(lib.NewResponseFail(0, "未查询到计划进度"))
		return
	}
	plan := models.Plan{}
	res2, err := lib.Engine.Table("plan").Where("id=?", reqPlan.PlanId).Get(&plan)
	if err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if res2 == false {
		ctx.JSON(lib.NewResponseFail(1, "未查询到计划"))
		return
	}

	//1.学期跑步进度----开始
	termPlanProgress := simplePlanProgress{StepsKey: "步数", PaceKey: "配速", CaloriesKey: "消耗", DurationKey: "时间"}
	//累积进度
	termPlanProgress.ProgressDistance = fmt.Sprintf("%.2fkm", float32(float32(planProgress.Distance)/1000)) //累计跑量
	//termPlanProgress.StepsValue = fmt.Sprintf("%d", planProgress.Steps)
	termPlanProgress.TargetUnit = "km"
	println("")
	termPlanProgress.StepsValue = fmt.Sprintf("%d", planProgress.Steps) //步数

	fmt.Printf("progress duration:%d,progress distance:%d", planProgress.Duration, planProgress.Distance/1000)
	if planProgress.Duration == 0 || planProgress.Distance == 0 { //当时间或距离为0配速为0，防止返回空
		termPlanProgress.PaceValue = fmt.Sprintf("%d", 0) //配速
		println("当前计划进度，跑步总时长或总距离为0，所以将配速置空")
	} else {
		pace64 := float64(float64(planProgress.Duration) / float64(float64(planProgress.Distance)/1000))
		termPlanProgress.PaceValue = lib.FomPace(pace64) //配速
	}

	termPlanProgress.CaloriesValue = fmt.Sprintf("%dkcal", int32(planProgress.Calories)) //卡路里
	termPlanProgress.DurationValue = SecondsFormHMs(planProgress.Duration)               //累计时间

	//termPlanProgress.Target = fmt.Sprintf("%.2f", float64(float64(plan.TotalDistance)/1000))
	termPlanProgress.ScaleTotalTimes = plan.TotalTimes                                   //总公里int类型
	termPlanProgress.Target = fmt.Sprintf("%d/%d次", planProgress.Times, plan.TotalTimes) //计划次数
	//去掉次数模式下的km单位
	termPlanProgress.TargetUnit = ""

	termPlanProgress.ScaleProgressTimes = planProgress.Times //计划进度
	//进度次数超过最低次数，只显示最低次数
	if termPlanProgress.ScaleProgressTimes >= termPlanProgress.ScaleTotalTimes {
		termPlanProgress.ScaleProgressTimes = termPlanProgress.ScaleTotalTimes
	}

	//学期跑步进度----结束

	//本周跑步进度----开始
	weeRecords := simplePlanProgress{}

	weeRecords.StepsKey = "步数"
	weeRecords.PaceKey = "配速"
	weeRecords.CaloriesKey = "消耗"
	weeRecords.DurationKey = "时间"
	//获取本周运动数据
	records := make([]models.PlanRecord, 0)
	errRecords := lib.Engine.Table("plan_record").
		Where("YEARWEEK( date_format(  create_at,'%Y-%m-%d' ),1 ) = YEARWEEK( now(),1 )").
		And("plan_id=?", reqPlan.PlanId).And("student_id=?", reqPlan.StudentId).And("status=?", 1).
		Find(&records)
	if errRecords != nil {
		fmt.Printf("%v", err.Error())
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	weekDuration := 0
	weekSteps := 0
	weekCalories := 0
	for _, value := range records {
		weeRecords.ScaleProgressTimes = weeRecords.ScaleProgressTimes + 1 //本周公里数累计
		//单次公里超过最大限制使用最大限制
		if plan.MaxSingleDistance != 0 && value.Distance > plan.MaxSingleDistance {
			weeRecords.ScaleProgressDistance = weeRecords.ScaleProgressDistance + plan.MaxSingleDistance
		} else {
			weeRecords.ScaleProgressDistance = weeRecords.ScaleProgressDistance + value.Distance
		}
		weekDuration = weekDuration + value.Duration      //本周时间累计
		weekSteps = weekSteps + value.Steps               //本周步数累计
		weekCalories = weekCalories + int(value.Calories) //本周累计卡路里
	}
	if plan.MaxWeekDistance != 0 && weeRecords.ScaleProgressDistance > plan.MaxWeekDistance {
		weeRecords.ScaleProgressDistance = plan.MaxWeekDistance
	}

	println("本周累计跑步数据：")
	println("累计次数：", weeRecords.ScaleProgressTimes, "累计时长：", weekDuration, "累计步数：", weekSteps, "累计卡路里：", weekCalories)
	//固定数据获取
	//weeRecords.StepsValue = fmt.Sprintf("%d", weekSteps)

	weeRecords.StepsValue = fmt.Sprintf("%d", weekSteps) //步数
	if weeRecords.ScaleProgressDistance != 0 {
		//weeRecords.PaceValue = fmt.Sprintf("%0.1f", float32(weekDuration)/(float32(weeRecords.ScaleProgressDistance)/1000)) //配速
		weeRecords.PaceValue = lib.FomPace(float64(weekDuration) / (float64(weeRecords.ScaleProgressDistance) / 1000))
	} else {
		weeRecords.PaceValue = fmt.Sprintf("%d", 0) //配速
		//weeRecords.PaceValue=lib.FomPace(float64(weekDuration)/(float64(weeRecords.ScaleProgressDistance)/1000))
	}

	weeRecords.CaloriesValue = fmt.Sprintf("%dkcal", weekCalories) //卡路里
	weeRecords.DurationValue = SecondsFormHMs(weekDuration)        //累计时间

	//变化数据开始
	weeRecords.TargetUnit = "次"
	weeRecords.ScaleTotalTimes = plan.MinWeekTimes //总公里int类型
	//进度次数超过最低次数，只显示最低次数
	thisWeekTimes := weeRecords.ScaleProgressTimes
	if thisWeekTimes > plan.MinWeekTimes {
		thisWeekTimes = plan.MinWeekTimes
	}
	weeRecords.Target = fmt.Sprintf("%d/%d", thisWeekTimes, plan.MinWeekTimes) //计划次数

	//周distance
	weeRecords.ProgressDistance = fmt.Sprintf("%.2fkm", float64(weeRecords.ScaleProgressDistance)/1000)

	//本周跑步进度----结束

	responseProgress := responseSimpleProgress{termPlanProgress, weeRecords, true}

	ctx.JSON(lib.NewResponseOK(responseProgress))

}

func SecondsFormHMs(seconds int) string {

	resTime := ""
	hour := seconds / 3600                    //获取小时
	minute := (seconds - hour*3600) / 60      //获取分钟
	second := seconds - hour*3600 - minute*60 //获取秒

	resTime = fmt.Sprintf("%d:%d:%d", hour, minute, second)
	return resTime
}

func (startRunRes *responPlanProgress) String() string {
	b, err := json.Marshal(*startRunRes)
	if err != nil {
		return fmt.Sprintf("%+v", *startRunRes)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *startRunRes)
	}
	return out.String()
}

//func timeStrToUnix(timeStr string) int64{
//
//	time.Now().Unix()
//
//
//	//now := time.Now().Format(timeLayout)
//	//nowTime, _ := time.ParseInLocation(timeLayout, now, time.Local)
//	//nowUnix := nowTime.Unix()//当前时间戳
//
//
//}
