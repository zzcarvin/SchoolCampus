package share

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"fmt"
	"github.com/kataras/iris"
	"time"
)

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

type newResponseShare struct {
	Record       models.PlanRecord `json:"record"`
	SumDistance  int               `json:"sumDistance"`
	ContinueDays int               `json:"continue_days"`
}

//分享
func share(ctx iris.Context) {

	query := lib.Engine.Table("plan_record")
	//var  studentId ,planId int
	//var errUrl error
	if ctx.URLParamExists("studentId") {
		query.And("student_id=?", ctx.URLParam("studentId"))
	}
	if ctx.URLParamExists("studentId") {
		query.And("plan_id=?", ctx.URLParam("planId"))
	}
	//计划跑第几天
	//学生id,计划id，
	//1、获取本次计划运动记录，排序 最长里程，最长时间，最小配速，最多卡路里。中第一个运动记录
	var shareRecord []responseShare
	err := query.Join("INNER", "student", "student.id=plan_record.student_id").
		Desc("distance").
		Desc("duration").
		Desc("calories").
		Asc("pace").
		Asc("create_at").
		Cols("student.name", "plan_record.distance", "plan_record.duration", "plan_record.calories", "plan_record.pace", "plan_record.create_at").
		Find(&shareRecord)
	if err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(lib.NewResponseFail(0, err.Error()))
		return
	}
	if shareRecord == nil {
		ctx.JSON(lib.NewResponseFail(0, "查到相关记录"))
		return
	}

	//获取总距离
	for index, _ := range shareRecord {
		shareRecord[0].AllDistance = shareRecord[0].AllDistance + shareRecord[index].LongestDistance
		println("")
		fmt.Printf("%v", shareRecord[index])
	}
	//获取天数
	now := time.Now().Format("2006-01-02 15:04:05")
	nowUnix := time.Now().Unix()
	first := shareRecord[0].CreateAt.Unix()
	shareRecord[0].Sequence = int((nowUnix - first) / (24 * 60 * 60))

	//获取现在
	shareRecord[0].Now = now

	ctx.JSON(lib.NewResponseOK(shareRecord[0]))

	//用今天的时间戳减第一次运动记录的时间戳，再除以24小时

}

func newShare(ctx iris.Context) {
	//获取student_id
	//id := ctx.Params().GetUint64Default("id", 0)
	//studentId:=0
	//planId:=0
	var studendId, planId string
	if ctx.URLParamExists("studentId") {
		//query.And("student_id=?", ctx.URLParam("studentId"))
		studendId = ctx.URLParam("studentId")
	}
	if ctx.URLParamExists("planId") {
		//query.And("plan_id=?", ctx.URLParam("planId"))
		planId = ctx.URLParam("planId")
	}

	print("student_id:", studendId, "planId:", planId)

	PlanRecord := models.PlanRecord{}
	sumDistance, err := lib.Engine.Table("plan_record").Where("student_id=?", studendId).And("status=?", 1).Sum(&PlanRecord, "distance")
	println("")
	fmt.Printf("sumDistance: %f", sumDistance)

	//重新排序寻找需要的最佳记录----开始

	//单次最长里程
	distanceRecord := models.PlanRecord{}
	bDistance, err := lib.Engine.Table("plan_record").Where("student_id=?", studendId).And("status=?", 1).Desc("distance").Get(&distanceRecord)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if bDistance == false {
		ctx.JSON(lib.NewResponseFail(0, "未找到最佳公里数运动记录"))
		return
	}
	//单次最长时间
	durationRecord := models.PlanRecord{}
	bduration, err := lib.Engine.Table("plan_record").Where("student_id=?", studendId).And("status=?", 1).Desc("duration").Get(&durationRecord)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if bduration == false {
		ctx.JSON(lib.NewResponseFail(0, "未找到最佳运动时长运动记录"))
		return
	}

	//单次最佳配速
	paceRecord := models.PlanRecord{}
	bpace, err := lib.Engine.Table("plan_record").Where("student_id=?", studendId).And("status=?", 1).Desc("calories").Get(&paceRecord)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if bpace == false {
		ctx.JSON(lib.NewResponseFail(0, "未找到最佳配速运动记录"))
		return
	}

	//单次最多消耗卡路里
	caloryRecord := models.PlanRecord{}
	bcalory, err := lib.Engine.Table("plan_record").Where("student_id=?", studendId).And("status=?", 1).Asc("pace").Get(&caloryRecord)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if bcalory == false {
		ctx.JSON(lib.NewResponseFail(0, "未找到最佳卡路里运动记录"))
		return
	}

	//重新排序寻找需要的最佳记录----结束

	//获取连续运动天数

	student := models.Student{}
	bstudent, err := lib.Engine.Table("student").Where("id=?", studendId).Get(&student)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, "查询学生错误"))
		return
	}
	if bstudent == false {
		ctx.JSON(lib.NewResponseFail(1, "查询学生失败"))
		return
	}

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	bestShare := responseShare{
		Name:            student.Name,
		Sequence:        student.Continue,
		AllDistance:     int(sumDistance),
		LongestDistance: distanceRecord.Distance,
		LongestDuration: durationRecord.Duration,
		MaxCalories:     float32(caloryRecord.Calories),
		FastestPace:     paceRecord.Pace,
		CreateAt:        time.Time{},
		Now:             now,
	}

	//返回距离最远的相关记录记录
	ctx.JSON(lib.NewResponseOK(bestShare))
}
