package record

import (
	"Campus/configs"
	"Campus/internal/lib"
	"Campus/internal/models"
	"fmt"
	"github.com/go-xorm/builder"
	"github.com/kataras/iris"
	"math"
	"strconv"
	"strings"
	"time"
)

type responseRecord struct {
	Id        int     `json:"id" xorm:"autoincr id"`
	SchoolId  int     `json:"school_id" xorm:"school_id"`
	PlanId    int     `json:"plan_id" xorm:"plan_id"`
	Name      string  `json:"name" xorm:"name"`
	StudentId int     `json:"student_id" xorm:"student_id"`
	Type      int     `json:"type" xorm:"type"`
	StartTime string  `json:"start_time" xorm:"start_time"`
	EndTime   string  `json:"end_time" xorm:"end_time"`
	Distance  int     `json:"distance" xorm:"distance"`
	Duration  int     `json:"duration" xorm:"duration"`
	Calories  float64 `json:"calories" xorm:"calories"`
	Steps     int     `json:"steps" xorm:"steps"`
	Pace      int     `json:"pace" xorm:"pace"`
	FormPace  string  `json:"form_pace" xorm:"form_pace"`
	//Points       []models.Points `json:"points" xorm:"points"`
	Points           []models.Points `json:"points" xorm:"points"`
	Frequency        float64         `json:"frequency" xorm:"frequency"`
	Frequencies      []int           `json:"frequencies" xorm:"frequencies"`
	CreateAt         string          `json:"create_at" xorm:"create_at created"`
	XFrequencies     []int           `json:"x_frequencies" xorm:"-"`
	XNumber          []int           `json:"x_number" xorm:"-"`
	Status           int             `json:"status" xorm:"status"`
	InvalidCode      []int           `json:"invalid_code" xorm:"invalid_code"`
	PlanRecordString []string        `json:"planrecordstring"  xorm:"-"`

	//新加字段
	DistanceObject  recordReasonFlo `json:"distance_object" xorm:"-"`
	DurationObject  recordReasonStr `json:"duration_object" xorm:"-"`
	PaceObject      recordReasonStr `json:"pace_object" xorm:"-"`
	StepPageObject  recordReasonInt `json:"step_page_object" xorm:"-"`
	FrequencyObject recordReasonInt `json:"frequency_object" xorm:"-"`
}

//请求体
// swagger:parameters  recordsDuration
type swaggerRequestDuration struct {
	// in: body
	Body newRequestDuration
}

type newRequestDuration struct {
	Type      int `json:"type"`
	PlanId    int `json:"plan_id"`
	StudentId int `json:"student_id"`
}

type requestDuration struct {
	Type   int `json:"type"`
	PlanId int `json:"plan_id"`
}

//返回周月年数据展示信息
//swagger:response  shellResponseRecordsDuration
type shellResponseRecordsDuration struct {
	// 简单计划需要的字段
	Body swaggerResponseRecordsDuration
}

type swaggerResponseRecordsDuration struct {
	// Required: true
	models.APPResponseType
	Data newResponseDuration
}

type newResponseDuration struct {
	//Records     []models.PlanRecord `json:"records"`
	StartDate    string    `json:"start_date"`
	EndDate      string    `json:"end_date"`
	SumDistance  int       `json:"sumDistance"`
	SumDuration  int       `json:"sumDuration"`
	Sumtimes     int       `json:"sumTimes"`
	SumColories  float64   `json:"sumColories"`
	AveragePace  int       `json:"averagePace"`
	XData        []float32 `json:"x_data"` //x轴数据
	XAxis        []string  `json:"x_axis"` //x轴
	YAxis        []float64 `json:"y_axis"` //y轴
	YUnit        string    `json:"y_unit"`
	DurationType string    `json:"duration_type"`
}

type strResponseDuration struct {
	StartDate    string    `json:"start_date"`
	EndDate      string    `json:"end_date"`
	SumDistance  string    `json:"sumDistance"`
	SumDuration  string    `json:"sumDuration"`
	Sumtimes     string    `json:"sumTimes"`
	SumColories  string    `json:"sumColories"`
	AveragePace  string    `json:"averagePace"`
	XData        []float32 `json:"x_data"` //x轴数据
	XAxis        []string  `json:"x_axis"` //x轴
	YAxis        []float64 `json:"y_axis"` //y轴
	YUnit        string    `json:"y_unit"`
	DurationType string    `json:"duration_type"`
}

type responseDuration struct {
	//Records     []models.PlanRecord `json:"records"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	SumDistance int     `json:"sumDistance"`
	SumDuration int     `json:"sumDuration"`
	Sumtimes    int     `json:"sumTimes"`
	SumColories float64 `json:"sumColories"`
	AveragePace int     `json:"averagePace"`
	XNumbers    []int   `json:"xNumbers"`
}

type continueDays struct {
	StudentId    int `json:"student_id" xorm:"student_id"`
	ContinueDays int `json:"Continue_days" xorm:"continue_days"`
}

type responseBest struct {
	Record       models.PlanRecord `json:"record"`
	SumDistance  int               `json:"sumDistance"`
	ContinueDays int               `json:"continue_days"`
}

type responseShare struct {
	Name            string    `json:"name" xorm:"name"`
	Sequence        int       `json:"sequence" xorm:"-"`
	AllDistance     int       `json:"all_distance" xorm:"-"`
	LongestDistance int       `json:"longest_distance" xorm:"distance"`
	LongestDuration int       `json:"longest_duration" xorm:"duration"`
	MaxCalories     float32   `json:"bigest_calories" xorm:"calories"`
	FastestPace     int       `json:"fastest_pace" xorm:"pace"`
	CreateAt        time.Time `json:"-" xorm:"create_at"`
	Now             string    `json:"now" xorm:"-"`
}

//返回个人所有运动记录
//swagger:response  swaggerResponseRecords
type swaggerResponseRecords struct {
	Body swaggerShellResponseRecords
}

type swaggerShellResponseRecords struct {
	models.APPResponseType
	Data []responseRecord
}

type responseDetail struct {
	Points []point `json:"points"`
}

type point struct {
	Id        int     `json:"-" xorm:"id"`
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
	Status    bool    `json:"status" xorm:"-"`
}

type recordReasonInt struct {
	RawData int    `json:"raw_data"`
	Reason  string `json:"reason"`
	Status  bool   `json:"status"`
}

type recordReasonStr struct {
	RawData string `json:"raw_data"`
	Reason  string `json:"reason"`
	Status  bool   `json:"status"`
}

type recordReasonFlo struct {
	RawData float64 `json:"raw_data"`
	Reason  string  `json:"reason"`
	Status  bool    `json:"status"`
}

func get(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)
	print("id:", id)
	PlanRecord := models.PlanRecord{}
	//根据id查询,这里用ID()函数出错：ID condition is error, expect 0 primarykeys, there are 1
	b, err := lib.Engine.Table("plan_record").Where("id=?", id).Get(&PlanRecord)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该运动记录"))
		return
	}
	ctx.JSON(lib.NewResponseOK(PlanRecord))
}

//获取多条记录
//新版一段时间内的运动记录
// swagger:route POST  /app/plan/record/search APP个人所有运动记录 recordsSearch
//
// 获取个人的所有运动数据，请求参数,URLParam student_id
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: swaggerResponseRecords
func search(ctx iris.Context) {

	//创建查询Session
	query := lib.Engine.Table("plan_record")

	//字段查询
	if ctx.URLParamExists("student_id") {
		query.And(builder.Like{"student_id", ctx.URLParam("student_id")})
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
	//page := ctx.URLParamIntDefault("page", 0)
	//size := ctx.URLParamIntDefault("size", 50)
	//query.Limit(size, page*size)
	page := ctx.URLParamIntDefault("page", 0)
	size := ctx.URLParamIntDefault("size", 0)
	query.Limit(size, page*size)

	//获取计划名称
	query.Join("INNER", "plan", "plan.id=plan_record.plan_id")
	//查询
	var planRecord []responseRecord
	err := query.Find(&planRecord)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//处理步频数组，在运动详细记录里展示步频折线图
	for index, _ := range planRecord {
		xFrequencies := make([]int, 0)
		x_times := make([]int, 0)
		if len(planRecord[index].Frequencies) <= 5 {
			//for i1, _ := range planRecord[index].Frequencies {
			//	xFrequencies = append(xFrequencies, planRecord[index].Frequencies[i1])
			//}
			for i := 0; i < 5; i++ {
				if len(planRecord[index].Frequencies) > i {
					xFrequencies = append(xFrequencies, planRecord[index].Frequencies[i])
				} else {
					xFrequencies = append(xFrequencies, 0)
				}
				x_times = append(x_times, i)
			}
			planRecord[index].XFrequencies = xFrequencies
			planRecord[index].XNumber = x_times
			println("")
			//fmt.Printf("展示的xFrequenceies:%v", xFrequencies)
		} else {
			//获取步频数组长度
			freLength := len(planRecord[index].Frequencies)
			xIndex := 1
			for i2, _ := range planRecord[index].Frequencies {
				println("i2:", i2, "(freLength/4)*xIndex):", (freLength/4)*xIndex, "xIndex:", xIndex)
				if i2 == 0 {
					//planRecord[index].XNumber[xIndex] = i2
					x_times = append(x_times, i2)
					xFrequencies = append(xFrequencies, planRecord[index].Frequencies[i2])
				} else if (i2 == ((freLength / 4) * xIndex)) && xIndex <= 4 {
					//planRecord[index].XNumber[xIndex] = i2
					x_times = append(x_times, i2)                                          //x轴
					xFrequencies = append(xFrequencies, planRecord[index].Frequencies[i2]) //y轴
					xIndex++
				}
			}
			planRecord[index].XFrequencies = xFrequencies
			planRecord[index].XNumber = x_times
			println("")
			//fmt.Printf("展示的xFrequenceies:%v", xFrequencies)
		}
	}

	//步频折线图第二版-----使用每段时间内的步数除分钟数
	//思路：找到每段结束，用循环往前加i个

	//fmt.Printf("展示的xFrequenceies:%v", xFrequencies)

	ctx.JSON(lib.NewResponseOK(planRecord))
}

//最佳运动记录
//没有运动记录时返回初始化运动记录
func best(ctx iris.Context) {

	//获取student_id
	id := ctx.Params().GetUint64Default("id", 0)
	print("student_id:", id)

	PlanRecord := models.PlanRecord{}
	sumDistance, err := lib.Engine.Table("plan_record").Where("student_id=?", id).Sum(&PlanRecord, "distance")
	if err != nil {
		println("查询运动距离之和出错：")
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, "运动记录错误"))
		return
	}
	//根据id查询

	recordNumbers := 1
	//重新排序寻找需要的最佳记录----开始

	//单次最长里程
	distanceRecord := models.PlanRecord{}
	bDistance, err := lib.Engine.Table("plan_record").Where("student_id=?", id).Desc("distance").Get(&distanceRecord)
	if err != nil {
		println("查询最佳距离出错：")
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, "运动记录错误"))
		return
	}
	if bDistance == false {
		//ctx.JSON(lib.NewResponseFail(0, "未找到最佳公里数运动记录"))
		recordNumbers = 0
	}
	//单次最长时间
	durationRecord := models.PlanRecord{}
	bduration, err := lib.Engine.Table("plan_record").Where("student_id=?", id).Desc("duration").Get(&durationRecord)
	if err != nil {
		println("查询最佳时长出错：")
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, "运动记录错误"))
		return
	}
	if bduration == false {
		//ctx.JSON(lib.NewResponseFail(0, "未找到最佳运动时长运动记录"))
		recordNumbers = 0
	}

	//单次最佳配速
	paceRecord := models.PlanRecord{}
	bpace, err := lib.Engine.Table("plan_record").Where("student_id=?", id).Desc("calories").Get(&paceRecord)
	if err != nil {
		println("查询最佳卡路里出错：")
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, "运动记录错误"))
		return
	}
	if bpace == false {
		//ctx.JSON(lib.NewResponseFail(0, "未找到最佳配速运动记录"))
		recordNumbers = 0
	}

	//单次最多消耗卡路里。没要卡路里？？？
	caloryRecord := models.PlanRecord{}
	bcalory, err := lib.Engine.Table("plan_record").Where("student_id=?", id).Asc("pace").Get(&caloryRecord)
	if err != nil {
		println("查询最佳配速出错：")
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, "运动记录错误"))
		return
	}
	if bcalory == false {
		//ctx.JSON(lib.NewResponseFail(0, "未找到最佳卡路里运动记录"))
		recordNumbers = 0
	}
	//单次最佳速度
	speedRecord := models.PlanRecord{}
	bspeed, err := lib.Engine.Table("plan_record").Where("student_id=?", id).Desc("speed").Get(&speedRecord)
	if err != nil {
		println("查询最佳速度出错：")
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, "运动记录错误"))
		return
	}
	if bspeed == false {
		//ctx.JSON(lib.NewResponseFail(0, "未找到最佳速度运动记录"))
		recordNumbers = 0
	}
	//单次最佳速度
	stepsRecord := models.PlanRecord{}
	bsteps, err := lib.Engine.Table("plan_record").Where("student_id=?", id).Desc("steps").Get(&stepsRecord)
	if err != nil {
		println("查询最佳步数出错：")
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, "运动记录错误"))
		return
	}
	if bsteps == false {
		//ctx.JSON(lib.NewResponseFail(0, "未找到最佳速度运动记录"))
		recordNumbers = 0
	}

	//重新排序寻找需要的最佳记录----结束

	//整理最佳运动记录
	bestRecord := models.PlanRecord{}
	//有运动记录
	if recordNumbers != 0 {
		bestRecord = models.PlanRecord{
			Distance: distanceRecord.Distance,
			Duration: durationRecord.Duration,
			Calories: caloryRecord.Calories,
			Pace:     paceRecord.Pace,
			Steps:    stepsRecord.Steps,
			Speed:    speedRecord.Speed,
		}
	}

	//获取连续运动天数

	student := models.Student{}
	bContinue, err := lib.Engine.Table("student").Where("id=?", id).Get(&student)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, "查询学生错误"))
		return
	}
	if bContinue == false {
		ctx.JSON(lib.NewResponseFail(1, "查询学生失败"))
		return
	}

	//for _,value:=range bl{
	//	println(len(value))
	//	fmt.Printf("%#v", value)
	//	println(value["continue_days"])
	//}
	//days,_:=strconv.Atoi(bl[0]["continue_days"])
	responseBestRecord := responseBest{
		Record:       bestRecord,
		SumDistance:  int(sumDistance),
		ContinueDays: student.Continue,
	}
	//返回距离最远的相关记录记录
	ctx.JSON(lib.NewResponseOK(responseBestRecord))

}

//新版一段时间内的运动记录
// swagger:route POST  /app/plan/record/newduration APP数据展示页 recordsDuration
//
// 获取周月年的运动统计数据
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: shellResponseRecordsDuration
func newDuration(ctx iris.Context) {

	var requestDurations newRequestDuration
	if err := ctx.ReadJSON(&requestDurations); err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		fmt.Print("%v", err)
		return
	}

	//获取student_id
	id := requestDurations.StudentId

	//验证
	valType := lib.ValidateRequest(requestDurations)
	if valType == false {
		ctx.JSON(lib.NewResponseFail(1, "时间类型格式错误"))
		return
	}

	//获取当前时间和一周前时间
	now := time.Now()
	beforeTime := ""
	var cycle int //周期
	var startDate string
	var endDate string
	xLength := 0 //x轴长度
	durationType := ""
	switch requestDurations.Type { //获取x轴的长度，和开始的时间
	case 1:
		cycle = int(now.Weekday()) //周期天数
		if cycle == 0 {
			cycle = 7
		}
		startTime := now.AddDate(0, 0, -(cycle - 1))
		beforeTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location()).Format("2006-01-02")
		endTime1 := now.AddDate(0, 0, 7-cycle)
		endDate = time.Date(endTime1.Year(), endTime1.Month(), endTime1.Day(), 0, 0, 0, 0, startTime.Location()).Format("2006年1月2日")
		startDate = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location()).Format("2006年1月2日")
		xLength = 7 //x轴长度
		durationType = "周统计"
		println("")
		fmt.Printf("week,startTime:%v,endTime:%v", beforeTime, endDate)
	case 2:
		cycle = int(now.Day())
		startTime := now.AddDate(0, 0, -(cycle - 1))
		beforeTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location()).Format("2006-01-02")
		startDate = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location()).Format("2006年1月2日")
		endTime1 := now.AddDate(0, 0, 0)
		endDate = time.Date(endTime1.Year(), endTime1.Month()+1, 0, 0, 0, 0, 0, startTime.Location()).Format("2006年1月2日")
		xLength = 5 //x轴长度
		durationType = "月统计"
		println("")
		fmt.Printf("month,startTime:%v,endTime:%v", beforeTime, endDate)
	case 3:
		cycle = int(now.Month())
		startTime := now.AddDate(0, -(cycle - 1), 0)
		beforeTime = time.Date(startTime.Year(), startTime.Month(), 0, 0, 0, 0, 0, startTime.Location()).Format("2006-01-02 15:04:05")
		endTime1 := now.AddDate(0, 0, 0)
		endDate = time.Date(endTime1.Year(), endTime1.Month()+1, 0, 0, 0, 0, 0, startTime.Location()).Format("2006年")
		startDate = endDate
		if cycle <= 6 {
			xLength = 6
		} else {
			xLength = cycle
		}
		durationType = "年统计"
		println("")
		fmt.Printf("month,startTime:%v,endTime:%v", beforeTime, endDate)
	default:
		ctx.JSON(lib.NewResponseFail(1, "时间类型格式错误"))
		return
	}
	println("今天是第", int(now.Month()), "月，周：", int(now.Weekday()), "日：", int(now.Day()))

	//本周到现在的每天的运动里程，本月到现在，每天的运动里程，本年到现在，每月的运动里程
	PlanRecord := []models.PlanRecord{}
	//根据id查询，And("status=?",1). 作弊数据也计入统计
	err := lib.Engine.Table("plan_record").
		Where("student_id=?", id).
		And("plan_id=?", requestDurations.PlanId).
		And(builder.Between{"create_at", beforeTime, now.Format("2006-01-02 15:04:05")}).Asc("create_at").
		Find(&PlanRecord)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	recordslens := len(PlanRecord)
	//if recordslens == 0 {
	//	//没有运动数据返回空数组
	//	nums := make([]float32, 0)
	//	dataRecoed := responseDuration{
	//		StartDate:   startDate,
	//		EndDate:     endDate,
	//		SumDistance: 0,
	//		SumDuration: 0,
	//		SumColories: 0,
	//		Sumtimes:    0,
	//		AveragePace: 0,
	//		XData:    nums,
	//	}
	//	//ctx.JSON(lib.NewResponseFail(0, "没有运动数据"))
	//	ctx.JSON(lib.NewResponseOK(dataRecoed))
	//	return
	//}

	//空数据报错！
	//println("record length:", recordslens, PlanRecord[0].CreateAt.Unix())
	//每个周期结束时间戳数组，type=1,type=3
	xLimitArr := make([]int64, cycle+1)
	lenLimit := len(xLimitArr)
	lastTime := time.Now()
	zeroTime := time.Date(lastTime.Year(), lastTime.Month(), lastTime.Day(), 0, 0, 0, 0, lastTime.Location()) //最后一天零时时间戳
	if requestDurations.Type == 3 {
		zeroTime = time.Date(lastTime.Year(), lastTime.Month(), 1, 0, 0, 0, 0, lastTime.Location()) //最后一天零时时间戳
	}

	//空数据报错
	//println("最后一次运动的时间戳:", PlanRecord[recordslens-1].CreateAt.Unix(), "当天零时时间戳：", zeroTime.Unix(), "form:", zeroTime.Format("2006-01-02 15:04:05"))
	i2 := -1
	for c1 := lenLimit - 1; c1 >= 0; c1-- { //每个单元具体时间段,周和年
		if requestDurations.Type == 3 { //年
			xLimitArr[c1] = zeroTime.AddDate(0, -(i2), 0).Unix()
			println("每次的结束时间戳:", xLimitArr[c1], "form:", zeroTime.AddDate(0, -(i2), 0).Format("2006-01-02 15:04:05"), "i2:", -i2)
		} else if requestDurations.Type == 1 {
			xLimitArr[c1] = GetZeroTime(time.Now()).AddDate(0, 0, -(i2)).Unix()
			println("每次的结束时间戳:", xLimitArr[c1], "form:", zeroTime.AddDate(0, 0, -(i2)).Format("2006-01-02 15:04:05"), "i2:", -i2)
		}
		i2 = i2 + 1
	}

	//获取本周是当月的第几周和每周的开始和结束边界，type=2
	var weeks int
	xLimitArrWeeks := make([]int64, 0)
	if requestDurations.Type == 2 {
		//处理按周的边界，weekDay()的周日是0！！！
		//新版每月，周处理
		nowTime := time.Now()
		firstWeekDay := time.Date(nowTime.Year(), nowTime.Month(), 1, 0, 0, 0, 0, nowTime.Location()) //本月第一天零时
		//firstWeekDayUnix:=firstWeekDay.Unix()//本月第一天零时时间戳
		firstWeekDayInt := int(firstWeekDay.Weekday()) //本月第一天是周几
		if firstWeekDayInt == 0 {
			firstWeekDayInt = 7
		}
		secondWeekDay := firstWeekDay.AddDate(0, 0, 8-int(firstWeekDay.Weekday())) //本月第二周第一天零时
		secondWeekDayUnix := secondWeekDay.Unix()
		//本月最后一天24时
		lastWeekDay := time.Date(nowTime.Year(), nowTime.Month()+1, 0, 24, 0, 0, 0, nowTime.Location())
		lastWeekDayUnix := lastWeekDay.Unix() //本月最后一天24时时间戳
		lastWeekDayInt := int(lastWeekDay.Weekday())
		if lastWeekDayInt == 0 {
			lastWeekDayInt = 7
		}
		xLimitArrWeeks = append(xLimitArrWeeks, firstWeekDay.Unix())  //第一天零时
		xLimitArrWeeks = append(xLimitArrWeeks, secondWeekDay.Unix()) //第二周第一天零时
		println("加入第一和第二个时间")
		for index, value := range xLimitArrWeeks {
			println("index:", index, "value:", value, "form date:", time.Unix(value, 0).Format("2006-01-02 15:04:05"))
		}
		println("开始加入中间时间")
		//获取中间每周边界
		for indexUnix := secondWeekDayUnix + 604800; indexUnix < lastWeekDayUnix; indexUnix = indexUnix + 604800 {
			xLimitArrWeeks = append(xLimitArrWeeks, indexUnix)
		}
		//最后一天24时
		xLimitArrWeeks = append(xLimitArrWeeks, lastWeekDayUnix) //最后一天24时
		//获取本周时本月第几周
		weekNumber := 1
		for index, value := range xLimitArrWeeks {
			println("index:", index, "value:", value, "form date:", time.Unix(value, 0).Format("2006-01-02 15:04:05"))
			if nowTime.Unix() > value {
				weekNumber++
				continue
			}
			if nowTime.Unix() > xLimitArrWeeks[index-1] {
				weeks = index
				break
			}
		}
		println("本月时第几周：", weeks)

	}
	println("xLimitArrWeeks:")

	for index, value := range xLimitArrWeeks {
		valUnix := time.Unix(value, 0)
		println("index:", index, "value:", valUnix.Format("2006-01-02 15:04:05"))
	}

	//获取柱状图数据，x轴数组-----开始
	k := 0
	cycle1 := cycle
	xLimitArrs := xLimitArr
	if requestDurations.Type == 2 {
		cycle1 = weeks
		xLimitArrs = xLimitArrWeeks
	}
	sumDays := make([]int, cycle1)
	singleDistance := make([]float32, cycle1) //数据
	println("cycle1:", cycle1)
	for i := 0; i < cycle1; i++ { //第一层循环，表示数组的长度
		println("进入第", i, "次循环")
		for in := k; in < recordslens; in++ { //第二层循环，表示数据。根据每个记录的创建时间，在每个数组边界里，则数据加一，进行下一个，否则退出，进行下一个大的循环。
			println(PlanRecord[k].CreateAt.Format("2006-01-02 15:04:05"), time.Unix(xLimitArrs[i], 0).Format("2006-01-02 15:04:05"), time.Unix(xLimitArrs[i+1], 0).Format("2006-01-02 15:04:05"))
			if PlanRecord[k].CreateAt.Unix() >= xLimitArrs[i] && PlanRecord[k].CreateAt.Unix() <= xLimitArrs[i+1] {
				println("运动数据，plan_record_id:", PlanRecord[k].Id, "create_at:", PlanRecord[k].CreateAt.Format("2006-01-02 15:04:05"), "在limit ,i", i, time.Unix(xLimitArrs[i], 0).Format("2006-01-02 15:04:05"), "和i+1之间，", time.Unix(xLimitArrs[i+1], 0).Format("2006-01-02 15:04:05"))
				sumDays[i] = sumDays[i] + 1
				singleDistance[i] = singleDistance[i] + float32(PlanRecord[k].Distance) //获取柱状图数据
				k = k + 1
				continue
			}
			break
		}
	}

	//补全x轴数组----开始
	//根据周，月，年，获取x轴长度
	println("打印xLength:", xLength)
	fullLength := 0 //X轴长度
	if requestDurations.Type == 1 {
		fullLength = 7
	} else if requestDurations.Type == 2 {
		fullLength = xLength
	} else if requestDurations.Type == 3 {
		if time.Now().Month() <= 6 {
			fullLength = 6
		} else {
			fullLength = xLength
		}
	}

	for i := len(singleDistance); i < fullLength; i++ {
		singleDistance = append(singleDistance, 0)
	}
	//打印完全数据
	println("将精度精确到0.1")

	XData := make([]string, len(singleDistance))
	xAxis := make([]string, len(singleDistance)) //x轴单位

	//将数据单位换算成公里并保留小数点后一位
	for index, value := range singleDistance {
		println("index:", index, "value:", value)
		XData[index] = fmt.Sprintf("%0.2f", singleDistance[index]/1000)
		singleDistance[index] = float32(math.Trunc(float64(singleDistance[index]/1000)*1e2+0.5) * 1e-2)
	}
	//x轴单位处理----开始
	if requestDurations.Type == 1 || requestDurations.Type == 3 { //周和年

		for index := 0; index < xLength; index++ {
			xStr := ""
			xKey := index + 1
			switch xKey {
			case 1:
				xStr = "一"
			case 2:
				xStr = "二"
			case 3:
				xStr = "三"
			case 4:
				xStr = "四"
			case 5:
				xStr = "五"
			case 6:
				xStr = "六"
			case 7:
				xStr = "七"
			case 8:
				xStr = "八"
			case 9:
				xStr = "九"
			case 10:
				xStr = "十"
			case 11:
				xStr = "十一"
			case 12:
				xStr = "十二"
			}
			//年单独加本月
			if requestDurations.Type == 3 && index+1 == len(xAxis) {
				xStr = "本月"
			}
			xAxis[index] = xStr
		}
	} else if requestDurations.Type == 2 { //月

		for i := 0; i < len(xLimitArrWeeks)-1; i++ {
			xAxis[i] = time.Unix(xLimitArrWeeks[i], 0).Format("01/02")
		}
	}
	println("打印xAxis:")
	for index, value := range xAxis {
		println("index:", index, "value:", value)
	}
	//x轴单位处理----结束
	println(" singleDistanceStr")
	for index, value := range XData {
		println("index:", index, "value:", value)

	}
	//补全x轴数组----结束
	//获取柱状图数据，x轴数组-----结束

	//获取y轴数据----开始
	yAxis := make([]float64, 2) //y轴单位
	yAxis[0] = 0
	//获取单次最低公里数
	student := models.Student{}
	res, err := lib.Engine.Table("student").Where("id=?", id).Get(&student)
	if err != nil {
		fmt.Printf("获取学生信息出错：%v", err.Error())
		ctx.JSON(lib.NewResponseOK(err))
	}
	if res == false {
		println("该学生不存在")
		ctx.JSON(lib.NewResponseOK("该学生不存在"))
		return
	}

	plan := models.Plan{}
	res, err = lib.Engine.Table("plan").Where("id=?", requestDurations.PlanId).Get(&plan)
	if err != nil {
		fmt.Printf("获取运动计划出错：%v", err.Error())
		ctx.JSON(lib.NewResponseOK(err))
	}
	if res == false {
		println("该计划不存在")
		ctx.JSON(lib.NewResponseOK("该计划不存在"))
		return
	}
	//minDistance:=0
	//if plan.Types == 1 { //公里模式
	//	if student.Gender == 1 { //男生
	//		minDistance = plan.BoySingleMindistance //单次最低
	//	} else {
	//		minDistance = plan.GirlSingleMindistance //单次最低
	//	}
	//}
	//
	//if plan.Types == 2 { //以次数模式
	//	if student.Gender == 1 { //男生
	//		minDistance = plan.BoySingleMindistance //单次最低公里数
	//	} else {
	//		minDistance = plan.GirlSingleMindistance //单次最低公里数
	//	}
	//}
	////yAxis[1]=fmt.Sprintf("%0.1fkm",float32(minDistance)/1000)//最低公里数
	//yAxis[1]=math.Trunc(float64(minDistance)/1000*1e2+0.5) * 1e-2

	maxDistance := singleDistance[0]
	println("")
	fmt.Printf("sing 0:%f,maxDistance:%f", singleDistance[0], maxDistance)
	println("")
	for index, value := range singleDistance {
		println("index:", index, "value:", value)
		if value > maxDistance {
			println("value:", value, ">", "max", maxDistance)
			maxDistance = value
		}
	}

	maxDis := lib.FloatRoundingToFloat(float64(maxDistance), 2)
	//yAxis[2]= fmt.Sprintf("%0.1fkm",maxDistance)
	//在最大的单次公里上加一公里
	if maxDis > 0 { //有数据，最大值就是数据中的最大值加1公里
		yAxis[1] = maxDis + 1
	} else { //没有数据,最大值：5
		yAxis[1] = 5
	}
	fmt.Printf("yAxis[1]:%f", yAxis[1])

	println("打印y轴数组")
	for index, value := range yAxis {
		println("yAxis,index:", index, "value:", value)
	}
	//获取y轴数据----结束

	//整理返回体数据
	var sumDistance, sumDuration, sumtimes int
	var sumColories float64
	for _, value := range PlanRecord {
		sumDistance = sumDistance + value.Distance
		sumDuration = sumDuration + value.Duration
		sumColories = sumColories + value.Calories
	}

	println("sumDuration:", sumDuration, "sumDistance:", sumDistance)
	//根据周月年，返回不同的x轴和y轴
	sumtimes = len(PlanRecord)

	//根据周月年获取不同统计

	averagePace := int((float64(sumDuration)) / (float64(float64(sumDistance) / 1000)))
	floPace := float64(sumDuration) / (float64(float64(sumDistance) / 1000))
	fmt.Printf("sumDuration:%f,sumDistance:%f,floPace:%f", sumDuration, sumDistance, floPace)
	strPace := lib.FomPace(floPace)
	if sumDistance <= 0 {
		strPace = "0'0''"
	}

	strColories := fmt.Sprintf("%.2f", float64(sumColories))
	strTimes := fmt.Sprintf("%d", sumtimes)
	strDistance := fmt.Sprintf("%.2f", float64(sumDistance)/1000)
	strDuration := lib.ConvertTime(sumDuration)

	fmt.Printf("flopace:%f,intflopace:%d,strPace:%v", floPace, averagePace, strPace)
	//averagePace:=int((float64(sumDuration)) / (float64(float64(sumDistance) / 1000)))
	//float64(finishRun.Duration) / float64(float64(finishRun.Distance)/1000)
	if averagePace <= 0 {
		averagePace = 0
	}

	//datas := newResponseDuration{
	//	StartDate:   startDate,
	//	EndDate:     endDate,
	//	SumDistance: sumDistance,
	//	SumDuration: sumDuration,
	//	SumColories: sumColories,
	//	Sumtimes:    sumtimes,
	//	AveragePace: averagePace,
	//	XData:    singleDistance,
	//	XAxis:xAxis,
	//	YAxis:yAxis,
	//	YUnit:"km",
	//	DurationType:durationType,
	//}

	resData := strResponseDuration{
		StartDate:    startDate,
		EndDate:      endDate,
		SumDistance:  strDistance,
		SumDuration:  strDuration,
		Sumtimes:     strTimes,
		SumColories:  strColories,
		AveragePace:  strPace,
		XData:        singleDistance,
		XAxis:        xAxis,
		YAxis:        yAxis,
		YUnit:        "km",
		DurationType: durationType,
	}

	//for index,value:=range singleDistance{
	//	println("index:",index,",value:",value)
	//}
	//println("")
	//for index,value:=range sumDays{
	//	println("index:",index,",value:",value)
	//}

	//返回记录
	ctx.JSON(lib.NewResponseOK(resData))

}

//分享
func share(ctx iris.Context) {

	var schoolId, studentId, planId, resUrl, ip string

	if ctx.URLParamExists("student_id") {
		studentId = ctx.URLParam("student_id")
	} else {
		ctx.JSON(lib.NewResponseFail(0, "获取不到student_id"))
		return
	}
	if ctx.URLParamExists("plan_id") {
		planId = ctx.URLParam("plan_id")
	} else {
		ctx.JSON(lib.NewResponseFail(0, "获取不到plan_id"))
		return
	}
	if ctx.URLParamExists("school_id") {
		schoolId = ctx.URLParam("school_id")
	} else {
		ctx.JSON(lib.NewResponseFail(0, "获取不到school_id"))
		return
	}

	cfg := configs.Conf.Web
	// todo 记得将ip地址放入配置文件,还有schoolId
	ip = "47.101.67.225"
	resUrl = ip + cfg.Addr + "/open/share/records?school_id=" + schoolId + "&student_id=" + studentId + "&plan_id" + planId

	ctx.JSON(lib.NewResponseOK(resUrl))
	return

}

//旧版一段时间内运动记录
func duration(ctx iris.Context) {
	//获取student_id
	id := ctx.Params().GetUint64Default("id", 0)
	print("student_id:", id)

	var requestDurations requestDuration
	if err := ctx.ReadJSON(&requestDurations); err != nil {
		//ctx.JSON(iris.Map{"type ReadJSON error": "错误的type类型"})
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		fmt.Print("%v", err)
		return
	}

	//验证
	valType := lib.ValidateRequest(requestDurations)
	if valType == false {
		ctx.JSON(lib.NewResponseFail(1, "时间类型格式错误"))
		return
	}

	//获取当前时间和一周前时间
	now := time.Now()
	beforeTime := ""
	var cycle int //周期
	var startDate string
	var endDate string

	switch requestDurations.Type { //获取x轴的长度，和开始的时间
	case 1:
		cycle = int(now.Weekday()) //周期天数
		startTime := now.AddDate(0, 0, -(cycle - 1))
		beforeTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location()).Format("2006-01-02")
		endTime1 := now.AddDate(0, 0, 7-cycle)
		endDate = time.Date(endTime1.Year(), endTime1.Month(), endTime1.Day(), 0, 0, 0, 0, startTime.Location()).Format("2006年1月2日")
		startDate = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location()).Format("2006年1月2日")
		println("")
		fmt.Printf("week,startTime:%v,endTime:%v", beforeTime, endDate)
	case 2:
		cycle = int(now.Day())
		startTime := now.AddDate(0, 0, -(cycle - 1))
		beforeTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location()).Format("2006-01-02")
		startDate = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location()).Format("2006年1月2日")
		endTime1 := now.AddDate(0, 0, 0)
		endDate = time.Date(endTime1.Year(), endTime1.Month()+1, 0, 0, 0, 0, 0, startTime.Location()).Format("2006年1月2日")
		println("")
		fmt.Printf("month,startTime:%v,endTime:%v", beforeTime, endDate)
	case 3:
		cycle = int(now.Month())
		startTime := now.AddDate(0, -(cycle - 1), 0)
		beforeTime = time.Date(startTime.Year(), startTime.Month(), 0, 0, 0, 0, 0, startTime.Location()).Format("2006-01-02 15:04:05")
		endTime1 := now.AddDate(0, 0, 0)
		endDate = time.Date(endTime1.Year(), endTime1.Month()+1, 0, 0, 0, 0, 0, startTime.Location()).Format("2006年")
		startDate = endDate
		println("")
		fmt.Printf("month,startTime:%v,endTime:%v", beforeTime, endDate)
	default:
		ctx.JSON(lib.NewResponseFail(1, "时间类型格式错误"))
		return
	}
	println("今天是第", int(now.Month()), "月，周：", int(now.Weekday()), "日：", int(now.Day()))

	//本周到现在的每天的运动里程，本月到现在，每天的运动里程，本年到现在，每月的运动里程
	PlanRecord := []models.PlanRecord{}
	//根据id查询
	err := lib.Engine.Table("plan_record").Where("student_id=?", id).And(builder.Between{"create_at", beforeTime, now.Format("2006-01-02 15:04:05")}).Asc("create_at").Find(&PlanRecord)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	recordslens := len(PlanRecord)
	if recordslens == 0 {
		//没有运动数据返回空数组
		nums := make([]int, 0)
		dataRecoed := responseDuration{
			StartDate:   startDate,
			EndDate:     endDate,
			SumDistance: 0,
			SumDuration: 0,
			SumColories: 0,
			Sumtimes:    0,
			AveragePace: 0,
			XNumbers:    nums,
		}
		//ctx.JSON(lib.NewResponseFail(0, "没有运动数据"))
		ctx.JSON(lib.NewResponseOK(dataRecoed))
		return
	}
	println("record length:", recordslens, PlanRecord[0].CreateAt.Unix())
	//每个周期结束时间戳数组，type=1,type=3
	xLimitArr := make([]int64, cycle+1)
	lenLimit := len(xLimitArr)
	lastTime := PlanRecord[recordslens-1].CreateAt
	zeroTime := time.Date(lastTime.Year(), lastTime.Month(), lastTime.Day(), 0, 0, 0, 0, lastTime.Location()) //最后一天零时时间戳
	if requestDurations.Type == 3 {
		zeroTime = time.Date(lastTime.Year(), lastTime.Month(), 1, 0, 0, 0, 0, lastTime.Location()) //最后一天零时时间戳
	}

	println("最后一次运动的时间戳:", PlanRecord[recordslens-1].CreateAt.Unix(), "当天零时时间戳：", zeroTime.Unix(), "form:", zeroTime.Format("2006-01-02 15:04:05"))
	i2 := -1
	for c1 := lenLimit - 1; c1 >= 0; c1-- { //每个单元具体时间段,周和年
		if requestDurations.Type == 3 {
			xLimitArr[c1] = zeroTime.AddDate(0, -(i2), 0).Unix()
			println("每次的结束时间戳:", xLimitArr[c1], "form:", zeroTime.AddDate(0, -(i2), 0).Format("2006-01-02 15:04:05"), "i2:", -i2)
		} else {
			xLimitArr[c1] = zeroTime.AddDate(0, 0, -(i2)).Unix()
			println("每次的结束时间戳:", xLimitArr[c1], "form:", zeroTime.AddDate(0, 0, -(i2)).Format("2006-01-02 15:04:05"), "i2:", -i2)
		}
		i2 = i2 + 1
	}

	//获取本周是当月的第几周和每周的开始和结束边界，type=2
	var weeks int
	xLimitArrWeeks := make([]int64, 0)
	if requestDurations.Type == 2 {
		weeksNumbers1 := cycle / 7
		weeksNumbers2 := cycle % 7
		weeks = weeksNumbers1
		if weeksNumbers1 == 0 {
			weeks = 1
		}
		if weeksNumbers2 > 0 && weeksNumbers1 != 0 {
			weeks = weeksNumbers1 + 1
		}
		println("当前type是", requestDurations.Type, ",是第", weeks, "周")
		for i := 0; i < weeks*2; i++ {
			xLimitArrWeeks = append(xLimitArrWeeks, zeroTime.AddDate(0, 0, i*7).Unix())
		}

	}

	//获取柱状图数据，x轴数组-----开始
	k := 0
	cycle1 := cycle
	xLimitArrs := xLimitArr
	if requestDurations.Type == 2 {
		cycle1 = weeks
		xLimitArrs = xLimitArrWeeks
	}
	sumDays := make([]int, cycle1)
	println("cycle1:", cycle1)
	for i := 0; i < cycle1; i++ { //第一层循环，表示数组的长度
		for in := k; in < recordslens; in++ { //第二层循环，表示数据。根据每个记录的创建时间，在每个数组边界里，则数据加一，进行下一个，否则退出，进行下一个大的循环。
			println(PlanRecord[k].CreateAt.Format("2006-01-02 15:04:05"), time.Unix(xLimitArrs[i], 0).Format("2006-01-02 15:04:05"), time.Unix(xLimitArrs[i+1], 0).Format("2006-01-02 15:04:05"))
			if PlanRecord[k].CreateAt.Unix() >= xLimitArrs[i] && PlanRecord[k].CreateAt.Unix() <= xLimitArrs[i+1] {
				sumDays[i] = sumDays[i] + 1
				k = k + 1
				continue
			}
			break
		}
	}
	//获取柱状图数据，x轴数组-----结束
	//整理返回体数据
	var sumDistance, sumDuration, sumtimes int
	var sumColories float64
	for _, value := range PlanRecord {
		sumDistance = sumDistance + value.Distance
		sumDuration = sumDuration + value.Duration
		sumColories = sumColories + value.Calories
	}
	sumtimes = len(PlanRecord)
	datas := responseDuration{
		StartDate:   startDate,
		EndDate:     endDate,
		SumDistance: sumDistance,
		SumDuration: sumDuration,
		SumColories: sumColories,
		Sumtimes:    sumtimes,
		AveragePace: int((float64(sumDuration)) / (float64(sumDistance % 1000))),
		XNumbers:    sumDays,
	}

	//返回记录
	ctx.JSON(lib.NewResponseOK(datas))

}

//获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

//获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

//新的获取学生所有运动记录，防止和api冲突
func newSearch(ctx iris.Context) {
	//创建查询Session
	query := lib.Engine.Table("plan_record")

	//字段查询
	if ctx.URLParamExists("student_id") {
		//query.And(builder.Like{"student_id", ctx.URLParam("student_id")})
		query.And("student_id=?", ctx.URLParam("student_id"))
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
	//page := ctx.URLParamIntDefault("page", 0)
	//size := ctx.URLParamIntDefault("size", 50)
	//query.Limit(size, page*size)

	//获取计划名称
	query.Join("INNER", "plan", "plan.id=plan_record.plan_id")
	//查询
	var planRecord []responseRecord
	err := query.Find(&planRecord)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//处理步频数组，在运动详细记录里展示步频折线图
	for index, _ := range planRecord {
		xFrequencies := make([]int, 0)
		x_times := make([]int, 0)
		if len(planRecord[index].Frequencies) <= 5 {
			//for i1, _ := range planRecord[index].Frequencies {
			//	xFrequencies = append(xFrequencies, planRecord[index].Frequencies[i1])
			//}
			for i := 0; i < 5; i++ {
				if len(planRecord[index].Frequencies) > i {
					xFrequencies = append(xFrequencies, planRecord[index].Frequencies[i])
				} else {
					xFrequencies = append(xFrequencies, 0)
				}
				x_times = append(x_times, i)
			}
			planRecord[index].XFrequencies = xFrequencies
			planRecord[index].XNumber = x_times
			println("")
			//fmt.Printf("展示的xFrequenceies:%v", xFrequencies)
		} else {
			//获取步频数组长度
			freLength := len(planRecord[index].Frequencies)
			xIndex := 1
			for i2, _ := range planRecord[index].Frequencies {
				println("i2:", i2, "(freLength/4)*xIndex):", (freLength/4)*xIndex, "xIndex:", xIndex)
				if i2 == 0 {
					//planRecord[index].XNumber[xIndex] = i2
					x_times = append(x_times, i2)
					xFrequencies = append(xFrequencies, planRecord[index].Frequencies[i2])
				} else if (i2 == ((freLength / 4) * xIndex)) && xIndex <= 4 {
					//planRecord[index].XNumber[xIndex] = i2
					x_times = append(x_times, i2)                                          //x轴
					xFrequencies = append(xFrequencies, planRecord[index].Frequencies[i2]) //y轴
					xIndex++
				}
			}
			planRecord[index].XFrequencies = xFrequencies
			planRecord[index].XNumber = x_times
			println("")
			//fmt.Printf("展示的xFrequenceies:%v", xFrequencies)
		}

	}

	//步频折线图第二版-----使用每段时间内的步数除分钟数
	//思路：找到每段结束，用循环往前加i个

	//fmt.Printf("展示的xFrequenceies:%v", xFrequencies)
	//println("len:", len(planRecord))

	//新的返回

	//获取学生性别
	student := models.Student{}
	b, err := lib.Engine.Table("student").Where("id=?", ctx.URLParam("student_id")).Get(&student)
	if err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, "查询学生错误"))
		return
	}

	if b == false {
		println("查询学生失败")
		ctx.JSON(lib.NewResponseFail(1, "查询学生失败"))
		return
	}

	for i, _ := range planRecord {

		//获取边界值---开始
		plan := models.Plan{}
		b, err := lib.Engine.Table("plan").Where("id=?", planRecord[i].PlanId).Get(&plan)
		if err != nil {
			fmt.Printf("%v", err)
			ctx.JSON(lib.NewResponseFail(1, "查询计划错误"))
			return
		}
		if b == false {
			println("查询计划失败")
			ctx.JSON(lib.NewResponseFail(1, "查询计划失败"))
			return
		}

		conf := configs.Conf
		minPace := plan.MinPace //最慢配速
		maxPace := plan.MaxPace //最快配速
		//正常最低步频120
		minFrequency := conf.Limit.MinFrequency
		maxFrequency := conf.Limit.MaxFrequency

		//计划有限制使用计划的限制

		//获取单次最低公里数---开始
		var minDistance int
		minDistance = plan.MinSingleDistance

		//获取单次最低公里数---结束
		//获取边界值---结束

		var PlanRecordString []string
		PlanRecordString = make([]string, 0)
		for _, value := range planRecord[i].InvalidCode {
			if value == 1 {
				PlanRecordString = append(PlanRecordString, "配速过低")
				planRecord[i].PaceObject.Status = true
				planRecord[i].PaceObject.Reason = lib.FomPace(float64(minPace))
				//b.WriteString( "配速过低,")
				//planRecord[i].PlanRecordString = "配速过低"

			}
			if value == 2 {
				PlanRecordString = append(PlanRecordString, "配速过高")
				planRecord[i].PaceObject.Status = true
				planRecord[i].PaceObject.Reason = lib.FomPace(float64(maxPace))
				//b.WriteString( "配速过高,")
				//planRecord[i].PlanRecordString = "配速过高"

			}
			if value == 3 {
				PlanRecordString = append(PlanRecordString, "步频过低")
				planRecord[i].FrequencyObject.Status = true
				planRecord[i].FrequencyObject.Reason = strconv.Itoa(minFrequency)

				//b.WriteString( "步频过低,")
				//planRecord[i].PlanRecordString = "步频过低"

			}
			if value == 4 {
				PlanRecordString = append(PlanRecordString, "步频过高")
				planRecord[i].FrequencyObject.Status = true
				planRecord[i].FrequencyObject.Reason = "步频过高"
				planRecord[i].FrequencyObject.Reason = strconv.Itoa(maxFrequency)
				//b.WriteString( "步频过高,")
				//planRecord[i].PlanRecordString = "步频过高"

			}
			if value == 5 {
				PlanRecordString = append(PlanRecordString, "超过24小时")
				planRecord[i].DurationObject.Status = true
				planRecord[i].DurationObject.Reason = lib.ConvertTime(3600 * 24)
				//b.WriteString("超过24小时,")
				//planRecord[i].PlanRecordString = "超过24小时"

			}
			if value == 6 {
				PlanRecordString = append(PlanRecordString, "没有经过所有打卡点")
				//b.WriteString("没有经过所有打卡点,")
				//planRecord[i].PlanRecordString = "没有经过所有打卡点"

			}
			if value == 7 {
				PlanRecordString = append(PlanRecordString, "未达到最低公里数")
				planRecord[i].DistanceObject.Status = true
				planRecord[i].DistanceObject.Reason = strconv.FormatFloat(lib.FloatRoundingToFloat(float64(minDistance)/1000, 2), 'f', 2, 64)
				//b.WriteString("未达到最低公里数,")
				//planRecord[i].PlanRecordString = "未达到最低公里数"

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
			planRecord[i].PlanRecordString = PlanRecordString
		}

		planRecord[i].DistanceObject.RawData = lib.FloatRoundingToFloat(float64(planRecord[i].Distance)/1000, 2)
		planRecord[i].DurationObject.RawData = lib.ConvertTime(planRecord[i].Duration)
		planRecord[i].PaceObject.RawData = planRecord[i].FormPace
		planRecord[i].StepPageObject.RawData = planRecord[i].Steps
		planRecord[i].FrequencyObject.RawData = int(planRecord[i].Frequency)
	}

	ctx.JSON(lib.NewResponseOK(planRecord))
}

func passPoints(ctx iris.Context) {
	//获取运动记录id
	recordId := 0
	if ctx.URLParamExists("record_id") {
		//recordId=int(ctx.Params().GetUint64Default("record_id", 0))
		recordId, _ = ctx.URLParamInt("record_id")

		fmt.Printf("%v", recordId)
	} else {
		ctx.JSON(lib.NewResponseFail(1, "无运动记录id"))
		return
	}

	//获取运动记录
	record := models.PlanRecord{}
	res, err := lib.Engine.Table("plan_record").ID(recordId).Get(&record)
	if err != nil {
		fmt.Printf("%v", err.Error())
		ctx.JSON(lib.NewResponseFail(1, "查询运动记录失败"))
		return
	}
	if res == false {
		ctx.JSON(lib.NewResponseFail(1, "运动记录不存在"))
		return
	}
	println("")
	fmt.Printf("%v", record)
	println("record,route id:", record.RouteId)
	//获取运动记录路径
	route := models.PlanRoute{}
	resRoute, err := lib.Engine.Table("plan_route").Where("id=?", record.RouteId).Get(&route)
	if err != nil {
		fmt.Printf("查询运动路径出错：%v", err.Error())
		ctx.JSON(lib.NewResponseFail(1, "查询运动路径失败"))
		return
	}
	if resRoute == false {
		ctx.JSON(lib.NewResponseFail(1, "运动路径不存在"))
		return
	}

	sqlString := ""
	for _, value := range route.Route {
		sqlString = sqlString + " or id=" + strconv.FormatInt(value, 10)
	}
	orString := sqlString[3:len(sqlString)]
	println(orString)

	routePoints := make([]point, len(route.Route))
	for index, value := range route.Route {
		println("index:", index, "id:", value)
		routePoints[index].Id = int(value)
	}

	points := []point{}
	errPoints := lib.Engine.Table("plan_points").Where(orString).Find(&points)
	if errPoints != nil {
		fmt.Printf("查询点位出错：%v", errPoints.Error())
		ctx.JSON(lib.NewResponseFail(1, "点位错误"))
		return
	}
	fmt.Printf("resss%v", points)

	//补充路径点的详细信息
	for index, value := range routePoints {
		for _, value2 := range points {
			if value.Id == value2.Id {
				routePoints[index].Longitude = value2.Longitude
				routePoints[index].Latitude = value2.Latitude
			}
		}
	}

	//筛选蓝牙点
	//for index, value := range points {
	//	for _, value2 := range record.PassPoints {
	//		if value.Id == value2 {
	//			points[index].Status = 1
	//		}
	//	}
	//}

	for index, value := range routePoints {
		for _, value2 := range record.PassPoints {
			if value.Id == value2 {
				routePoints[index].Status = true
			}
		}
	}

	responPoint := responseDetail{Points: routePoints}
	//返回点详细信息
	ctx.JSON(lib.NewResponseOK(responPoint))

}
