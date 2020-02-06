package progress

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"bytes"
	"fmt"
	"github.com/go-xorm/builder"
	"strconv"

	"github.com/kataras/iris"
	"log"
	"strings"
	"time"
)

// swagger:parameters  ProgressCreateRequest
type ProgressCreateRequest struct {
	// in: body
	Body models.PlanProgress
}

// 响应结构体
//
// swagger:response    ProgressCreateResponse
type ProgressCreateResponse struct {
	// in: body
	Body pointsresponseMessage
}
type pointsresponseMessage struct {
	// Required: true
	models.ResponseType
	Data models.PlanProgress
}

//计划进度次数返回体
//次数单个结构体
type recordTimes struct {
	CompleteTimes []int `json:"complete_times"`
	ValidateTimes []int `json:"validate_times"`
}

//距离单个结构体
type recordDistance struct {
	completeDistance []int `json:"complete_distance"`
}

//最外层返回结构体
type responseRecord struct {
	TimesData    recordTimes `json:"times_data"`    //次数数组，x轴数据
	DistanceData []int       `json:"distance_data"` //距离数组，y轴数据
	XUnit        []string    `json:"x_unit"`        //x轴单位
	//YUnit []string `json:"y_unit"`//y轴单位
}

func create(ctx iris.Context) {
	// swagger:route   /api/plan/progress progress ProgressCreateRequest
	//
	// 创建计划
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: ProgressCreateResponse
	progress := models.PlanProgress{}
	err := ctx.ReadJSON(&progress)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return

	}

	//ctx.JSON(lib.NewResponseOK(planFence))

	res, err := lib.Engine.Table("plan_progress").Insert(&progress)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	lib.NewResponseOK(res)
	ctx.JSON(lib.NewResponseOK(progress))
}

// swagger:route DELETE /api/plan/progress:id  progress ProgressDelete
//
//	 删除计划 要带id过来
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
	PlanProgress := models.PlanProgress{}
	affected, err := lib.Engine.Table("plan_progress").ID(id).Delete(&PlanProgress)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	ctx.JSON(lib.NewResponseOK(affected))
}

// swagger:parameters  ProgressUpdateRequest
type ProgressUpdateRequest struct {
	// in: body
	Body models.PlanProgress
}

// 响应结构体
//
// swagger:response    ProgressUpdateResponse
type ProgressUpdateResponse struct {
	// in: body
	Body pointsresponseMessage
}

func update(ctx iris.Context) {
	// swagger:route PUT /api/plan/progress progress ProgressUpdateRequest
	// 修改学生计划
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: ProgressUpdateResponse
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	PlanProgress := models.PlanProgress{}

	//解析student
	err := ctx.ReadJSON(&PlanProgress)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//TODO 验证数据有效性

	//插入数据
	res, err2 := lib.Engine.Table("plan_fence").ID(id).Update(&PlanProgress)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(res))
}

// swagger:route GET /api/plan/progress:id  progress ProgressGet
//
// 获取学生计划进度
//     Produces:
//     - application/json
//
//     Responses:
//       200: Response
func get(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)
	print("id:", id)
	planAndProgress := models.PlanAndProgress{}
	//根据id查询
	b, err := lib.Engine.Table("plan_progress").
		Join("INNER", "plan", "plan.id=plan_progress.id").
		Where("student_id=?", id).
		Cols("plan.name", "plan.date", "plan.duration", "plan.stride_frequency", "plan.pace",
			"plan_progress.distance", "plan_progress.duration", "plan_progress.times", "plan_progress.calories", "plan_progress.steps").
		Get(&planAndProgress)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该计划记录"))
		return
	}
	ctx.JSON(lib.NewResponseOK(planAndProgress))

}

// swagger:route GET /api/plan/progress progress ProgressSearch
//
// 获取学生计划进度(按字段查询 别忘了+s   /api/plan/progresss
//     Produces:
//     - application/json
//
//     Responses:
//       200: Response

func search(ctx iris.Context) {

	//创建查询Session指针
	query := lib.Engine.Table("plan_progress")

	//字段查询
	if ctx.URLParamExists("name") {
		query.And(builder.Like{"name", ctx.URLParam("name")})
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

	//查询
	var PlanProgress []models.PlanProgress
	err := query.Find(&PlanProgress)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(PlanProgress))
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

func gethistogram(ctx iris.Context) {

	var (
		gender       int
		departmentid int
		month        int
		year         int
	)

	sumChan := make(chan Sorts, 49)

	//totalfrequency := make([]int, 0)

	zerotime := zerotime()

	totalfrequencystruct := []Sorts{}

	count := 1
	counts := 1

	if ctx.URLParamExists("departmentid") {
		departmentid = int(ctx.URLParamInt64Default("departmentid", 0))
		//query :=lib.Engine.Table("plan_record").Where("departmentid_id",departmentid)
	}

	if ctx.URLParamExists("gender") {
		//gender := ctx.Params().GetUint64Default("gender",0)
		gender = int(ctx.URLParamInt64Default("gender", 0))

		//fmt.Println("性别", gender)
	}
	if ctx.URLParamExists("month") {
		month = int(ctx.URLParamInt64Default("month", 0))

	}
	if ctx.URLParamExists("year") {
		year = int(ctx.URLParamInt64Default("year", 0))

	}

	//输入到管道缓存
	for i := 0; i < 48; i++ {

		zerotime1 := zerotime
		zerotime = zerotime + 30*60
		tm1 := time.Unix(zerotime1, 0).Format("2006-01-02 15:04:05")
		tm2 := time.Unix(zerotime, 0).Format("2006-01-02 15:04:05")

		counts1 := counts
		counts++
		if ctx.URLParamExists("month") {
			now := time.Now()
			currentYear, currentMonth, _ := now.Date()
			currentLocation := now.Location()

			firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
			lastOfMonth := firstOfMonth.AddDate(0, 1, -1).Day()
			var buffer bytes.Buffer
			ftm1 := tm1[0:8]
			ltm1 := tm1[10:len(tm1)]
			buffer.WriteString(ftm1)
			buffer.WriteString("01")
			buffer.WriteString(ltm1)

			var buffer1 bytes.Buffer
			ftm2 := tm1[0:8]
			ltm2 := tm2[10:len(tm2)]
			buffer1.WriteString(ftm2)

			buffer1.WriteString(strconv.Itoa(lastOfMonth))
			buffer1.WriteString(ltm2)
			tm1 = buffer.String()
			tm2 = buffer1.String()

		}
		if ctx.URLParamExists("year") {
			var buffer bytes.Buffer
			ftm1 := tm1[0:5]
			ltm1 := tm1[10:len(tm1)]
			buffer.WriteString(ftm1)
			buffer.WriteString("01-01")
			buffer.WriteString(ltm1)

			var buffer1 bytes.Buffer
			ftm2 := tm1[0:5]
			ltm2 := tm2[10:len(tm2)]
			buffer1.WriteString(ftm2)
			buffer1.WriteString("12-31")
			buffer1.WriteString(ltm2)
			tm1 = buffer.String()
			tm2 = buffer1.String()

		}

		go func() {

			j, sum, err := findRecord(counts1, gender, departmentid, tm1, tm2, month, year)

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
			//log.Println("1111111周", value)
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

func findRecord(counts int, gender int, departmentid int, tm1, tm2 string, month, year int) (j int, sum int, err error) {
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

		//session.Exec("SELECT COALESCE(sum(`status`),0) FROM `plan_record` WHERE (DATE_FORMAT(end_time,'%H:%i:%S') > '15:00:00'  and DATE_FORMAT(end_time,'%H:%i:%S') <= '18:00:00')")  执行48次 把数据库中时间与24小时分成的48个区间比较，
		//形成每半小时的统计
		//"DATE_FORMAT(end_time,'%H:%i:%S' 目的是把数据库中的年月日时分秒序列化，并且只以时分秒的格式输出，不带年月日，便于比对48个小时区间 因为比较的是全年每一天48个时间段总的统计图
		yearsum, err2 := session.Table("plan_record").
			And("DATE_FORMAT(end_time,'%H:%i:%S') > ?", newtm1).
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

type BoyAndGirl struct {
	Boytotal             int64   `json:"Boytotal"`
	Boytotalcompletion   int64   `json:"Boytotalcompletion"`
	Boycompletiondegree  float64 `json:"Boycompletiondegree"`
	Girltotal            int64   `json:"girltotal"`
	Girltotalcompletion  int64   `json:"girltotalcompletion"`
	Girlcompletiondegree float64 `json:"girlcompletiondegree"`
}

func completiondegree(ctx iris.Context) {
	boyandgirl := BoyAndGirl{}
	//student := models.Student{}
	Progress := models.PlanProgress{}
	query := lib.Engine.Table("plan_progress")
	query2 := lib.Engine.Table("plan_progress")

	if ctx.URLParamExists("gender") {
		//gender := ctx.URLParamInt64Default("gender",0)

		boytotal, err := query.
			Join("INNER", "student", "student.id = plan_progress.student_id").
			Where("gender=?", 1).
			Cols("progress.status").Count()

		if err != nil {
			ctx.JSON(lib.NewResponseFail(1, err.Error()))
			return
		}
		boytotalcompletion, err1 := query.And("gender=?", 1).
			Join("INNER", "student", "student.id = plan_progress.student_id").
			Cols("progress.status").Sum(Progress, "status")
		if err1 != nil {
			ctx.JSON(lib.NewResponseFail(1, err1.Error()))
			return
		}

		boycompletiondegree := boytotalcompletion / float64(boytotal)

		boyandgirl.Boytotal = boytotal
		boyandgirl.Boytotalcompletion = int64(boytotalcompletion)
		boyandgirl.Boycompletiondegree = boycompletiondegree

		girltotal, err3 := query2.
			Join("INNER", "student", "student.id = plan_progress.student_id").
			Where("gender=?", 2).
			Cols("progress.status").Count()

		if err3 != nil {
			ctx.JSON(lib.NewResponseFail(1, err3.Error()))
			return

		}

		//girltotalcompletion,err4:=query.And("gender=?", 2).
		//Join("INNER","student","student.id = plan_progress.student_id").
		//Cols("progress.status").Sum(Progress,"status")
		//
		//if err4 != nil {
		//ctx.JSON(lib.NewResponseFail(1, err4.Error()))
		//return
		//

		girltotalcompletion, err4 := query2.And("gender=?", 2).
			Join("INNER", "student", "student.id = plan_progress.student_id").
			Cols("progress.status").Sum(Progress, "status")

		if err4 != nil {
			ctx.JSON(lib.NewResponseFail(1, err4.Error()))
			return

		}

		girlcompletiondegree := girltotalcompletion / float64(girltotal)
		boyandgirl.Girltotal = girltotal
		boyandgirl.Girltotalcompletion = int64(girltotalcompletion)
		boyandgirl.Girlcompletiondegree = girlcompletiondegree

		ctx.JSON(lib.NewResponseOK(boyandgirl))
		return

		//b, err := lib.Engine.Table("account").
		//	Join("INNER", "role", "role.id = account.role_id").
		//	Where("account.id=?", id).
		//	Cols("role.name", "role.privilege", "account.avatar", "role.introduce", "role.roles","account.Username").
		//	Get(&roleback)

	}

}

//func main() {
//
//	var listInt []int
//
//	rand.Seed(time.Now().UnixNano())
//
//
//	dataCh := make(chan int, 48)
//	stopCh := make(chan bool)
//
//	//发送者
//	for i := 0; i < 48; i++ {
//		zerotime1 := zerotime
//		zerotime = zerotime + 30*60
//		tm1 := time.Unix(zerotime1, 0).Format("2006-01-02 15:04:05")
//		tm2 := time.Unix(zerotime, 0).Format("2006-01-02 15:04:05")
//
//		go func() {
//			for {
//				value := rand.Intn(10)
//				select {
//				case <-stopCh:
//					return
//				case dataCh <- value:
//				}
//			}
//		}()
//	}
//
//	//接收者
//	go func() {
//
//		count := 0
//		listInt = make([]int, 0)
//		for value := range dataCh {
//			if count == 48 {
//				close(stopCh)
//				return
//			}
//			log.Println(value)
//			listInt = append(listInt, value)
//			count++
//		}
//	}()
//
//
//
//
//	log.Println("listInt:", listInt)
//}
//package main
//
//import (
//"log"
//"math/rand"
//"sync"
//"time"
//)
//
//func main() {
//
//	var listInt []int
//
//	rand.Seed(time.Now().UnixNano())
//	wg := new(sync.WaitGroup)
//	wg.Add(1)
//	dataCh := make(chan int, 48)
//	stopCh := make(chan bool)
//
//	//发送者
//	for i := 0; i < 48; i++ {
//		go func() {
//			for {
//				value := rand.Intn(10)
//				select {
//				case <-stopCh:
//					return
//				case dataCh <- value:
//				}
//			}
//		}()
//	}
//
//	//接收者
//	go func() {
//		defer wg.Done()
//		count := 0
//		listInt = make([]int, 0)
//		for value := range dataCh {
//			if count == 48 {
//				close(stopCh)
//				return
//			}
//			log.Println(value)
//			listInt = append(listInt, value)
//			count++
//		}
//	}()
//
//	wg.Wait()
//
//	log.Println("listInt:", listInt)
//}
type getprogressstruct struct {
	//每周有效跑量
	WeekDistance int `json:"week_distance" xorm:"not null comment('每周有效跑量') INT(11)"`

	//每周有效跑步次数
	WeekTimes int `json:"week_times" xorm:"not null comment('每周有效跑步次数') INT(11)"`

	//完成进度
	CompleteProgress float32 `json:"complete_progress" xorm:"comment('完成进度') FLOAT(32)"`

	//周完成进度
	WeekCompleteProgress string `json:"week_complete_progress" xorm:"comment('周完成进度') VARCHAR(200)"`

	//总里程
	TotalDistance int `json:"total_distance" xorm:"not null comment('总里程') INT(4)"`

	Name string `json:"name" xorm:"name"`

	Code string `json:"code" xorm:"code"`

	Times int `json:"times" xorm:"times"` //总跑步次数

	Distance int `json:"distance" xorm:"distance"` //总 跑步量

	Id string `json:"student_id" xorm:"id"`
}

func getprogress(ctx iris.Context) {
	var progress []getprogressstruct
	//student := models.Student{}
	type progressstruct struct {
		Res      int64               `json:"res"`
		Progress []getprogressstruct `json:"progress"`
	}

	//query := lib.Engine.Table("plan_progress")
	query := lib.Engine.Table("plan_progress")
	//搜索框
	if ctx.URLParamExists("studentCode") {
		query.And(builder.Like{"code", ctx.URLParam("studentCode")})
	}
	if ctx.URLParamExists("studentName") {
		query.And(builder.Like{"student.name", ctx.URLParam("studentName")})
	}
	if ctx.URLParamExists("planId") {
		planId, _ := strconv.Atoi(ctx.URLParam("planId"))
		if planId != -1 {
			query.And("student.plan_id=?", planId)
			query.And("plan_progress.plan_id=?", planId)
		}
	}
	if ctx.URLParamExists("planStatus") {
		query.And("plan.stop", ctx.URLParam("planStatus"))
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

	//res, err := query.Join("INNER", "student", "student.id=plan_progress.student_id").
	//	Join("INNER", "plan", "plan_progress.plan_id=plan.id").
	//	Cols("plan_progress.week_distance", "plan_progress.week_times", "plan_progress.complete_progress",
	//		"plan_progress.week_complete_progress", "plan.total_distance", "student.name", "student.code", "plan_progress.times", "plan_progress.distance").
	//	FindAndCount(&progress)
	//使用外连接替换内连接
	//TODO 历史计划也展示
	res, err := query.Join("INNER", "student", "student.id=plan_progress.student_id ").
		Join("INNER", "plan", "plan.id=plan_progress.plan_id").
		Select("plan.name as plan_name,student.id,student.code,student.name,plan_progress.week_distance,plan_progress.week_times,plan_progress.complete_progress,plan_progress.week_complete_progress, plan.total_distance, plan_progress.times, plan_progress.distance").
		And("student.plan_id!=0").
		FindAndCount(&progress)
	if err != nil {
		ctx.JSON(lib.NewResponseOK(err.Error()))
	}
	println("len list:", len(progress))

	progressreturn := progressstruct{res, progress}

	//ctx.JSON(lib.NewResponseOK(&progress))
	ctx.JSON(lib.NewResponseOK(progressreturn))
	//fmt.Println(progress)
	fmt.Println(progressreturn)
}

//计划进度图表接口
//1.前七天，2.本月，3.本学期
func progressChart(ctx iris.Context) {

	PlanId := 0
	TimeType := 1

	//获取计划id,时间段
	if ctx.URLParamExists("plan_id") {

		PlanId = int(ctx.URLParamInt64Default("plan_id", 0))
	} else {
		ctx.JSON(lib.NewResponseFail(1, "计划id错误"))
		return
	}

	if ctx.URLParamExists("time_type") {
		TimeType = int(ctx.URLParamInt64Default("time_type", 1))
	} else {
		ctx.JSON(lib.NewResponseFail(1, "时间类型错误"))
		return
	}

	if PlanId == 0 {
		ctx.JSON(lib.NewResponseFail(1, "计划id错误"))
		return
	}

	//获取计划详情
	plan := models.Plan{}
	res, err := lib.Engine.Table("plan").Where("id=?", PlanId).Get(&plan)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	if res != true {
		println("找不到计划")
		return
	}
	println("time_type:", TimeType)
	responseRecord := responseRecord{}

	switch TimeType {
	case 1:
		responseRecord = weekChartProgress(plan)
	case 2:
		responseRecord = monthChartProgress(plan)
	case 3:
		responseRecord = termChartProgress(plan)

	}
	fmt.Printf("%v", responseRecord)

	ctx.JSON(lib.SuccessResponse(responseRecord, "获取计划次数和公里数信息成功。"))

}

//前七天数据
func weekChartProgress(plan models.Plan) responseRecord {

	//获取从今天到之前七天的时间
	timeLimits := make([]string, 0)
	timeStr := time.Now().Format("2006-01-02")
	nowZero, _ := time.Parse("2006-01-02", timeStr)
	timesXUnit := make([]string, 0) //x轴单位
	for i := 6; i >= 0; i-- {
		indexZero := nowZero.AddDate(0, 0, -i)
		timeLimits = append(timeLimits, indexZero.Format("2006-01-02"))
		timesXUnit = append(timesXUnit, indexZero.Format("01/02"))
	}
	timeLimits = append(timeLimits, nowZero.AddDate(0, 0, 1).Format("2006-01-02"))
	//timesXUnit = append(timesXUnit, nowZero.AddDate(0, 0, 1).Format("01/02"))
	//for index,value:=range timeLimits{
	//	println("")
	//	fmt.Printf("%d,%v",index,value)
	//}
	//获取前七天的运动记录
	timeStartStr := timeLimits[0]
	timeEndStr := timeLimits[7]
	records := []models.PlanRecord{}
	err := lib.Engine.Table("plan_record").Where("plan_id=?", plan.Id).And("create_at BETWEEN '" + timeStartStr + "' AND '" + timeEndStr + "'").Find(&records)
	if err != nil {
		fmt.Printf("%v", err)
		return responseRecord{}
	}

	println("长度：", len(timeLimits))
	//获取每天的数据
	completeTimesArr := make([]int, 7)
	validateTimesArr := make([]int, 7)
	completeDistance := make([]int, 7)
	i3 := 0
	for i := 0; i < len(timeLimits); i++ {

		for i2 := i3; i2 < len(records); i2++ {

			if records[i2].CreateAt.Format("2006-01-02") >= timeLimits[i] && records[i2].CreateAt.Format("2006-01-02") < timeLimits[i+1] {
				completeTimesArr[i]++
				completeDistance[i] += records[i2].Distance

				if records[i2].Status == 1 {
					validateTimesArr[i]++
				}
				i3++
			} else {

				break
			}
		}

	}

	recTimes := recordTimes{
		CompleteTimes: completeTimesArr,
		ValidateTimes: validateTimesArr,
	}
	//recDistance:=recordDistance{completeDistance:completeDistance}

	weekResResponse := responseRecord{
		TimesData:    recTimes,
		DistanceData: completeDistance,
		XUnit:        timesXUnit,
	}
	fmt.Printf("%v", weekResResponse)

	return weekResResponse

}

//本月数据
func monthChartProgress(plan models.Plan) responseRecord {

	//1.获取本月几周的开始时间和结束时间
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	timeLimits := lib.CurrentMonthEveryWeekLimits(firstOfMonth, lastOfMonth, now)
	lastDay, errt := time.Parse("2006-01-02", timeLimits[len(timeLimits)-1])
	if errt != nil {
		fmt.Printf("%v", errt)
	}

	timeLimits[len(timeLimits)-1] = lastDay.AddDate(0, 0, 1).Format("2006-01-02")
	timesXUnit := make([]string, 0) //x轴单位
	for index, value := range timeLimits {
		println(index, value)
		if index == len(timeLimits)-1 {
			break
		} else {
			timeLimitsStart, _ := time.Parse("2006-01-02", value)
			timeLimitsEnd, _ := time.Parse("2006-01-02", timeLimits[index+1])
			timesXUnit = append(timesXUnit, timeLimitsStart.Format("01/02")+"-"+timeLimitsEnd.Format("01/02"))
		}

	}
	//2.获取本月的运动记录
	timeStartStr := timeLimits[0]
	timeEndStr := timeLimits[len(timeLimits)-1]
	records := []models.PlanRecord{}
	err := lib.Engine.Table("plan_record").Where("plan_id=?", plan.Id).And("create_at BETWEEN '" + timeStartStr + "' AND '" + timeEndStr + "'").Find(&records)
	if err != nil {
		fmt.Printf("%v", err)
		return responseRecord{}
	}

	//3.获取每周的数据
	completeTimesArr := make([]int, 7)
	validateTimesArr := make([]int, 7)
	completeDistance := make([]int, 7)
	i3 := 0
	for i := 0; i < len(timeLimits); i++ {

		for i2 := i3; i2 < len(records); i2++ {

			if records[i2].CreateAt.Format("2006-01-02") >= timeLimits[i] && records[i2].CreateAt.Format("2006-01-02") < timeLimits[i+1] {
				completeTimesArr[i]++
				completeDistance[i] += records[i2].Distance

				if records[i2].Status == 1 {
					validateTimesArr[i]++
				}
				i3++
			} else {

				break
			}
		}

	}

	recTimes := recordTimes{
		CompleteTimes: completeTimesArr,
		ValidateTimes: validateTimesArr,
	}
	//recDistance:=recordDistance{completeDistance:completeDistance}

	weekResResponse := responseRecord{
		TimesData:    recTimes,
		DistanceData: completeDistance,
		XUnit:        timesXUnit,
	}
	fmt.Printf("%v", weekResResponse)

	//4.整理数据

	return weekResResponse

}

//本学期数据
func termChartProgress(plan models.Plan) responseRecord {

	//获取本计划每个月的开始时间和结束时间
	timeLayout := "2006-01-02 15:04:05"
	now := plan.DateBegin
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	println(lastOfMonth.Format(timeLayout))

	timeLimits := make([]string, 0)
	timeLimits = append(timeLimits, plan.DateBegin.Format(timeLayout))

	for plan.DateEnd.Format(timeLayout) >= lastOfMonth.Format(timeLayout) {
		println(plan.DateEnd.Format(timeLayout), lastOfMonth.Format(timeLayout), "追加", lastOfMonth.Format(timeLayout))
		//timeLimits= append(timeLimits, lastOfMonth.Format(timeLayout))
		nextMonth := lastOfMonth.AddDate(0, 0, 1)
		lastOfMonth = lib.GetMonthLastDay(nextMonth)
		firstDay := lib.GetMonthFirstDay(lastOfMonth)
		timeLimits = append(timeLimits, firstDay.Format(timeLayout))
	}
	timeLimits = append(timeLimits, plan.DateEnd.Format(timeLayout))

	timesXUnit := make([]string, 0) //x轴单位
	for index, value := range timeLimits {

		timeLimits[index] = value[0:10]
		println(index, timeLimits[index])

	}

	for index, value := range timeLimits {
		if index == len(timeLimits)-1 {
			break
		} else {
			timeLimitsStart, _ := time.Parse("2006-01-02", value)
			timeLimitsEnd, _ := time.Parse("2006-01-02", timeLimits[index+1])
			timesXUnit = append(timesXUnit, timeLimitsStart.Format("01/02")+"-"+timeLimitsEnd.Format("01/02"))
		}
	}

	//获取整个计划运动记录
	timeStartStr := timeLimits[0]
	timeEndStr := timeLimits[len(timeLimits)-1]
	records := []models.PlanRecord{}
	err := lib.Engine.Table("plan_record").Where("plan_id=?", plan.Id).And("create_at BETWEEN '"+timeStartStr+"' AND '"+timeEndStr+"'").Cols("id", "status", "create_at", "distance").Find(&records)
	if err != nil {
		fmt.Printf("%v", err)
		return responseRecord{}
	}

	//获取每个月运动记录
	completeTimesArr := make([]int, 7)
	validateTimesArr := make([]int, 7)
	completeDistance := make([]int, 7)
	i3 := 0
	for i := 0; i < len(timeLimits); i++ {

		for i2 := i3; i2 < len(records); i2++ {

			if records[i2].CreateAt.Format("2006-01-02") >= timeLimits[i] && records[i2].CreateAt.Format("2006-01-02") < timeLimits[i+1] {
				completeTimesArr[i]++
				completeDistance[i] += records[i2].Distance

				if records[i2].Status == 1 {
					validateTimesArr[i]++
				}
				i3++
			} else {

				break
			}
		}

	}

	recTimes := recordTimes{
		CompleteTimes: completeTimesArr,
		ValidateTimes: validateTimesArr,
	}
	//recDistance:=recordDistance{completeDistance:completeDistance}

	weekResResponse := responseRecord{
		TimesData:    recTimes,
		DistanceData: completeDistance,
		XUnit:        timesXUnit,
	}
	fmt.Printf("%v", weekResResponse)
	//整理数据
	return weekResResponse

}
