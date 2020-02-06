package feedback

import (
	"Campus/internal/controllers/v3/app/run"
	"Campus/internal/lib"
	"Campus/internal/models"
	"fmt"
	"github.com/go-xorm/builder"
	"github.com/kataras/iris"
	"strings"
)

type requestCreate struct {
	StudentId int    `json:"student_id"`
	Content   string `json:"content"`
}

func create(ctx iris.Context) {

	feedback := models.Feedback{}

	//解析
	err := ctx.ReadJSON(&feedback)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	//if feedback.Content == "" {
	//	ctx.JSON(lib.NewResponseFail(1, "您只有一次申诉机会"))
	//	return
	//
	//}
	unSql := `SELECT id FROM feedback where record_id = ?`
	res, unerr := lib.Engine.QueryString(unSql, feedback.RecordId)
	if unerr != nil {
		println(unerr.Error())

		return
	}
	if len(res) != 0 {
		ctx.JSON(lib.NewResponseFail(1, "您只有一次申诉机会"))
		return

	}
	fmt.Println("打印\n\n\n\n\n\n", res)

	//res,unerr :=db.DB.Query(unSql,feedback.RecordId)
	//if unerr != nil {
	//	return
	//}
	//if len(res) == 0 {
	//	ctx.JSON(lib.NewResponseFail(1, "您只有一次申诉机会"))
	//
	//
	//
	//}

	//b, errid := lib.Engine.Exist(feedback.RecordId)
	//if errid != nil {
	//	ctx.JSON(lib.NewResponseFail(1, err.Error()))
	//
	//	return
	//}
	//if b == true {
	//	ctx.JSON(lib.NewResponseFail(1, "您只有一次申诉机会"))
	//	println("您只有一次申诉机会")
	//	return
	//}

	//插入数据
	affected, err := lib.Engine.Table("feedback").Insert(&feedback)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if affected == 0 {
		ctx.JSON(lib.NewResponseFail(1, "反馈添加失败"))
		return
	}
	ctx.JSON(lib.NewResponseOK(feedback))

}

//更新字段内容
func update(ctx iris.Context) {
	id := ctx.Params().GetUint64Default("id", 0)
	feedfack := models.Feedback{}
	err := ctx.ReadJSON(&feedfack)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	//TODO 验证数据有效性

	//插入数据
	res, err1 := lib.Engine.Table("feedback").ID(id).Update(feedfack)
	if err1 != nil {
		ctx.JSON(lib.NewResponseFail(1, err1.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(res))
}

//feedback的status状态修改
func feedbackStatus(ctx iris.Context) {
	id := ctx.Params().GetUint64Default("id", 0)
	feedfack := models.Feedback{}
	err := ctx.ReadJSON(&feedfack)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	//TODO 验证数据有效性

	//插入数据
	res, err1 := lib.Engine.Table("feedback").MustCols("status").ID(id).Update(feedfack)
	if err1 != nil {
		ctx.JSON(lib.NewResponseFail(1, err1.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(res))
}

//本函数适用于专门更新record表的status状态的，目前的反馈模块的逻辑是，app上传跑步记录，根据跑步记录的数据来判断跑步记录的状态，分为0，1
//1为正常，当record表中的status判断为0时，此时进入反馈模块，会识别status为0的record记录显示为异常，此时app端也有此条异常记录显示，此时学生可以申诉此记录，申诉内容
//可以与异常记录同时显示在一行中，然后学校后端页面可以对异常记录根据反馈信息进行审核，当审核通过时，直接更改record的status即可。函数如下
func statusupdate(ctx iris.Context) {
	id := ctx.Params().GetUint64Default("id", 0)
	planrecord := models.PlanRecord{}
	err := ctx.ReadJSON(&planrecord)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

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

	//TODO 验证数据有效性
	fmt.Println(planrecord)
	//插入数据
	res, err1 := lib.Engine.Table("plan_record").MustCols("status").ID(id).Update(planrecord)
	if err1 != nil {
		ctx.JSON(lib.NewResponseFail(1, err1.Error()))
		return
	}

	//增加或减少该学生的运动计划进度

	println("旧运动记录的状态：", oldRecord.Status, "要修改的状态：", planrecord.Status)
	if planrecord.Status == oldRecord.Status {
		ctx.JSON(lib.NewResponseFail(1, "状态未发生变化，不做修改"))
		return
	}

	//先获取该运动进度
	progress := models.PlanProgress{}

	respro, err := lib.Engine.Table("plan_progress").Where("plan_id=?", oldRecord.PlanId).And("student_id=?", oldRecord.StudentId).Get(&progress)
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

	plan := models.Plan{}
	planRes, err := lib.Engine.Table("plan").Where("id=?", oldRecord.PlanId).Get(&plan)
	if err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(1, "查询计划错误"))
		return
	}
	if planRes == false {
		println("查询计划失败")
		ctx.JSON(lib.NewResponseFail(1, "查询计划失败"))
		return
	}

	//更新运动计划进度
	planProgress := models.PlanProgress{}
	//TODO complete没有计算
	//无效改有效
	record1 := models.PlanRecord{}

	//查询计划
	b, err := lib.Engine.Table("plan_record").Where("id=?", oldRecord.Id).Get(&record1)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, "查询跑步记录错误"))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "查询跑步记录失败"))
		return
	}

	if planrecord.Status == 1 {
		run.UpdateProgressJob(record1, plan)
	}

	println("修改后的计划进度里程：", planProgress.Distance)

	ctx.JSON(lib.NewResponseOK(res))
}
func get(ctx iris.Context) {
	////取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	feedback := models.Feedback{}
	//根据id查询
	b, err := lib.Engine.Table("feedback").ID(id).Get(&feedback)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该条异常"))
		return
	}
	fmt.Println(feedback)
	ctx.JSON(feedback)
}

func search(ctx iris.Context) {

	//创建查询Session指针
	query := lib.Engine.Table("feedback")

	//字段查询
	if ctx.URLParamExists("record_id") {
		query.And(builder.Like{"record_id", ctx.URLParam("record_id")})
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
	var feedback []models.Feedback
	err := query.Find(&feedback)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(feedback))
}

//用户申诉更新进度函数
func UpdateProgress(recordId int) (models.PlanProgress, error) {

	c := lib.GetRedisConn()
	defer c.Close()
	fmt.Println("获取到数据，开始计算", recordId)

	//修改多计划开始***********************************************************************************
	record1 := models.PlanRecord{}
	student1 := models.Student{}
	planProgress1 := models.PlanProgress{}
	plan1 := models.Plan{}

	b, err := lib.Engine.Table("plan_record").Where("id=?", recordId).Get(&record1)
	if err != nil {
		fmt.Printf("%v", err)
		fmt.Println("\n\n\n 查询跑步记录错误 \n\n\n")
		//ctx.JSON(lib.NewResponseFail(1, "查询跑步记录错误"))
		return models.PlanProgress{}, err
	}
	if b == false {
		println("查询跑步记录失败")
		//ctx.JSON(lib.NewResponseFail(1, "查询跑步记录失败"))
		return models.PlanProgress{}, err
	}
	println("获取的运动记录：")
	fmt.Printf("record:%v", record1)

	b, err = lib.Engine.Table("plan").Where("id=?", record1.PlanId).Get(&plan1)
	if err != nil {
		fmt.Printf("%v", err)
		fmt.Println("\n\n\n 查询计划错误 \n\n\n")
		//ctx.JSON(lib.NewResponseFail(1, "查询计划错误"))
		return models.PlanProgress{}, err
	}
	if b == false {
		println("查询计划失败")
		//ctx.JSON(lib.NewResponseFail(1, "查询计划失败"))
		return models.PlanProgress{}, err
	}

	progress := models.PlanProgress{}

	respro, err := lib.Engine.Table("plan_progress").Where("plan_id=?", record1.PlanId).And("student_id=?", record1.StudentId).Get(&planProgress1)
	if err != nil {
		fmt.Printf("%v", err)
		return models.PlanProgress{}, err
	}
	if respro == false {
		println("学生计划进度查询失败")
		return models.PlanProgress{}, err
	}

	println("原先计划进度里程：", progress.Distance)

	progressFlg := true
	//当天数据
	dayRecords := []models.PlanRecord{}
	errRecords := lib.Engine.Table("plan_record").Where("status=?", 1).And("student_id=?", record1.StudentId).
		And("plan_id=?", record1.PlanId).And("TO_DAYS(create_at) = TO_DAYS(NOW())").Find(&dayRecords)
	if errRecords != nil {
		fmt.Printf("获取运动记录错误：%v", errRecords)
		return models.PlanProgress{}, err
	}
	//超过日跑次数，不累计进度，每日最高次数为0未考虑！！！！
	if len(dayRecords) > plan1.MaxDayTimes {
		progressFlg = false
		return models.PlanProgress{}, err
	}

	//未满足当日最高次数，允许累计
	if progressFlg {

		//本周跑步数据
		statusanddistance, err := lib.Engine.Table("plan_record").Where("status=1").And("student_id=?", student1.Id).
			And("YEARWEEK( DATE_FORMAT(  `plan_record`.`create_at`, '%Y-%m-%d' ),1 ) = YEARWEEK( NOW(),1 )").SumsInt(record1, "status", "distance")
		if err != nil {
			fmt.Printf("查询周记录错误：%v", err)
			fmt.Println("\n\n\n查询周记录错误\n\n\n")
			//ctx.JSON(lib.NewResponseFail(1, err.Error()))
			return models.PlanProgress{}, err
		}
		println("\n\n\n\nstatusanddistance\n\n\n\n\n")
		println(statusanddistance)

		finishDistance := record1.Distance
		//先比较单次最高里程，超过单次最高里程，本次里程以单次最高里程为准
		if plan1.MaxSingleDistance != 0 && record1.Distance > plan1.MaxSingleDistance {
			finishDistance = plan1.MaxSingleDistance
		}

		totalTimes := planProgress1.Times + 1
		totalDistance := planProgress1.Distance + finishDistance
		weekTimes := planProgress1.WeekTimes + 1
		weekDistance := planProgress1.Distance + finishDistance

		//超过周跑最高次数，本周跑步次数为每周最高次数，周里程累加，总次数不累加，总里程累加
		if plan1.MaxWeekTimes != 0 && weekTimes > plan1.MaxWeekTimes {
			println("超过本周最大次数")
			totalTimes = planProgress1.Times
		}

		//超过本周最高里程，总里程只累加到单周最高里程，本周里程为最高里程
		if plan1.MaxWeekDistance != 0 && weekDistance > plan1.MaxWeekDistance {
			println("超过本周最大里程")
			//总进度重新计算：总里程=总里程-本周总里程+每周最高里程
			totalDistance = totalDistance - weekDistance + plan1.MaxWeekDistance
		}
		println("record calory:", record1.Calories, "progress calory:", planProgress1.Calories)
		//小于计划要求，更新计划进度
		planProgress := models.PlanProgress{
			Distance:     totalDistance,
			Duration:     record1.Duration + planProgress1.Duration,
			Calories:     record1.Calories + planProgress1.Calories,
			Steps:        record1.Steps + planProgress1.Steps,
			Times:        totalTimes,
			WeekTimes:    weekTimes,
			WeekDistance: weekDistance,
		}

		fmt.Printf("获取的更新数据：%v", planProgress)

		return planProgress, nil

	} else {
		return models.PlanProgress{}, err
	}

}
