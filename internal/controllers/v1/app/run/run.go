package run

import (
	"Campus/configs"
	"Campus/internal/lib"
	"Campus/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type requestStartRun struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	PlanId    int     `json:"plan_id"`
}

type responseStartRun struct {
	Fence       models.PlanFence `json:"fence"`
	Route       []route          `json:"route"`
	PlanStartId int              `json:"plan_start_id"`
	Exist       bool             `json:"exist"`
}

type route struct {
	Points  []*models.PlanPoints `json:"points"`
	IBeacon []*models.PlanPoints `json:"iBeaconPoints"`
}

type existFences struct {
	Exist bool `json:"exist"`
}

type passPoints struct {
	Id     int `json:"id" xorm:"id"`
	Status int `json:"status" xorm:"status"`
}

type iBeaconsPoints struct {
	Id         int       `json:"id" xorm:"autoincr pk id"`
	Name       string    `json:"name" xorm:"name"`
	FenceId    int       `json:"fence_id" xorm:"fence_id"`
	Longitude  float32   `json:"longitude" xorm:"longitude DOUBLE"`
	Latitude   float32   `json:"latitude" xorm:"latitude DOUBLE"`
	Type       int       `json:"type" xorm:"type"`
	Uuid       string    `json:"uuid" xorm:"uuid"`
	CreateAT   time.Time `json:"create_at" xorm:"created 'create_at'"`
	FaceStatus bool      `json:"face_status" xorm:"-"`
}

//开始跑步的返回体的String,用于格式化打印返回体
func (startRunRes *responseStartRun) String() string {
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

//开始跑步接口
func startRun(ctx iris.Context) {
	//读取数据
	startRun := requestStartRun{}
	err := ctx.ReadJSON(&startRun)
	if err != nil {

		println(err)
		ctx.JSON(lib.NewResponseFail(1, "读取startRun参数错误"))
		return
	}
	//简单验证
	//验证
	valType := lib.ValidateRequest(startRun)
	if valType == false {
		ctx.JSON(lib.NewResponseFail(1, "时间类型格式错误"))
		return
	}
	fmt.Printf("获取的startRun：%v,lon:%v,lat:%v", startRun, startRun.Longitude, startRun.Latitude)
	//TODO 没有计划返回

	//1.根据经纬度获取最近的围栏
	fenceInfo := make([]models.PlanFence, 0)
	err = lib.Engine.Table("plan_fence").Find(&fenceInfo)
	if err != nil {
		println("查询围栏信息出错")
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if len(fenceInfo) == 0 {
		ctx.JSON(lib.NewResponse(0, "没有电子围栏", nil))
		return
	}
	//遍历围栏的中心点，查询最近的围栏，len(fenceInfo)
	distance := make([]float64, len(fenceInfo))
	cloestFence := fenceInfo[0]
	//所有围栏中心点与定位点的距离
	for i := 0; i < len(fenceInfo); i++ {
		distance[i] = lib.GetDistance(startRun.Latitude, fenceInfo[i].Latitude, startRun.Longitude, fenceInfo[i].Longitude)
		if err != nil {
			println("查询定位点与围栏中心点距离出错")
			ctx.JSON(lib.NewResponseFail(1, err.Error()))
		}
		//println("循环得到的距离：", distance[i])
		fmt.Printf("i%v,围栏的距离%v：,围栏：%v", i, distance[i], fenceInfo[i])
	}
	//寻找最近的围栏并赋值给closestFence
	min := distance[0]
	for i2 := 0; i2 < len(fenceInfo); i2++ {
		if distance[i2] < min {
			min = distance[i2]
			cloestFence = fenceInfo[i2]
		}
	}

	fmt.Printf("最近的围栏%v：,距离%v：", cloestFence, min)
	// 如果最近距离的围栏还是太远比如3KM直接返回，距离太远
	if min > 3 {
		println("附近没有围栏，返回")
		existFences := existFences{Exist: false}
		ctx.JSON(lib.NewResponse(0, "当前定位距离最近围栏超过3公里", existFences))
		fmt.Printf("围栏距离太远返回的数据：%v", existFences)
		return
	}
	println("fence_id:", cloestFence.Id)

	//TODO 获取学生计划的单次最低距离
	token := ctx.Values().Get("jwt").(*jwt.Token)
	studentId := token.Claims.(jwt.MapClaims)["id"].(float64)
	println("学生：", studentId)
	fmt.Printf("学生：%d", studentId)
	//性别
	studend := models.Student{}
	bl1, err := lib.Engine.Table("student").Where("id=?", studentId).Get(&studend)
	if err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(0, err.Error()))
		return
	}
	if bl1 == false {
		println("没有找到学生，无法获取最低跑步距离")
		ctx.JSON(lib.NewResponseFail(0, "没有找到学生，无法获取最低跑步距离"))
		return
	}
	//计划里的单次最低距离
	plan := models.Plan{}
	bl2, err := lib.Engine.Table("plan").Where("id=?", startRun.PlanId).Get(&plan)
	if err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(0, err.Error()))
		return
	}
	if bl2 == false {
		println("没有找到计划，无法获取最低跑步距离")
		ctx.JSON(lib.NewResponseFail(0, "没有找到计划，无法获取最低跑步距离"))
		return
	}

	//跑步时间段判断 start
	//日期判断
	timeLayout := "2006-01-02 15:04:05"
	nowUnix1 := time.Now().Unix()
	startDateUnix := plan.DateBegin.Unix()
	planEndDate, _ := time.ParseInLocation(timeLayout, plan.DateEnd.Format(timeLayout)[0:10]+" 23:59:59", time.Local)
	endDateUnix := planEndDate.Unix()
	if nowUnix1 > endDateUnix || nowUnix1 < startDateUnix {
		println("不在计划跑规定时间，已返回")
		ctx.JSON(lib.NewResponseFail(0, "不在计划跑规定时间，已返回"))
		return
	}

	//时间段判断
	var timeFrame []models.PlanTimeFrame
	err = lib.Engine.Table("plan_time_frame").Where("plan_id=?", plan.Id).Find(&timeFrame)
	if err != nil {
		fmt.Printf("%v", err.Error())
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	now := time.Now().Format(timeLayout)
	nowTime, _ := time.ParseInLocation(timeLayout, now, time.Local)
	nowUnix := nowTime.Unix()
	isOnTime := false
	if len(timeFrame) > 0 {
		for _, value := range timeFrame {
			planDurationStart := now[0:11] + value.DurationBegin + ":00"
			planDurationEnd := now[0:11] + value.DurationEnd + ":00"
			planStartTime, _ := time.ParseInLocation(timeLayout, planDurationStart, time.Local)
			planEndTime, _ := time.ParseInLocation(timeLayout, planDurationEnd, time.Local)
			planStartUnix := planStartTime.Unix()
			planEndUnix := planEndTime.Unix()
			println(" ")
			fmt.Printf("今天的计划允许开始时间：%v,今天计划允许的结束时间：%v,now:%v,现在的时间戳：%v,今天计划允许开始的时间戳:%v,今天计划结束的时间戳：%v", planDurationStart, planDurationEnd, now, nowUnix, planStartUnix, planEndUnix)
			if nowUnix > planStartUnix && nowUnix < planEndUnix {
				isOnTime = true
				break
			}

		}
	}
	if !isOnTime {
		ctx.JSON(lib.NewResponseFail(0, "不在计划跑步的允许时间范围，已返回"))
		return
	}
	//跑步时间段判断 end

	//获取单次最低公里数
	var minDistance float64
	//新，按公里和按次数的最低公里数一样了。。。
	minDistance = float64(plan.MinSingleDistance) / 1000

	println("本次动态路径：")
	fmt.Printf("性别：%v,目标公里数minDistance:%v", studend.Gender, minDistance)
	//动态路径开始
	dyroutes, err := DynamicRoute(minDistance, startRun.Longitude, startRun.Latitude)
	if err != nil {
		fmt.Printf("动态路径设置错误：%v", err)
		ctx.JSON(lib.NewResponseFail(0, err.Error()))
		return
	}

	//去掉最后一个点
	lessDyroutes := dyroutes[0 : len(dyroutes)-1]

	println("动态路径第二版,点的数量：", len(dyroutes))
	//动态路径结束

	//合并plan_start和plan_route-------------------------开始
	//合并的表名字plan_route,字段主要为plan_start的字段，将route_id字段更换为route,即路径点的数组。
	//1.获取动态路径点的id数组，route
	pointsId := make([]int64, 0)
	for _, value := range dyroutes {
		pointsId = append(pointsId, int64(value.Id))
	}
	//2.插入plan_route
	//从token中获取userId

	//没有计划的是自由跑
	//bl,err:=lib.Engine.Table("plan").Where("id=?",startRun.PlanId).Exist()
	//if bl!=true{
	//	println("没有计划，返回")
	//	ctx.JSON(lib.NewResponseFail(1, "运动计划不存在"))
	//	return
	//}
	//修复提交冲突
	newStart := models.PlanRoute{}
	newStart.FenceId = cloestFence.Id
	newStart.PlanId = startRun.PlanId
	newStart.Route = pointsId
	newStart.StudentId = int(studentId)
	_, err2 := lib.Engine.Table("plan_route").Insert(&newStart)
	if err2 != nil {
		println("插入plan_route表失败")
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	//合并plan_start和plan_route-------------------------结束

	//获取蓝牙点
	//iBeaconRoutes:= make([]route, 0)
	iBeaconPoints := make([]*models.PlanPoints, 0)
	for _, value := range dyroutes {

		if len(value.Name) > 0 {
			iBeaconPoints = append(iBeaconPoints, value)
		}
	}

	//TODO 需要添加一条默认路径
	p1 := make([]route, 0)
	responseRoutes := route{lessDyroutes, iBeaconPoints}
	p1 = append(p1, responseRoutes)
	//返回
	responseStartRun := responseStartRun{
		Fence:       cloestFence,
		Route:       p1,
		PlanStartId: newStart.Id,
		Exist:       true,
	}
	//fmt.Printf("返回体设置完成%v:",responseStartRun)

	//println(cloestFence.Points)
	//respose,err:=json.Marshal(responseStartRun)
	//if err!=nil{
	//	println(err)
	//
	//}
	println("打印返回体")
	fmt.Printf("%v", responseStartRun)
	nums, err := ctx.JSON(lib.NewResponse(0, "返回最近围栏信息和围栏内路径成功", responseStartRun))
	println("ctx JSON num", nums)
	if err != nil {
		fmt.Printf("%v", err)
	}

	//格式化打印返回结构体，用于查询安卓连接非设置线段问题
	fmt.Println("StartRun:", responseStartRun.String())
	//ctx.JSON(lib.NewResponse(0, "返回最近围栏信息和围栏内路径成功", responseStartRun))
}

//跑步结束接口
func finishRun(ctx iris.Context) {

	//获取全部数据
	finishRun := models.RequestFinishRun{}
	err := ctx.ReadJSON(&finishRun)
	if err != nil {
		ctx.JSON(lib.NewResponse(1, "参数错误", err))
		return
	}

	//验证
	//println("startid:",finishRun.PlanStartId)
	fmt.Printf("步频数组frequencies:%v", finishRun.Paces)
	frequencies := make([]int, 0)
	for _, value := range finishRun.Paces {
		frequencies = append(frequencies, value)
	}

	//计算步频、配速，速度-----开始

	var finishPace, finishFrequency, finishSpeed float64
	if finishRun.Distance == 0 || finishRun.Duration == 0 { //当距离为0
		finishRun.Pace = 0
		finishFrequency = 0
		finishSpeed = 0
	} else {
		//平均配速（计算公式：时间/公里，单位为秒和公里）
		finishPace = float64(finishRun.Duration) / float64(float64(finishRun.Distance)/1000)
		finishRun.Pace = int(finishPace)

		//平均步频，保留两位小数
		finishFrequency = float64(finishRun.Steps) / float64(float64(finishRun.Duration)/60)
		n10 := math.Pow10(2)
		finishFrequency = math.Trunc((finishFrequency+0.5/n10)*n10) / n10

		//平均速度，保留两位小数
		finishSpeed = float64(float64(finishRun.Distance)/1000) / float64(float64(finishRun.Duration)/3600)
		finishSpeed = math.Trunc((finishSpeed+0.5/n10)*n10) / n10
	}
	println("")
	fmt.Printf("get pace: %v", finishRun.Pace)

	//计算步频、配速，速度-----结束

	//查询开始时间 Where("id=?", finishRun.PlanStartId)
	//routeId:=finishRun.PlanStartId
	start_planTime := models.PlanRoute{}
	b, err := lib.Engine.Table("plan_route").Where("id=?", finishRun.PlanStartId).Get(&start_planTime)
	if err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, "查询开始时间错误"))
		return
	}
	fmt.Printf("%v", start_planTime)
	if b == false {
		println("查询开始时间失败")
		ctx.JSON(lib.NewResponseFail(1, "查询开始时间失败"))
		return
	}

	//补充保存打卡点---开始

	sqlString := ""
	for _, value := range start_planTime.Route {
		sqlString = sqlString + " or id=" + strconv.FormatInt(value, 10)
	}
	orString := sqlString[3:len(sqlString)]

	//1.获取围栏内点的详细信息。记得以后给所有点加软删
	routePoints := make([]iBeaconsPoints, 0)
	alliBeaconsErr := lib.Engine.Table("plan_points").Where(orString).Find(&routePoints)
	if alliBeaconsErr != nil {
		ctx.JSON(lib.NewResponseFail(1, "查询打卡点失败"))
		return
	}

	for index, value := range finishRun.PassPoints {
		println("")
		fmt.Printf("经过点index:%v,value:%v", index, value)
	}

	//2.遍历经过点(蓝牙点)，二重遍历所有点，获取经过点的id。
	passiBeaconsId := make([]int, 0)
	for _, value := range finishRun.PassPoints {
		for _, value1 := range routePoints {
			if value.Longitude == value1.Longitude && value.Latitude == value1.Latitude {
				passiBeaconsId = append(passiBeaconsId, value1.Id)
				break
			}
		}
	}

	println("打印所有经过点的id")
	for _, value := range passiBeaconsId {
		println("id:", value)
	}

	//补充保存打卡点---结束

	//先用计划id查询计划表，获取计划类型，然后获取总计量单位。
	//然后查询计划进度表，获取当前进度
	//最后比较，当进度大于等于计划要求，更新进度状态为完成

	plan := models.Plan{}
	b, err = lib.Engine.Table("plan").Where("id=?", finishRun.PlanId).Get(&plan)
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

	planProgress1 := models.PlanProgress{}
	b, err = lib.Engine.Table("plan_progress").Where("plan_id=?", finishRun.PlanId).And("student_id=?", finishRun.StudentId).Get(&planProgress1)
	if err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, "查询计划进度错误"))
		return
	}
	if b == false {
		println("查询计划进度失败")
		ctx.JSON(lib.NewResponseFail(1, "查询计划进度失败"))
		return
	}
	//获取学生性别
	student := models.Student{}
	b, err = lib.Engine.Table("student").Where("id=?", finishRun.StudentId).Get(&student)
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

	//获取计划中的计量标准和计划中的进度
	//比较当前计划进度与计划
	var plantotal int
	var progresstotal int

	progresstotal = planProgress1.Times
	plantotal = plan.TotalTimes

	//更新进度表为完成状态---开始

	if progresstotal >= plantotal {

		ProgressStatus := models.PlanProgress{
			Status: 1,
		}

		res, err := lib.Engine.Table("plan_progress").Where("plan_id=?", finishRun.PlanId).And("student_id=?", finishRun.StudentId).Update(&ProgressStatus)
		if err != nil {
			fmt.Printf("%v", err)
			ctx.JSON(lib.NewResponseFail(1, err.Error()))
			return
		}
		if res == 1 {
			println("学生计划进度更新成功")
			//return
		}
		if res != 1 {
			println("计划进度更新失败")
		}
	}
	//更新进度表为完成状态----结束

	//查询该用户的院系以及性别

	//无效判断---开始
	//默认最低最高限制
	//正常最低minPace := 720

	//*************更改为男女****************
	conf := configs.Conf
	boyMinPace := conf.Limit.BoyMinPace
	girlMinPace := conf.Limit.GirlMinPace
	minPace := conf.Limit.MinPace //最慢配速
	maxPace := conf.Limit.MinPace //最快配速
	boyMaxPace := conf.Limit.BoyMaxPace
	girlMaxPace := conf.Limit.GirlMaxPace
	//正常最低步频120
	minFrequency := conf.Limit.MinFrequency
	maxFrequency := conf.Limit.MaxFrequency
	startPlantimeUnix := start_planTime.CreateAt.Unix() //跑步开始时间戳

	recordStatus := 0 //默认0，无效，1有效
	//计划有限制使用计划的限制

	if plan.MinPace != 0 {
		minPace = plan.MinPace
	}
	if plan.MaxPace != 0 {
		maxPace = plan.MaxPace
	}
	println("")
	fmt.Printf("pace:%v,frequency:%v", finishRun.Pace, finishFrequency)
	println("")
	fmt.Printf("minpace:%v,maxpace:%v,minFrequency:%v,maxFrequency:%v,boyMaxPace:%v,boyMinPace%v,girlMinPace%v,girlMaxPace%v",
		minPace, maxPace, minFrequency, maxFrequency, boyMaxPace, boyMinPace, girlMinPace, girlMaxPace)

	//获取单次最低公里数---开始
	var minDistance int
	minDistance = plan.MinSingleDistance

	println("目标最低公里数：", minDistance, "运动公里数：", finishRun.Distance)
	//获取单次最低公里数---结束

	//无效判断
	invalidCode := make([]int, 0)
	//判断女生的配速是否正常
	if student.Gender == 2 {

		if finishRun.Pace > girlMinPace { //配速太慢
			invalidCode = append(invalidCode, 1)
		}
		if finishRun.Pace < girlMaxPace { //配速太快
			invalidCode = append(invalidCode, 2)

		}

	}
	//判断男生的配速是否正常
	if student.Gender == 1 {

		if finishRun.Pace > boyMinPace { //配速太慢
			invalidCode = append(invalidCode, 1)
		}
		if finishRun.Pace < boyMaxPace { //配速太快
			invalidCode = append(invalidCode, 2)

		}
	}

	/*****************************************************/

	//老版本没有性别的配速判断

	//if finishRun.Pace > minPace { //配速太慢
	//	invalidCode = append(invalidCode, 1)
	//}
	//if finishRun.Pace < maxPace { //配速太快
	//	invalidCode = append(invalidCode, 2)
	//}

	/*****************************************************/

	if finishFrequency < float64(minFrequency) { //步频太低
		invalidCode = append(invalidCode, 3)
	}
	if finishFrequency > float64(maxFrequency) { //步频太高
		invalidCode = append(invalidCode, 4)
	}
	if (time.Now().Unix() - startPlantimeUnix) > 86400 { //超过24小时
		invalidCode = append(invalidCode, 5)
	}
	if finishRun.PointsStatus == false {
		invalidCode = append(invalidCode, 6) //没有经过所有打卡点
	}
	if minDistance > finishRun.Distance {
		invalidCode = append(invalidCode, 7) //没有完成最低公里数
	}
	if len(invalidCode) == 0 {
		recordStatus = 1
	}

	//无效判断---结束

	//格式化pace---开始
	formFrequency := lib.FomPace(float64(finishRun.Pace))

	//格式化pace--结束

	//超过24小时只插入运动记录
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	endTime, _ := time.ParseInLocation(timeLayout, finishRun.EndTime, loc)
	record := models.PlanRecord{
		PlanId:       finishRun.PlanId,
		StudentId:    finishRun.StudentId,
		StartTime:    start_planTime.CreateAt,
		EndTime:      endTime,
		Distance:     finishRun.Distance,
		Duration:     finishRun.Duration,
		Calories:     finishRun.Calories,
		Pace:         finishRun.Pace,
		FormPace:     formFrequency,
		Points:       finishRun.Points,
		Speed:        finishSpeed,
		Steps:        finishRun.Steps,
		Frequency:    finishFrequency,
		Frequencies:  frequencies,
		Status:       recordStatus,
		InvalidCode:  invalidCode,
		Gender:       student.Gender,
		DepartmentId: student.DepartmentId,
		RouteId:      finishRun.PlanStartId,
		PassPoints:   passiBeaconsId,
	}
	//计算运动时间间隔
	record.Times = record.EndTime.Sub(record.StartTime).Hours()

	//添加运动记录
	res, err := lib.Engine.Table("plan_record").Insert(record)
	if err != nil {
		//println("planrecord表插入失败")
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if res == 0 {
		ctx.JSON(lib.NewResponse(0, "运动记录上传失败", nil))
		return
	}

	//运动记录上传后，如果运动记录正常，更新计划进度----开始

	//获取单次最大公里数----开始
	// 如果有，且本次运动公里数超过最大公里数，只存入单次最大公里数
	singleMaxDistance := 0
	singleMaxDistance = plan.MaxSingleDistance
	//存在单次最大公里数且本次公里数超过单次最大公里数
	if singleMaxDistance != 0 && finishRun.Distance > singleMaxDistance {
		finishRun.Distance = singleMaxDistance
	}
	//获取单次最大公里数----结束
	//判断是否24小时内

	//24小时内，累加计划进度

	if len(invalidCode) == 0 {

		//小于计划要求，更新计划进度
		planProgress := models.PlanProgress{
			Distance: finishRun.Distance + planProgress1.Distance,
			Duration: finishRun.Duration + planProgress1.Duration,
			Calories: finishRun.Calories + planProgress1.Calories,
			Steps:    finishRun.Steps + planProgress1.Steps,
			Times:    planProgress1.Times + 1,
		}

		////修改多计划开始***********************************************************************************
		//record1 := models.PlanRecord{}
		//statusanddistance, err := lib.Engine.Table("plan_record").Where("status=1").And("student_id=?", finishRun.StudentId).
		//	And("YEARWEEK( DATE_FORMAT(  `plan_record`.`create_at`, '%Y-%m-%d' ),1 ) = YEARWEEK( NOW(),1 )").SumsInt(record1, "status", "distance")
		//if err != nil {
		//	fmt.Printf("查询周记录错误：%v", err)
		//	ctx.JSON(lib.NewResponseFail(1, err.Error()))
		//	return
		//}
		//println("\n\n\n\nstatusanddistance\n\n\n\n\n")
		//planProgress.WeekTimes = int(statusanddistance[0])
		//planProgress.WeekDistance = int(statusanddistance[1])
		//if plan.TotalDistance == 0 {
		//	ctx.JSON("\n\n\n注意！！！！女生的总公里数为0\n\n\n")
		//	return
		//
		//}
		//planProgress.CompleteProgress = float32(planProgress.Distance) / float32(plan.TotalDistance)
		//var completeweek = 1
		//weektotals := studentFunc.EveryMonthProgress(ctx, finishRun)
		//if len(weektotals) == 0 {
		//	ctx.JSON(lib.NewResponseFail(1, "获取该学生总周数失败 "))
		//	return
		//
		//}
		//
		//for _, weektotal := range weektotals {
		//	if weektotal.Status == 1 {
		//
		//		completeweek = completeweek + 1
		//	}
		//}
		//var buffer bytes.Buffer
		//completeweekstring := strconv.Itoa(completeweek)
		//
		//buffer.WriteString(completeweekstring)
		//buffer.WriteString("/")
		//buffer.WriteString(strconv.Itoa(len(weektotals)))
		//
		//planProgress.WeekCompleteProgress = buffer.String()
		//
		////修改多计划结束***********************************************************************************

		res, err := lib.Engine.Table("plan_progress").Where("plan_id=?", finishRun.PlanId).And("student_id=?", finishRun.StudentId).Update(&planProgress)
		if err != nil {
			fmt.Printf("更新计划进度错误：%v", err)
			ctx.JSON(lib.NewResponseFail(1, err.Error()))
			return
		}
		if res == 1 {
			println("学生计划进度更新成功")
		}
		if res == 0 {
			println("学生计划进度更新失败")
		}

	}
	//运动记录上传后，如果运动记录正常，更新计划进度----结束

	//更新连续运动天数和上次运动时间------开始

	now := time.Now()
	tm1 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := tm1.AddDate(0, 0, -1).Unix() //昨天零时时间戳
	today := tm1.Unix()                       //今天零时时间戳
	tomorrow := tm1.AddDate(0, 0, 1).Unix()   //明天零时时间戳
	lastUnix := student.LastSport.Unix()
	println("yesterday:", yesterday, "today:", today, "tomorrow:", tomorrow, "lastSport:", lastUnix)
	upStudent := models.Student{}
	if lastUnix < yesterday { //最后运动时间前天之前
		upStudent.Continue = 1
		println("前天之前")
	} else if yesterday < lastUnix && lastUnix < today { //最后运动时间昨天
		upStudent.Continue = student.Continue + 1
		println("昨天")
	} else if lastUnix > today && lastUnix < tomorrow { //最后运动时间今天
		upStudent.Continue = student.Continue
		println("今天")
	}
	if lastUnix > tomorrow {
		//大于今天报错
		ctx.JSON(lib.NewResponseFail(1, "上次运动时间错误"))
		return
	}
	upStudent.LastSport = now
	//更新
	affectds, err := lib.Engine.Table("student").Cols("continue", "last_sport").Where("id=?", finishRun.StudentId).Update(&upStudent)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, "更新学生连续运动天数和上次运动时间错误"))
		fmt.Printf("更新学生连续运动天数错误：%v", err)
		return
	}
	if affectds == 0 {
		ctx.JSON(lib.NewResponseFail(1, "更新学生连续运动天数和上次运动时间失败，无数据更新"))
		return
	}
	//更新连续运动天数和上次运动时间------结束

	ctx.JSON(lib.NewResponse(0, "运动记录上传成功", ""))
	return

}

func get(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	records := models.PlanRecord{}
	//根据id查询
	b, err := lib.Engine.Table("plan_record").ID(id).Get(&records)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该跑步记录"))
		return
	}
	ctx.JSON(lib.NewResponseOK(records))
}

//ctx.URLParam("name")
func search(ctx iris.Context) {
	//从token中获取userId
	token := ctx.Values().Get("jwt").(*jwt.Token)
	userId := token.Claims.(jwt.MapClaims)["id"].(float64)
	records := []models.PlanRecord{}
	//根据id查询
	err := lib.Engine.Table("plan_record").ID(userId).Find(&records)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(records))
}

//将标准时间转换成时间戳
func timeToUnix(formTime string) int64 {

	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, formTime, loc)
	sr := theTime.Unix()
	//打印输出时间戳 1420041600
	return sr
}

//设置动态路径函数。先根据定位点获取最近的围栏，然后获取该围栏的所有的线段和点。接着给所有的点加上线段字段，线段加上点字段。
//主要算法：查询当前点的线段是否已经过，经过+1,如果线段经过次数较少，加入min_line。然后从min_line随机选取一个线段。并将该线段距离加上里程。
// 从选择的线段的不同点作为下一个点，往下循环，直到total超过要求总里程。
//参数：路径要求总里程单位KM，定位点经纬度
//返回参数：路径点，错误。
func DynamicRoute(total float64, startLongitude float64, startLatitude float64) ([]*models.PlanPoints, error) {

	//1.根据经纬度获取最近的围栏---------------开始
	fenceInfo := make([]models.PlanFence, 0)
	err := lib.Engine.Table("plan_fence").Find(&fenceInfo)
	if err != nil {
		println("查询围栏信息出错")
		return nil, err
	}
	if len(fenceInfo) == 0 {
		return nil, errors.New("没有电子围栏")
	}
	//遍历围栏的中心点，查询最近的围栏，len(fenceInfo)
	distance := make([]float64, len(fenceInfo))
	closestDistance := distance[0] //最近的围栏
	cloestFence := fenceInfo[0]
	//所有围栏中心点与定位点的距离
	for i := 0; i < len(fenceInfo); i++ {
		distance[i] = lib.GetDistance(startLongitude, fenceInfo[i].Longitude, startLatitude, fenceInfo[i].Latitude)
	}
	//寻找最近的围栏并赋值给closestFence
	min := distance[0]
	for i2 := 0; i2 < len(fenceInfo); i2++ {
		if distance[i2] < min {
			min = distance[i2]
			cloestFence = fenceInfo[i2]
		}
	}
	// 如果最近距离的围栏还是太远比如5KM直接返回，距离太远
	if closestDistance > 3 {
		return nil, errors.New("附近没有围栏")
	}
	//获取最近的围栏----------------------结束

	//获取所有点和线段-----------------------开始
	all_points := make([]models.PlanPoints, 0)
	err = lib.Engine.Table("plan_points").Where("fence_id=?", cloestFence.Id).Cols("id", "name", "fence_id", "longitude", "latitude", "type", "uuid", "create_at").Find(&all_points)
	if err != nil {
		return nil, err
	}
	println("围栏内点位数量：", len(all_points))
	if len(all_points) == 0 {
		return nil, errors.New("围栏内没有点位，请先设置点位")
	}
	if len(all_points) < 4 {
		return nil, errors.New("围栏内点位数量少于4个，请先设置点位")
	}
	all_lines := make([]models.PlanLine, 0)
	err = lib.Engine.Table("plan_line").Where("fence_id=?", cloestFence.Id).Cols("id", "fence_id", "point1", "point2", "create_at").Find(&all_lines)
	if err != nil {
		return nil, err
	}
	println("围栏内线段数量：", len(all_lines))
	if len(all_lines) == 0 {
		return nil, errors.New("围栏内没有线段，请先设置线段")
	}
	if len(all_lines) < 4 {
		return nil, errors.New("围栏内线段数量少于4个，请先设置线段")
	}
	//获取所有线段--------------------------结束

	//1、找到最近的点（这个点必须有线）
	var pt *models.PlanPoints
	//获取最近的点-----------------------开始
	distanceArr := make([]float64, len(all_points))
	for i := 0; i < len(all_points); i++ {
		distanceArr[i] = lib.GetDistance(startLatitude, all_points[i].Latitude, startLongitude, all_points[i].Longitude)
		//println("第index:", i, "个距离，distance:", distanceArr[i])
		println("")
		fmt.Printf("第%d个点，距离定位点距离%.2f", i, distanceArr[i])
	}
	//寻找最近的围栏并赋值给closestFence
	println("开始寻找最近的点：")
	firstMin := distanceArr[0]                //所有点与定位点的距离数组，设置最近距离为数组第一个
	cloestPoint := all_points[0]              //初始设置第一个点为距离最近的点
	for i2 := 0; i2 < len(all_points); i2++ { //遍历所有的点
		//println("firstMin:", min, "index distance:", distanceArr[i2])
		println("")
		if distanceArr[i2] < firstMin {
			fmt.Printf("当前点距离比前一个点距离更近，当前点的距离：%.2f,最近点的距离：%.2f", distanceArr[i2], firstMin)
			fmt.Printf("当前点：%v,距离最近的点：%v", all_points[i2], cloestPoint)
			firstMin = distanceArr[i2]
			cloestPoint = all_points[i2]

		} else {
			fmt.Printf("当前点距离比前一个点距离更远，当前点的距离：%.2f,最近点的距离：%.2f", distanceArr[i2], firstMin)
			fmt.Printf("当前点：%v,距离最近的点：%v", all_points[i2], cloestPoint)
		}

	}

	//获取最近的点-----------------------结束
	fmt.Printf("最近的点：%v", cloestPoint)
	println("寻找最近的点结束")
	//点和线的关系,给所有点加上线段，给所有线段的点补齐所有属性
	for pIndex, p := range all_points {
		////找到最近点
		if p.Id == cloestPoint.Id {
			pt = &all_points[pIndex]
		}
		//p.Lines = make([]*models.PlanLine, 0)
		for index, l := range all_lines {
			if l.Point1 == p.Id {
				//这里给的是复制的l.Point1Ptr,而不是原来数组的l.Point1Ptr
				all_points[pIndex].Lines = append(all_points[pIndex].Lines, &all_lines[index])
				all_lines[index].Point1Ptr = &all_points[pIndex]
			} else if l.Point2 == p.Id {
				all_points[pIndex].Lines = append(all_points[pIndex].Lines, &all_lines[index])
				all_lines[index].Point2Ptr = &all_points[pIndex]
			}
		}
	}

	//清空计数？
	for _, l := range all_lines {
		l.CoverTimes = 0
	}

	//所有轨迹点
	points := make([]*models.PlanPoints, 0)
	points = append(points, pt)

	//长度
	length := 0.00

	//2、生成路径
	for length < total {
		//从经过次数最少的线中随机选择
		var min_lines = make([]*models.PlanLine, 0)
		var min_cover = 10000

		for _, l := range pt.Lines { //将经过最少的线加入min_lines
			//println("l.Covertimes:",l.CoverTimes,",min_cover:",min_cover)
			if l.CoverTimes < min_cover {
				min_lines = make([]*models.PlanLine, 0)
				min_cover = l.CoverTimes
			}
			if l.CoverTimes == min_cover {
				min_lines = append(min_lines, l)
				l.CoverTimes = l.CoverTimes + 1 //走过的路线加一
			}
		}

		if len(min_lines) == 0 {
			return nil, errors.New("点的线段只有一条,min_lines length =0")
		}
		//随机取线
		rand.Seed(time.Now().Unix())

		line := models.PlanLine{}
		if len(min_lines) == 0 {
			line = *min_lines[0]
			fmt.Printf("点的线段只有一条,min_lines length =0")
			//return nil, errors.New("点的线段只有一条,min_lines length =0")
		} else {
			idx := rand.Intn(len(min_lines))
			line = *min_lines[idx]
		}

		//向前移动
		if line.Point1 == pt.Id {
			pt = line.Point2Ptr
		} else if line.Point2 == pt.Id {
			pt = line.Point1Ptr
		}
		println("")
		//fmt.Printf("加点前%v：", points)
		//TODO 当出现线段的point1和point2都是一个点时会出错。记得让前端检查数据不要传这种点。
		points = append(points, pt)
		if pt == nil {
			fmt.Printf("当前点的内存地址是空！%v:", pt)
			return nil, errors.New("当前点的内存地址是空！可能该线段的点已被删除，但线段仍存在。")
		}
		println("")
		//fmt.Printf("加点后%v:", points)
		length = length + lib.GetDistance(line.Point1Ptr.Latitude, line.Point2Ptr.Latitude, line.Point1Ptr.Longitude, line.Point2Ptr.Longitude)

	}

	return points, nil
}
