package student

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"github.com/go-xorm/builder"
	"github.com/kataras/iris"
	"strings"
)

type studentRequest struct {
	Name string `json:"name"`
	Code int    `json:"code"`
}

type responseRecord struct {
	Id           int             `json:"id" xorm:"autoincr id"`
	SchoolId     int             `json:"school_id" xorm:"school_id"`
	PlanId       int             `json:"plan_id" xorm:"plan_id"`
	Name         string          `json:"name" xorm:"name"`
	StudentId    int             `json:"student_id" xorm:"student_id"`
	Type         int             `json:"type" xorm:"type"`
	StartTime    string          `json:"start_time" xorm:"start_time"`
	EndTime      string          `json:"end_time" xorm:"end_time"`
	Distance     int             `json:"distance" xorm:"distance"`
	Duration     int             `json:"duration" xorm:"duration"`
	Calories     float64         `json:"calories" xorm:"calories"`
	Steps        int             `json:"steps" xorm:"steps"`
	Pace         int             `json:"pace" xorm:"pace"`
	FormPace     string          `json:"form_pace" xorm:"form_pace"`
	Points       []models.Points `json:"points" xorm:"points"`
	Frequency    float64         `json:"frequency" xorm:"frequency"`
	Frequencies  []int           `json:"frequencies" xorm:"frequencies"`
	CreateAt     string          `json:"create_at" xorm:"create_at created"`
	XFrequencies []int           `json:"x_frequencies" xorm:"-"`
	XNumber      []int           `json:"x_number" xorm:"-"`
	Status       int             `json:"status" xorm:"status"`
}

//学生基本信息 ，学号，姓名，班级，年级，院系

//学生所有跑步记录
func records(ctx iris.Context) {
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

//运动数据统计
func recordsStatistics(ctx iris.Context) {

}

//每周计划完成情况（所有周）
func everyWeek(ctx iris.Context) {

}
