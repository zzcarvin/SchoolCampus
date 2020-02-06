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

type requestFinishRun struct {
	PlanStartId  int             `json:"plan_start_id"`
	SchoolId     int             `json:"school_id"`
	PlanId       int             `json:"plan_id"`
	StudentId    int             `json:"student_id"`
	StartTime    string          `json:"start_time"`
	EndTime      string          `json:"end_time"`
	Distance     int             `json:"distance"`
	Duration     int             `json:"duration"`
	Times        int             `json:"times"`
	Calories     float64         `json:"calories"`
	Steps        int             `json:"steps"`
	Pace         int             `json:"pace"`
	Speed        float64         `json:"speed"`
	Paces        []int           `json:"paces"`
	Points       []models.Points `json:"points"`
	PointsStatus bool            `json:"points_status"`

	//增加经过点
	PassPoints []models.Points `json:"pass_points"`

	FaceStatus bool `json:"face_status" xorm:"-"`

	IBeacon []models.Points `json:"ibeacon"`
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

//计算每周次数和公里数
type dayData struct {
	times    int
	distance int
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

	//重新计算时间，以跑步开始时间和结束时间为准

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
	//把所有记录记录
	logJson, _ := json.Marshal(finishRun)
	lib.RecordLogger.Info("planRecord：" + string(logJson))

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

	//补充保存打卡点---开始

	sqlString := ""
	for _, value := range start_planTime.Route {
		sqlString = sqlString + " or id=" + strconv.FormatInt(value, 10)
	}
	orString := sqlString[3:]

	//1.获取围栏内点的详细信息。记得以后给所有点加软删
	routePoints := make([]iBeaconsPoints, 0)
	alliBeaconsErr := lib.Engine.Table("plan_points").Where(orString).Find(&routePoints)
	if alliBeaconsErr != nil {
		ctx.JSON(lib.NewResponseFail(1, "查询打卡点失败"))
		return
	}

	//将所有经过的打卡点入库
	passPointIds := make([]int, 0)
	for _, value := range finishRun.PassPoints {
		for _, value1 := range routePoints {
			if value.Longitude == value1.Longitude && value.Latitude == value1.Latitude {
				passPointIds = append(passPointIds, value1.Id)
				break
			}
		}
	}

	//检查是否有蓝牙点位，如果有蓝牙点位，再去检测iBeacon点位至少有一个,如果没有就是报蓝牙点位没有打到，invalid_code:9
	// TODO: 后续如果android把蓝牙问题解决了，这边需要收集用户所有的蓝牙点位，然后和传过来的蓝牙点位，全部比对成功才算蓝牙都打到
	//isHaveBeacons := false
	//for _, vv := range routePoints {
	//	if vv.Type == 1 {
	//		isHaveBeacons = true
	//	}
	//}

	//获取计划详情
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

	//查询该用户的院系以及性别

	//无效判断---开始
	//默认最低最高限制
	//正常最低minPace := 720

	//*************更改为男女****************
	conf := configs.Conf
	//男女配速默认都是0，如果没有值就不做有效无效判断
	minPace := 0 //最慢配速
	maxPace := 0 //最快配速
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
	fmt.Printf("minpace:%v,maxpace:%v,minFrequency:%v,maxFrequency:%v",
		minPace, maxPace, minFrequency, maxFrequency)

	//获取单次最低公里数---开始
	var minDistance int
	minDistance = plan.MinSingleDistance

	println("目标最低公里数：", minDistance, "运动公里数：", finishRun.Distance)
	//获取单次最低公里数---结束

	//无效判断
	invalidCode := make([]int, 0)
	//判断配速是否正常
	if minPace != 0 {
		if finishRun.Pace > minPace { //配速太慢
			invalidCode = append(invalidCode, 1)
		}
	}

	if maxPace != 0 {
		if finishRun.Pace < maxPace { //配速太快
			invalidCode = append(invalidCode, 2)
		}
	}

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
	if finishRun.FaceStatus == false {
		invalidCode = append(invalidCode, 8) //跑步结束没有通过人脸比对
	}
	//if isHaveBeacons == true {
	//	if len(finishRun.IBeacon) == 0 {
	//		invalidCode = append(invalidCode, 9) //有蓝牙点位，至少打一个
	//	}
	//}

	//新加无效---开始
	nowTime := time.Now()
	timeLong := int(lib.GetSecond(start_planTime.CreateAt, nowTime))

	if plan.MaxTimeLong != 0 {
		if timeLong > plan.MaxTimeLong {
			invalidCode = append(invalidCode, 10) //超过最长时间
		}
	}
	if plan.MinTimeLong != 0 {
		if timeLong < plan.MinTimeLong {
			invalidCode = append(invalidCode, 11) //低于最短时间
		}
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
		PassPoints:   passPointIds,
	}
	//计算运动时间间隔
	record.Times = record.EndTime.Sub(record.StartTime).Hours()

	//添加运动记录
	res, err := lib.Engine.Table("plan_record").Insert(&record)
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

	//新加无效---结束

	if len(invalidCode) == 0 {
		//修改多计划开始***********************************************************************************
		//recordid := record.Id
		//key := "finishRunJob"
		//fmt.Printf("已推送 recordid: %d", recordid)
		//lib.MainLogger.Info("PutJob：" + lib.Int2str(recordid))
		//err := PutJob(recordid, key)
		//if err != nil {
		//	ctx.JSON(lib.NewResponseFail(1, err.Error()))
		//	return
		//}
		//修改多计划结束**********************************************************************************

		UpdateProgressJob(record, plan)

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

type Job struct {
	Plan_record_id int    `json:"Plan_record_id"`
	Args           string `json:"Args"`
}

//type Job struct {
//	Class string        `json:"Class"`
//	Args  string        `json:"Args"`
//}
//计划记录id入队列函数
func PutJob(recordid int, key string) error {
	println("*******put*******")
	fmt.Printf("record_id:%d", recordid)
	c := lib.GetRedisConn()
	defer c.Close()
	//b, err := json.Marshal(job)
	//if err != nil {
	//	fmt.Println("PutJob error:", job.Args, err)
	//	return err
	//}
	_, err := c.Do("RPUSH", key, recordid)
	if err != nil {
		fmt.Println("PutJob error:", recordid, err)
		return err
	}
	return nil
}

//合并周进度和学期进度
//1.计算本周进度和本学期进度
//2.计算进度完成度，不满足总里程，总次数，周次数，周进度，必跑日，进度状态为0，都满足改为1

//分成两部分
//1.获取更新本周数据和本学期数据
//1.1超过日跑次数不更新周进度和学期进度，直接获取更新进度完成状态。
//1.2获取本周数据根据创建时间排序，本周开始时间和结束时间，遍历本周数据，将每天的数据放进map[2019/12/05]{times,distance}里，超过每天次数，map[2019/12/05]int不累加。
//获取本周一，获取本周二，获取本周三以此类推。
// 计算本周的次数和里程数，超过本周最高或本周最高里程以最高为准。超过本周次数和里程时本学期次数和里程也不累加。
//2.根据总次数总里程，每周次数，每周里程，必跑日，获取更新进度完成状态。
//2.1总次数，总里程判断。
//2.2必跑日判断。
//2.3每周次数，公里数判断
//三项都满足更新状态为完成，否则为未完成。
func GetJob(key string) {

	c := lib.GetRedisConn()
	defer c.Close()

	reply, err := c.Do("LPOP", key)
	if err != nil || reply == nil {
		if err != nil {
			fmt.Println("GetJob error:", err)
			//保留错误
			lib.MainLogger.Info("GetJob Error：" + err.Error())
		}

		//return 0, errors.New("no Job.")
	} else {

		var j int
		decoder := json.NewDecoder(bytes.NewReader(reply.([]byte)))
		if err := decoder.Decode(&j); err != nil {
			lib.MainLogger.Info("Job Decode Error:" + err.Error())
			//lib
		}
		fmt.Println("获取到数据，开始计算", j)
		if j != 0 {
			//UpdateProgressJob(j)
		}

	}

}

//分成两部分
//1.获取更新本周数据和本学期数据
//1.1超过日跑次数不更新周进度和学期进度，直接获取更新进度完成状态。
//1.2获取本周数据根据创建时间排序，本周开始时间和结束时间，遍历本周数据，将每天的数据放进map[2019/12/05]{times,distance}里，超过每天次数，map[2019/12/05]int不累加。
//获取本周一，获取本周二，获取本周三以此类推。
// 计算本周的次数和里程数，超过本周最高或本周最高里程以最高为准。超过本周次数和里程时本学期次数和里程也不累加。
//2.根据总次数总里程，每周次数，每周里程，必跑日，获取更新进度完成状态。
//2.1总次数，总里程判断。
//2.2必跑日判断。
//2.3每周次数，公里数判断
//三项都满足更新状态为完成，否则为未完成。

//一周次数和公里数
//入参：一周的某一时刻
//返回参数：一周的次数和公里数
func getWeekData(present time.Time, weekRecords []models.PlanRecord, plan models.Plan) (times int, distance int) {
	fmt.Printf("运动记录创建时间：%v", present)
	//获取周一至周日的日期
	monday := lib.GetFirstDateOfWeek(present)
	mondayStr := monday.Format("2006-01-02")
	tuesdayStr := monday.AddDate(0, 0, 1).Format("2006-01-02")
	wednesdayStr := monday.AddDate(0, 0, 2).Format("2006-01-02")
	thursdayStr := monday.AddDate(0, 0, 3).Format("2006-01-02")
	FridayStr := monday.AddDate(0, 0, 4).Format("2006-01-02")
	SaturdayStr := monday.AddDate(0, 0, 5).Format("2006-01-02")
	SundayStr := monday.AddDate(0, 0, 6).Format("2006-01-02")

	//println("周一：", mondayStr)

	newWeekData := make([]dayData, 7)
	//for i:=0;i<6;i++{
	//	newWeekData= append(newWeekData, dayData{})
	//}
	//newWeekData= append(newWeekData, dayData{})
	fmt.Printf("周五的数据：%v", newWeekData[4])
	//遍历本周数据,累加每天数据
	for _, value := range weekRecords {
		switch value.CreateAt.Format("2006-01-02") {
		case mondayStr:
			dayDataDetail(plan.MaxDayTimes, plan.MaxSingleDistance, value.Distance, &newWeekData[0])
		case tuesdayStr:
			dayDataDetail(plan.MaxDayTimes, plan.MaxSingleDistance, value.Distance, &newWeekData[1])
		case wednesdayStr:
			dayDataDetail(plan.MaxDayTimes, plan.MaxSingleDistance, value.Distance, &newWeekData[2])
		case thursdayStr:
			dayDataDetail(plan.MaxDayTimes, plan.MaxSingleDistance, value.Distance, &newWeekData[3])
		case FridayStr:
			dayDataDetail(plan.MaxDayTimes, plan.MaxSingleDistance, value.Distance, &newWeekData[4])
		case SaturdayStr:
			dayDataDetail(plan.MaxDayTimes, plan.MaxSingleDistance, value.Distance, &newWeekData[5])
		case SundayStr:
			dayDataDetail(plan.MaxDayTimes, plan.MaxSingleDistance, value.Distance, &newWeekData[6])

		}
	}

	//遍历每天数据，获取本周所有次数和公里数，超过本周最高限制使用最高限制
	weekTimes := 0
	weekDistance := 0
	for _, value := range newWeekData {
		if plan.MaxWeekTimes != 0 && weekTimes >= plan.MaxWeekTimes {
			weekTimes = plan.MaxWeekTimes
		} else {
			weekTimes += value.times
		}

		if plan.MaxWeekDistance != 0 && weekDistance >= plan.MaxWeekDistance {
			weekDistance = plan.MaxWeekDistance
		} else {
			weekDistance += value.distance
		}
	}
	//fmt.Printf("weekTimes：%d,weekDistance:%d", weekTimes, weekDistance)
	return weekTimes, weekDistance
}

//每天的详细数据
func dayDataDetail(maxDayTimes int, maxSingleDistance int, recordDistance int, dayData *dayData) {

	//超过单日最高次数，不累加直接返回
	if maxDayTimes != 0 && dayData.times >= maxDayTimes {
		dayData.times = maxDayTimes
		return
	} else {
		dayData.times++
		//超过单次公里数，只累加单次最大公里数
		if maxSingleDistance != 0 && recordDistance >= maxSingleDistance {
			dayData.distance += maxSingleDistance
		}
		dayData.distance += recordDistance
		//fmt.Printf("次数%d,公里数：%d", dayData.times, dayData.distance)
	}

	return
}

func UpdateProgressJob(record1 models.PlanRecord, plan1 models.Plan) {
	lib.MainLogger.Info("UpdateProgressJob Start：" + lib.Int2str(record1.Id))
	//查询进度
	planProgress1 := models.PlanProgress{}
	_, err := lib.Engine.Table("plan_progress").Where("student_id=?", record1.StudentId).And("plan_id=?", record1.PlanId).Get(&planProgress1)
	if err != nil {
		lib.MainLogger.Info("查询周计划错误：" + err.Error())
		return
	}

	//根据计划限制，调整数据---开始
	//学期进度数据--开始
	progressFlg := true
	//获取当天有效跑步次数
	var dayRecords []models.PlanRecord
	errRecords := lib.Engine.Table("plan_record").Where("status=?", 1).And("student_id=?", record1.StudentId).
		And("plan_id=?", record1.PlanId).And("TO_DAYS(create_at) = TO_DAYS(NOW())").Find(&dayRecords)
	if errRecords != nil {
		fmt.Printf("获取运动记录错误：%v", errRecords)
		return
	}
	//超过日跑次数，不累计进度
	if len(dayRecords) > plan1.MaxDayTimes && plan1.MaxDayTimes != 0 {
		progressFlg = false
		lib.MainLogger.Info("超过当日上线：：" + lib.Int2str(record1.Id))
	}

	//未超过日跑次数更新
	if progressFlg {
		//本周数据
		now := record1.CreateAt
		fmt.Printf("运动记录时间：%v", now)
		var weekRecords []models.PlanRecord
		errRecords := lib.Engine.Table("plan_record").Where("status=?", 1).And("student_id=?", record1.StudentId).
			And("plan_id=?", record1.PlanId).And("YEARWEEK(create_at,1) = YEARWEEK(NOW(),1)").Find(&weekRecords)
		if errRecords != nil {
			fmt.Printf("获取运动记录错误：%v", errRecords)
			return
		}

		weekTimes, weekDistance := getWeekData(now, weekRecords, plan1)
		println("week_times:", weekTimes, "weekDistance:", weekDistance)

		newProgress := models.PlanProgress{
			Duration:     record1.Duration + planProgress1.Duration,
			Calories:     record1.Calories + planProgress1.Calories,
			Steps:        record1.Steps + planProgress1.Steps,
			WeekDistance: weekDistance,
			WeekTimes:    weekTimes,
		}

		//根据计划限制，调整数据---结束

		//1.比跑日，验证值：mustRunDayBool
		mustRunDayBool := false
		if plan1.MustRunDay != "0" {

			var allRecords []models.PlanRecord
			errAllRecords := lib.Engine.Table("plan_record").Where("status=?", 1).And("student_id=?", record1.StudentId).
				And("plan_id=?", record1.PlanId).Find(&allRecords)
			if errAllRecords != nil {
				fmt.Printf("获取运动记录错误：%v", errRecords)
				return
			}

			//println("打印学生", record1.StudentId, "的所有运动记录：***************************")
			//for index, value := range allRecords {
			//	fmt.Printf("第%d个运动记录：%v", index, value)
			//}

			//必跑日，根据数字判断周几，循环每周数据判断是否每周的那天是否跑了而且有效！
			//1.获取每个必跑日的时间，然后查询时间运动记录，
			planStartWeekDayInt := int(plan1.DateBegin.Weekday())
			if plan1.DateBegin.Weekday() == 0 {
				planStartWeekDayInt = 7
			}

			//获取第一个必跑日，获取第一周开始时间，并获取周几，
			//必跑日小于开始日，
			planMustRunDayInt, err := strconv.Atoi(plan1.MustRunDay)
			var planMustRunUnix int64
			if err != nil {
				fmt.Printf("转换字符错误：%v", err)
				return
			}

			if planMustRunDayInt < planStartWeekDayInt {
				//必跑日小于开始日
				planMustRunUnix = plan1.DateBegin.Unix() + 7*86400 - int64((planStartWeekDayInt-planMustRunDayInt)*86400)
			} else if planMustRunDayInt == planStartWeekDayInt {
				//必跑日==开始日
				planMustRunUnix = plan1.DateBegin.Unix()
			} else {
				//必跑日大于开始日
				planMustRunUnix = plan1.DateBegin.Unix() - int64((planMustRunDayInt-planStartWeekDayInt)*86400)
			}

			mustRunDayUnixArr := make([]int64, 0)
			//必跑日数组
			for mustUnix := planMustRunUnix; mustUnix < plan1.DateBegin.Unix(); mustUnix += 7 * 86400 {
				mustRunDayUnixArr = append(mustRunDayUnixArr, mustUnix)
			}

			//判断每周必跑日是否都跑了
			mustDayRecordNumber := 0
			for i := 0; i < len(mustRunDayUnixArr)/2; i++ {
				for _, valueAllRecords := range allRecords {
					if valueAllRecords.CreateAt.Unix() >= mustRunDayUnixArr[i*2] && valueAllRecords.CreateAt.Unix() <= mustRunDayUnixArr[i*2+1] {
						mustDayRecordNumber += 1
						break
					}
				}

			}

			if mustDayRecordNumber == len(mustRunDayUnixArr) {
				mustRunDayBool = true
			}

		}

		//2.总里程，总次数，验证值：totalDistanceBool，totalTimesBool
		totalDistanceBool := false
		totalTimesBool := false
		println(totalDistanceBool, totalTimesBool)
		if planProgress1.Times >= plan1.TotalTimes {
			totalTimesBool = true
		}
		if planProgress1.Distance >= plan1.TotalDistance {
			totalDistanceBool = true
		}

		//3.所有周的最低跑步次数，最低跑步里程，有一周没完成就让总进度保持不到100，且status为0
		//所有周数据,验证值：allWeekStatus
		//每周的运动记录
		weeksRecords := lib.GetWeeksRecords(record1.StudentId, plan1)
		weeksRecords = weeksRecords[0 : len(weeksRecords)-1]
		//每周的公里数，次数，状态
		allWeeksData := make([]lib.SimpleWeekContent, len(weeksRecords))
		//完成的周数
		weeksStatus := 0
		allWeekStatus := false
		for index, value := range weeksRecords {
			//本周有运动记录
			if len(value) != 0 {
				allWeeksData[index].Sequence = index
				allWeeksData[index].Times, allWeeksData[index].Distance = getWeekData(value[0].CreateAt, value, plan1)
				if allWeeksData[index].Times >= plan1.MinWeekTimes && allWeeksData[index].Distance >= plan1.MinWeekDistance {
					allWeeksData[index].Status = 1
				} else {
					allWeeksData[index].Status = 0
				}
			} else {
				//本周没有运动记录
				allWeeksData[index].Status = 0
			}
			if allWeeksData[index].Status == 1 {
				weeksStatus++
			}
		}
		if weeksStatus == len(weeksRecords) {
			allWeekStatus = true
		}

		//完成的周数
		newProgress.WeekCompleteProgress = fmt.Sprintf("%d/%d", weeksStatus, len(weeksRecords))

		//使用每周的公里数和次数累计作为总次数和总公里数
		for _, value := range allWeeksData {
			newProgress.Times += value.Times
			newProgress.Distance += value.Distance
		}
		println("times:", newProgress.Times, "distance:", newProgress.Distance)
		newProgress.CompleteProgress = float32(newProgress.Times) / float32(plan1.TotalTimes)

		//全部条件满足后，更新进度状态。
		if mustRunDayBool && totalTimesBool && totalDistanceBool && allWeekStatus {
			newProgress.Status = 1
		} else {
			newProgress.Status = 0
		}

		//debug信息
		//fmt.Println(fmt.Sprintf("周距离%d", newProgress.WeekDistance))
		//fmt.Println(fmt.Sprintf("周次数%d", newProgress.WeekTimes))
		//fmt.Println(fmt.Sprintf("周完成进度%s", newProgress.WeekCompleteProgress))
		//fmt.Println(fmt.Sprintf("总距离%d", newProgress.Distance))
		//fmt.Println(fmt.Sprintf("总次数%d", newProgress.Times))
		//fmt.Println(fmt.Sprintf("总完成进度%f", newProgress.CompleteProgress))
		//fmt.Println(fmt.Sprintf("总耗时%d", newProgress.Duration))
		//fmt.Println(fmt.Sprintf("总卡路里%f", newProgress.Calories))
		//fmt.Println(fmt.Sprintf("总步数%d", newProgress.Steps))
		//fmt.Println(fmt.Sprintf("完成状态%d", newProgress.Status))

		res, err := lib.Engine.Table("plan_progress").Where("plan_id=?", record1.PlanId).
			And("student_id=?", record1.StudentId).
			Cols("week_distance", "week_times", "week_complete_progress", "distance", "times", "complete_progress", "duration", "calories", "steps", "status").
			Update(&newProgress)
		if err != nil {
			lib.MainLogger.Info("更新计划进度错误：" + lib.Int2str(record1.Id) + err.Error())
			return
		}
		if res == 1 {
			lib.MainLogger.Info("Success:UpdateProgressJob End：" + lib.Int2str(record1.Id))
		}
		if res == 0 {
			lib.MainLogger.Info("Fail:UpdateProgressJob End：" + lib.Int2str(record1.Id))
		}

	}
}

//用来初始化学生跑步进度
func UpdateStudentProgress(studentId int, plan1 models.Plan) {
	planProgress1 := models.PlanProgress{}
	planId := plan1.Id
	allData, err := lib.Engine.Table("plan_record").
		Where("student_id=?", studentId).
		And("plan_id=?", planId).And("status=1").
		Sums(&planProgress1, "duration", "calories", "steps")

	if err != nil {
		fmt.Printf("查询进度有问题：%v", err)
		return
	}

	weekRecords := []models.PlanRecord{}
	errRecords := lib.Engine.Table("plan_record").Where("status=?", 1).And("student_id=?", studentId).
		And("plan_id=?", planId).And("YEARWEEK(create_at,1) = YEARWEEK(NOW(),1)").Find(&weekRecords)
	if errRecords != nil {
		fmt.Printf("获取运动记录错误：%v", errRecords)
		return
	}

	//本周数据
	now := time.Now()
	fmt.Printf("运动记录时间：%v", now)
	weekTimes, weekDistance := getWeekData(now, weekRecords, plan1)
	println("week_times:", weekTimes, "weekDistance:", weekDistance)

	newProgress := models.PlanProgress{
		Duration:     int(allData[0]),
		Calories:     float64(allData[1]),
		Steps:        int(allData[2]),
		WeekDistance: weekDistance,
		WeekTimes:    weekTimes,
	}

	//1.比跑日，验证值：mustRunDayBool
	mustRunDayBool := false
	if plan1.MustRunDay != "0" {

		var allRecords []models.PlanRecord
		errAllRecords := lib.Engine.Table("plan_record").Where("status=?", 1).And("student_id=?", studentId).
			And("plan_id=?", planId).Find(&allRecords)
		if errAllRecords != nil {
			fmt.Printf("获取运动记录错误：%v", errAllRecords)
			return
		}

		//必跑日，根据数字判断周几，循环每周数据判断是否每周的那天是否跑了而且有效！
		//1.获取每个必跑日的时间，然后查询时间运动记录，
		planStartWeekDayInt := int(plan1.DateBegin.Weekday())
		if plan1.DateBegin.Weekday() == 0 {
			planStartWeekDayInt = 7
		}

		//获取第一个必跑日，获取第一周开始时间，并获取周几，
		//必跑日小于开始日，
		planMustRunDayInt, err := strconv.Atoi(plan1.MustRunDay)
		var planMustRunUnix int64
		if err != nil {
			fmt.Printf("转换字符错误：%v", err)
			return
		}

		if planMustRunDayInt < planStartWeekDayInt {
			//必跑日小于开始日
			planMustRunUnix = plan1.DateBegin.Unix() + 7*86400 - int64((planStartWeekDayInt-planMustRunDayInt)*86400)
		} else if planMustRunDayInt == planStartWeekDayInt {
			//必跑日==开始日
			planMustRunUnix = plan1.DateBegin.Unix()
		} else {
			//必跑日大于开始日
			planMustRunUnix = plan1.DateBegin.Unix() - int64((planMustRunDayInt-planStartWeekDayInt)*86400)
		}

		mustRunDayUnixArr := make([]int64, 0)
		//必跑日数组
		for mustUnix := planMustRunUnix; mustUnix < plan1.DateBegin.Unix(); mustUnix += 7 * 86400 {
			mustRunDayUnixArr = append(mustRunDayUnixArr, mustUnix)
		}

		//判断每周必跑日是否都跑了
		mustDayRecordNumber := 0
		for i := 0; i < len(mustRunDayUnixArr)/2; i++ {
			for _, valueAllRecords := range allRecords {
				if valueAllRecords.CreateAt.Unix() >= mustRunDayUnixArr[i*2] && valueAllRecords.CreateAt.Unix() <= mustRunDayUnixArr[i*2+1] {
					mustDayRecordNumber += 1
					break
				}
			}

		}

		if mustDayRecordNumber == len(mustRunDayUnixArr) {
			mustRunDayBool = true
		}

	}

	//2.总里程，总次数，验证值：totalDistanceBool，totalTimesBool
	totalDistanceBool := false
	totalTimesBool := false

	//3.所有周的最低跑步次数，最低跑步里程，有一周没完成就让总进度保持不到100，且status为0
	//所有周数据,验证值：allWeekStatus
	//每周的运动记录
	weeksRecords := lib.GetWeeksRecords(studentId, plan1)
	weeksRecords = weeksRecords[0 : len(weeksRecords)-1]
	//每周的公里数，次数，状态
	allWeeksData := make([]lib.SimpleWeekContent, len(weeksRecords))
	//完成的周数
	weeksStatus := 0
	allWeekStatus := false
	for index, value := range weeksRecords {
		//本周有运动记录
		if len(value) != 0 {
			allWeeksData[index].Sequence = index
			allWeeksData[index].Times, allWeeksData[index].Distance = getWeekData(value[0].CreateAt, value, plan1)
			if allWeeksData[index].Times >= plan1.MinWeekTimes && allWeeksData[index].Distance >= plan1.MinWeekDistance {
				allWeeksData[index].Status = 1
			} else {
				allWeeksData[index].Status = 0
			}
		} else {
			//本周没有运动记录
			allWeeksData[index].Status = 0
		}
		if allWeeksData[index].Status == 1 {
			weeksStatus++
		}
	}
	if weeksStatus == len(weeksRecords) {
		allWeekStatus = true
	}

	//完成的周数
	newProgress.WeekCompleteProgress = fmt.Sprintf("%d/%d", weeksStatus, len(weeksRecords))

	//使用每周的公里数和次数累计作为总次数和总公里数
	for _, value := range allWeeksData {
		newProgress.Times += value.Times
		newProgress.Distance += value.Distance
	}
	if newProgress.Times >= plan1.TotalTimes {
		totalTimesBool = true
	}
	if newProgress.Distance >= plan1.TotalDistance {
		totalDistanceBool = true
	}

	println("times:", newProgress.Times, "distance:", newProgress.Distance)
	newProgress.CompleteProgress = float32(newProgress.Times) / float32(plan1.TotalTimes)

	//全部条件满足后，更新进度状态。
	if mustRunDayBool && totalTimesBool && totalDistanceBool && allWeekStatus {
		newProgress.Status = 1
	} else {
		newProgress.Status = 0
	}

	fmt.Println(fmt.Sprintf("周距离%d", newProgress.WeekDistance))
	fmt.Println(fmt.Sprintf("周次数%d", newProgress.WeekTimes))
	fmt.Println(fmt.Sprintf("周完成进度%s", newProgress.WeekCompleteProgress))
	fmt.Println(fmt.Sprintf("总距离%d", newProgress.Distance))
	fmt.Println(fmt.Sprintf("总次数%d", newProgress.Times))
	fmt.Println(fmt.Sprintf("总完成进度%f", newProgress.CompleteProgress))
	fmt.Println(fmt.Sprintf("总耗时%d", newProgress.Duration))
	fmt.Println(fmt.Sprintf("总卡路里%f", newProgress.Calories))
	fmt.Println(fmt.Sprintf("总步数%d", newProgress.Steps))
	fmt.Println(fmt.Sprintf("完成状态%d", newProgress.Status))

	res, err := lib.Engine.Table("plan_progress").Where("plan_id=?", planId).
		And("student_id=?", studentId).
		Cols("week_distance", "week_times", "week_complete_progress", "distance", "times", "complete_progress", "duration", "calories", "steps", "status").
		Update(&newProgress)
	if err != nil {
		fmt.Printf("更新计划进度错误：%v", err)
		return
	}
	if res == 1 {
		println("学生计划进度更新成功")
	}
	if res == 0 {
		println("学生计划进度为最新")
	}
}
