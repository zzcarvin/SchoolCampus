package lib

import (
	"Campus/internal/models"
	"fmt"
	"github.com/robfig/cron"
	"strconv"
	"time"
)

type depRecLis struct {
	Id         int     `json:"id" xorm:'id'`        //Id
	Name       string  `json:"name" xorm:"name"`    //姓名
	Percentage float32 `json:"percentage" xorm:"-"` //百分比
	//Classes    []int               `json:"classes" xorm:"-"`     //班级id
	Times    int `json:"times" xorm:"-"`    //完成次数
	Distance int `json:"distance" xorm:"-"` //完成公里数
	//每个单位的数据
	Progresses []progressData `json:"progresses" xorm:"-"`
}

type progressData struct {
	PlanId     int     `json:"plan_id"`
	ClassId    int     `json:"class_id"`
	Times      int     `json:"times"`
	Distance   int     `json:"distance"`
	Percentage float32 `json:"percentage"`
}

type weekData struct {
	WeekTimes    int `json:"week_times" xorm:"week_times"`
	WeekDistance int `json:"week_distance" xorm:"week_distance"`
}

//定时任务
func NewCron() {
	cron.New()

	c := cron.New()
	//每周运行，周日到周一之间早上凌晨零点，minutes, hours, day of month, month, day of week
	//second, minute, hour, day of month, month, day of week
	//秒，分钟，小时，本月第几天，月份，本周第几天
	//每周一零时执行
	//specEveryWeek := "0 0 0 ? * MON"
	//每天零时零分运行，这段cron 现在的执行时间是：每天23点的第59分钟，第59秒
	specEveryDay := "59 59 23 * * ?"

	//每天获取小程序list数据
	err := c.AddFunc(specEveryDay, cronWeekProgress)
	if err != nil {
		fmt.Printf("每天更新progress week错误:%v", err)
	}
	errTerm := c.AddFunc(specEveryDay, cronTermProgress)
	if errTerm != nil {
		fmt.Printf("每天更新progress term错误:%v", errTerm)
		return
	}

	////每周定时清除学生计划进度上周数据
	//errClear := c.AddFunc(specEveryWeek, ClearWeekProgress)
	//if errClear != nil {
	//	fmt.Printf("每周清除学生每周计划进度数据错误：%v", errClear)
	//	return
	//}

	c.Start()
	//defer c.Stop()

}

//定时查询本周，本学期所有院系的跑步次数和有效公里数
//1.获取所有院系
//2.遍历院系，获取每个院系的所有计划。
//2.1.在每个院系内遍历所有计划，查询计划进度表，获取该计划的该院系的所有每周次数之和，和每周进度之和,百分比。
func cronWeekProgress() {
	println("每天凌晨零点零分执行，获取每周数据：")
	//获取所有院系
	departmentUnion := make([]depRecLis, 0)
	err := Engine.Table("department").Find(&departmentUnion)
	if err != nil {
		fmt.Printf("查询院系错误：%v", err)
		return
	}

	fmt.Printf("所有的院系：%v", departmentUnion)
	println("")

	//获取当前学年和学期
	termYear := 0
	term := 0
	yearPlans := make([]models.Plan, 0)
	errPlans := Engine.Table("plan").Desc("create_time").Find(&yearPlans)
	if errPlans != nil {
		fmt.Printf("查询该系所有计划详细信息出错：%v", errPlans)
		return
	}
	if len(yearPlans) != 0 {
		termYear = yearPlans[0].TermYear
		term = yearPlans[0].Term
	}

	//遍历所有院系，获取一个院系的所有计划，获取院系的所有计划的次数和公里数
	for index, value := range departmentUnion {

		//获取该系的所有正在运行的计划
		plans := make([]models.Plan, 0)
		planIds := make([]models.Student, 0)
		err := Engine.Table("student").Where("department_id=?", value.Id).And("plan_id!=?", 0).GroupBy("plan_id").Find(&planIds)
		if err != nil {
			fmt.Printf("查询该系所有计划出错：%v", err)
			return
		}

		//必须给progress添加否则会报错
		if len(planIds) == 0 {

		} else {
			planIdsStr := "id=" + strconv.Itoa(planIds[0].PlanId)

			for i := 1; i < len(planIds); i++ {
				planIdsStr = planIdsStr + " OR " + "id=" + strconv.Itoa(planIds[i].PlanId)
			}
			println("planIdsStr:", planIdsStr)
			errPlans := Engine.Table("plan").Where(planIdsStr).Find(&plans)
			if errPlans != nil {
				fmt.Printf("查询该系所有计划详细信息出错：%v", errPlans)
				return
			}
			//有计划时获取学年和学期

			println(value.Name + "所有正在运行的计划：")
			fmt.Printf("%v", plans)
			println("")

			//遍历该院系计划，初始化计划的数据
			for _, valuuep2 := range plans {
				departmentUnion[index].Progresses = append(departmentUnion[index].Progresses, progressData{PlanId: valuuep2.Id})
			}

			//获取每周次数和每周公里数
			for indexPlans, valuePlans := range plans {
				recordTimes := weekData{}
				resData, err := Engine.Table("plan_progress").Where("department_id=?", value.Id).And("plan_id=?", valuePlans.Id).SumsInt(&recordTimes, "week_times", "week_distance")
				if err != nil {
					fmt.Printf("查询计划进度错误：%v", err)
					return
				}

				fmt.Printf("recordTimes:%v", resData)
				//计算每周运动次数，公里数和完成度
				departmentUnion[index].Progresses[indexPlans].Times = int(resData[0])
				departmentUnion[index].Progresses[indexPlans].Distance = int(resData[1])
				stuNum := depStuNum(departmentUnion[index].Id)
				if stuNum != 0 {
					departmentUnion[index].Progresses[indexPlans].Percentage = float32(departmentUnion[index].Progresses[indexPlans].Times) / float32(valuePlans.MinWeekTimes*stuNum)
				}

				//累加次数和公里数
				departmentUnion[index].Times += departmentUnion[index].Progresses[indexPlans].Times
				departmentUnion[index].Distance += departmentUnion[index].Progresses[indexPlans].Distance
				//计算百分比
				if indexPlans == len(plans)-1 {
					var sumPercentage float32
					for _, valuePercen := range departmentUnion[index].Progresses {
						sumPercentage += valuePercen.Percentage
					}
					departmentUnion[index].Percentage = sumPercentage / float32(len(departmentUnion[index].Progresses))
				}

			}

		}

		//插入或更新院系计划进度表。需要添加字段学年和学期，不然不能定位一个院系在哪一个时间的记录。应该需要一个表来存储年纪和学期。
		//1.先查询该系，该学期有没有记录，有更新，没有插入
		//获取第几周
		_, sequenceWeek := time.Now().ISOWeek()

		//更新或插入院系每周数据---开始

		res, err := Engine.Table("week_progress").
			Where("department_id=?", value.Id).And("term_year=?", termYear).And("term=?", term).And("sequence=?", sequenceWeek).
			Exist()
		if err != nil {
			fmt.Printf("查询week_progress出错：%v", err)
			return
		}

		if res != true {
			//插入
			thisWeekProgress := models.WeekProgress{
				DepartmentId: value.Id,
				Times:        departmentUnion[index].Times,
				Distance:     departmentUnion[index].Distance,
				Percentage:   departmentUnion[index].Percentage,
				Sequence:     sequenceWeek,
				TermYear:     termYear,
				Tear:         term,
			}
			affected, err := Engine.Table("week_progress").Insert(thisWeekProgress)
			if err != nil {
				fmt.Printf("插入错误：%v", err)
				return
			}
			if affected == 0 {
				println("插入失败")
				return
			}
		} else {
			//更新 TODO
			thisWeekProgress := models.WeekProgress{
				Times:      departmentUnion[index].Times,
				Distance:   departmentUnion[index].Distance,
				Percentage: departmentUnion[index].Percentage,
			}
			_, err := Engine.Table("week_progress").Cols("times", "distance", "percentage").
				Where("department_id=?", value.Id).And("term_year=?", termYear).And("term=?", term).And("sequence=?", sequenceWeek).
				Update(&thisWeekProgress)
			if err != nil {
				fmt.Printf("更新错误：%v", err)
				return
			}
		}

		//更新或插入院系每周数据---结束

		//更新或插入所有班级每周数据---开始
		classProgress(value.Id, termYear, term)
		//更新或插入所有班级每周数据---结束

	}

	fmt.Printf("获取每周院系所有计划数据：%v", departmentUnion)

}

//获取该系所有学生的数量
func depStuNum(departmentId int) int {

	students := make([]models.Student, 0)
	err := Engine.Table("student").Where("department_id=?", departmentId).Find(&students)
	if err != nil {
		fmt.Printf("查询该系所有学生出错：%v", err)
		return 0
	}
	return len(students)
}

func classProgress(departmentId int, termYear int, term int) {
	//获取一个系所有班
	departmentUnion := make([]depRecLis, 0)
	err := Engine.Table("classes").Where("department_id=?", departmentId).Find(&departmentUnion)
	if err != nil {
		fmt.Printf("查询班级错误：%v", err)
		return
	}

	fmt.Printf("所有的班级：%v", departmentUnion)
	for index, value := range departmentUnion {

		//获取该班的所有正在运行的计划
		plans := make([]models.Plan, 0)
		planIds := make([]models.Student, 0)
		err := Engine.Table("student").Where("class_id=?", value.Id).And("plan_id!=?", 0).GroupBy("plan_id").Find(&planIds)
		if err != nil {
			fmt.Printf("查询该系所有计划出错：%v", err)
			return
		}

		//必须给progress添加否则会报错
		if len(planIds) == 0 {

		} else {
			planIdsStr := "id=" + strconv.Itoa(planIds[0].PlanId)

			for i := 1; i < len(planIds); i++ {
				planIdsStr = planIdsStr + " OR " + "id=" + strconv.Itoa(planIds[i].PlanId)
			}
			println("planIdsStr:", planIdsStr)
			errPlans := Engine.Table("plan").Where(planIdsStr).Find(&plans)
			if errPlans != nil {
				fmt.Printf("查询该系所有计划详细信息出错：%v", errPlans)
				return
			}
			//有计划时获取学年和学期

			println(value.Name + "所有正在运行的计划：")
			fmt.Printf("%v", plans)
			println("")

			//遍历该院系计划，初始化计划的数据
			for _, valuuep2 := range plans {
				departmentUnion[index].Progresses = append(departmentUnion[index].Progresses, progressData{PlanId: valuuep2.Id})
			}

			//获取每周次数和每周公里数
			for indexPlans, valuePlans := range plans {
				recordTimes := weekData{}
				resData, err := Engine.Table("plan_progress").Where("class_id=?", value.Id).And("plan_id=?", valuePlans.Id).SumsInt(&recordTimes, "week_times", "week_distance")
				if err != nil {
					fmt.Printf("查询计划进度错误：%v", err)
					return
				}

				fmt.Printf("recordTimes:%v", resData)
				//计算每周运动次数，公里数和完成度
				departmentUnion[index].Progresses[indexPlans].Times = int(resData[0])
				departmentUnion[index].Progresses[indexPlans].Distance = int(resData[1])
				stuNum := depStuNum(departmentUnion[index].Id)
				if stuNum != 0 {
					departmentUnion[index].Progresses[indexPlans].Percentage = float32(departmentUnion[index].Progresses[indexPlans].Times) / float32(valuePlans.MinWeekTimes*stuNum)
				}

				//累加次数和公里数
				departmentUnion[index].Times += departmentUnion[index].Progresses[indexPlans].Times
				departmentUnion[index].Distance += departmentUnion[index].Progresses[indexPlans].Distance
				//计算百分比
				if indexPlans == len(plans)-1 {
					var sumPercentage float32
					for _, valuePercen := range departmentUnion[index].Progresses {
						sumPercentage += valuePercen.Percentage
					}
					departmentUnion[index].Percentage = sumPercentage / float32(len(departmentUnion[index].Progresses))
				}

			}

		}

		//插入或更新院系计划进度表。需要添加字段学年和学期，不然不能定位一个院系在哪一个时间的记录。应该需要一个表来存储年纪和学期。
		//1.先查询该系，该学期有没有记录，有更新，没有插入
		//获取第几周
		_, sequenceWeek := time.Now().ISOWeek()

		//更新或插入院系每周数据---开始

		res, err := Engine.Table("week_progress").
			Where("department_id=?", departmentId).And("class_id=?", value.Id).And("term_year=?", termYear).And("term=?", term).And("sequence=?", sequenceWeek).
			Exist()
		if err != nil {
			fmt.Printf("查询week_progress出错：%v", err)
			return
		}

		if res != true {
			//插入
			thisWeekProgress := models.WeekProgress{
				DepartmentId: departmentId,
				ClassId:      value.Id,
				Times:        departmentUnion[index].Times,
				Distance:     departmentUnion[index].Distance,
				Percentage:   departmentUnion[index].Percentage,
				Sequence:     sequenceWeek,
				TermYear:     termYear,
				Tear:         term,
			}
			affected, err := Engine.Table("week_progress").Insert(thisWeekProgress)
			if err != nil {
				fmt.Printf("插入错误：%v", err)
				return
			}
			if affected == 0 {
				println("插入失败")
				return
			}
		} else {
			//更新 TODO
			thisWeekProgress := models.WeekProgress{
				Times:      departmentUnion[index].Times,
				Distance:   departmentUnion[index].Distance,
				Percentage: departmentUnion[index].Percentage,
			}
			_, err := Engine.Table("week_progress").Cols("times", "distance", "percentage").
				Where("department_id=?", departmentId).And("class_id=?", value.Id).And("term_year=?", termYear).And("term=?", term).And("sequence=?", sequenceWeek).
				Update(&thisWeekProgress)
			if err != nil {
				fmt.Printf("更新错误：%v", err)
				return
			}
		}

		//更新或插入院系每周数据---结束

	}

}

//每学期全系定时任务
func cronTermProgress() {
	println("每天凌晨零点零一分执行，获取每学期的数据：")
	//获取所有院系
	departmentUnion := make([]depRecLis, 0)
	err := Engine.Table("department").Find(&departmentUnion)
	if err != nil {
		fmt.Printf("查询院系错误：%v", err)
		return
	}

	fmt.Printf("所有的院系：%v", departmentUnion)
	println("")

	//获取当前学年和学期
	termYear := 0
	term := 0
	yearPlans := make([]models.Plan, 0)
	errPlans := Engine.Table("plan").Desc("create_time").Find(&yearPlans)
	if errPlans != nil {
		fmt.Printf("查询该系所有计划详细信息出错：%v", errPlans)
		return
	}
	if len(yearPlans) != 0 {
		termYear = yearPlans[0].TermYear
		term = yearPlans[0].Term
	}

	//遍历所有院系，获取一个院系的所有计划，获取院系的所有计划的次数和公里数
	for index, value := range departmentUnion {

		//获取该系的所有正在运行的计划
		plans := make([]models.Plan, 0)
		planIds := make([]models.Student, 0)
		err := Engine.Table("student").Where("department_id=?", value.Id).And("plan_id!=?", 0).GroupBy("plan_id").Find(&planIds)
		if err != nil {
			fmt.Printf("查询该系所有计划出错：%v", err)
			return
		}

		//必须给progress添加否则会报错
		if len(planIds) == 0 {

		} else {
			planIdsStr := "id=" + strconv.Itoa(planIds[0].PlanId)

			for i := 1; i < len(planIds); i++ {
				planIdsStr = planIdsStr + " OR " + "id=" + strconv.Itoa(planIds[i].PlanId)
			}
			println("planIdsStr:", planIdsStr)
			errPlans := Engine.Table("plan").Where(planIdsStr).Find(&plans)
			if errPlans != nil {
				fmt.Printf("查询该系所有计划详细信息出错：%v", errPlans)
				return
			}
			//有计划时获取学年和学期

			println(value.Name + "所有正在运行的计划：")
			fmt.Printf("%v", plans)
			println("")

			//遍历该院系计划，初始化计划的数据
			for _, valuuep2 := range plans {
				departmentUnion[index].Progresses = append(departmentUnion[index].Progresses, progressData{PlanId: valuuep2.Id})
			}

			//获取每周次数和每周公里数
			for indexPlans, valuePlans := range plans {
				recordTimes := weekData{}
				resData, err := Engine.Table("plan_progress").Where("department_id=?", value.Id).And("plan_id=?", valuePlans.Id).SumsInt(&recordTimes, "times", "distance")
				if err != nil {
					fmt.Printf("查询计划进度错误：%v", err)
					return
				}

				fmt.Printf("recordTimes:%v", resData)
				//计算每周运动次数，公里数和完成度
				departmentUnion[index].Progresses[indexPlans].Times = int(resData[0])
				departmentUnion[index].Progresses[indexPlans].Distance = int(resData[1])
				stuNum := depStuNum(departmentUnion[index].Id)
				if stuNum != 0 {
					departmentUnion[index].Progresses[indexPlans].Percentage = float32(departmentUnion[index].Progresses[indexPlans].Times) / float32(valuePlans.MinWeekTimes*stuNum)
				}

				//累加次数和公里数
				departmentUnion[index].Times += departmentUnion[index].Progresses[indexPlans].Times
				departmentUnion[index].Distance += departmentUnion[index].Progresses[indexPlans].Distance
				//计算百分比
				if indexPlans == len(plans)-1 {
					var sumPercentage float32
					for _, valuePercen := range departmentUnion[index].Progresses {
						sumPercentage += valuePercen.Percentage
					}
					departmentUnion[index].Percentage = sumPercentage / float32(len(departmentUnion[index].Progresses))
				}

			}

		}

		//插入或更新院系计划进度表。需要添加字段学年和学期，不然不能定位一个院系在哪一个时间的记录。应该需要一个表来存储年纪和学期。
		//1.先查询该系，该学期有没有记录，有更新，没有插入
		//更新或插入院系每学期数据---开始

		res, err := Engine.Table("term_progress").
			Where("department_id=?", value.Id).And("term_year=?", termYear).And("term=?", term).
			Exist()
		if err != nil {
			fmt.Printf("查询week_progress出错：%v", err)
			return
		}

		if res != true {
			thisWeekProgress := models.TermProgress{
				DepartmentId: value.Id,
				Times:        departmentUnion[index].Times,
				Distance:     departmentUnion[index].Distance,
				Percentage:   departmentUnion[index].Percentage,
				TermYear:     termYear,
				Tear:         term,
			}
			affected, err := Engine.Table("term_progress").Insert(thisWeekProgress)
			if err != nil {
				fmt.Printf("插入错误：%v", err)
				return
			}
			if affected == 0 {
				println("插入失败")
				return
			}
		} else {
			//更新
			thisWeekProgress := models.TermProgress{
				Times:      departmentUnion[index].Times,
				Distance:   departmentUnion[index].Distance,
				Percentage: departmentUnion[index].Percentage,
			}
			_, err := Engine.Table("term_progress").Cols("times", "distance", "percentage").
				Where("department_id=?", value.Id).And("term_year=?", termYear).And("term=?", term).
				Update(&thisWeekProgress)
			if err != nil {
				fmt.Printf("更新错误：%v", err)
				return
			}
		}

		//更新或插入院系每学期数据---结束

		//更新或插入所有班级每学期数据---开始
		classTermProgress(value.Id, termYear, term)
		//更新或插入所有班级每学期数据---结束

	}

	fmt.Printf("获取每周院系所有计划数据：%v", departmentUnion)

}

//每学期一个系全班
func classTermProgress(departmentId int, termYear int, term int) {
	//获取一个系所有班
	departmentUnion := make([]depRecLis, 0)
	err := Engine.Table("classes").Where("department_id=?", departmentId).Find(&departmentUnion)
	if err != nil {
		fmt.Printf("查询班级错误：%v", err)
		return
	}

	fmt.Printf("所有的班级：%v", departmentUnion)
	for index, value := range departmentUnion {

		//获取该班的所有正在运行的计划
		plans := make([]models.Plan, 0)
		planIds := make([]models.Student, 0)
		err := Engine.Table("student").Where("class_id=?", value.Id).And("plan_id!=?", 0).GroupBy("plan_id").Find(&planIds)
		if err != nil {
			fmt.Printf("查询该系所有计划出错：%v", err)
			return
		}

		//必须给progress添加否则会报错
		if len(planIds) == 0 {

		} else {
			planIdsStr := "id=" + strconv.Itoa(planIds[0].PlanId)

			for i := 1; i < len(planIds); i++ {
				planIdsStr = planIdsStr + " OR " + "id=" + strconv.Itoa(planIds[i].PlanId)
			}
			println("planIdsStr:", planIdsStr)
			errPlans := Engine.Table("plan").Where(planIdsStr).Find(&plans)
			if errPlans != nil {
				fmt.Printf("查询该系所有计划详细信息出错：%v", errPlans)
				return
			}
			//有计划时获取学年和学期

			println(value.Name + "所有正在运行的计划：")
			fmt.Printf("%v", plans)
			println("")

			//遍历该院系计划，初始化计划的数据
			for _, valuuep2 := range plans {
				departmentUnion[index].Progresses = append(departmentUnion[index].Progresses, progressData{PlanId: valuuep2.Id})
			}

			//获取每周次数和每周公里数
			for indexPlans, valuePlans := range plans {
				recordTimes := weekData{}
				resData, err := Engine.Table("plan_progress").Where("class_id=?", value.Id).And("plan_id=?", valuePlans.Id).SumsInt(&recordTimes, "times", "distance")
				if err != nil {
					fmt.Printf("查询计划进度错误：%v", err)
					return
				}

				fmt.Printf("recordTimes:%v", resData)
				//计算每周运动次数，公里数和完成度
				departmentUnion[index].Progresses[indexPlans].Times = int(resData[0])
				departmentUnion[index].Progresses[indexPlans].Distance = int(resData[1])
				stuNum := depStuNum(departmentUnion[index].Id)
				if stuNum != 0 {
					departmentUnion[index].Progresses[indexPlans].Percentage = float32(departmentUnion[index].Progresses[indexPlans].Times) / float32(valuePlans.MinWeekTimes*stuNum)
				}

				//累加次数和公里数
				departmentUnion[index].Times += departmentUnion[index].Progresses[indexPlans].Times
				departmentUnion[index].Distance += departmentUnion[index].Progresses[indexPlans].Distance
				//计算百分比
				if indexPlans == len(plans)-1 {
					var sumPercentage float32
					for _, valuePercen := range departmentUnion[index].Progresses {
						sumPercentage += valuePercen.Percentage
					}
					departmentUnion[index].Percentage = sumPercentage / float32(len(departmentUnion[index].Progresses))
				}

			}

		}

		//插入或更新院系计划进度表。需要添加字段学年和学期，不然不能定位一个院系在哪一个时间的记录。应该需要一个表来存储年纪和学期。
		//1.先查询该系，该学期有没有记录，有更新，没有插入

		//更新或插入院系每学期数据---开始

		res, err := Engine.Table("term_progress").
			Where("department_id=?", departmentId).And("class_id=?", value.Id).And("term_year=?", termYear).And("term=?", term).
			Exist()
		if err != nil {
			fmt.Printf("查询week_progress出错：%v", err)
			return
		}

		if res != true {
			//插入 TODO
			thisWeekProgress := models.TermProgress{
				DepartmentId: departmentId,
				ClassId:      value.Id,
				Times:        departmentUnion[index].Times,
				Distance:     departmentUnion[index].Distance,
				Percentage:   departmentUnion[index].Percentage,
				TermYear:     termYear,
				Tear:         term,
			}
			affected, err := Engine.Table("term_progress").Insert(thisWeekProgress)
			if err != nil {
				fmt.Printf("插入错误：%v", err)
				return
			}
			if affected == 0 {
				println("插入失败")
				return
			}
		} else {
			//更新 TODO
			thisWeekProgress := models.TermProgress{
				Times:      departmentUnion[index].Times,
				Distance:   departmentUnion[index].Distance,
				Percentage: departmentUnion[index].Percentage,
			}
			_, err := Engine.Table("term_progress").Cols("times", "distance", "percentage").
				Where("department_id=?", departmentId).And("class_id=?", value.Id).And("term_year=?", termYear).And("term=?", term).
				Update(&thisWeekProgress)
			if err != nil {
				fmt.Printf("更新错误：%v", err)
				return
			}
		}

		//更新或插入院系每学期数据---结束

	}

}

//清除所有学生上周数据
//func ClearWeekProgress() {
//	println("清空学生上周跑步进度：")
//	//周数据清空
//	zeroWeek := models.PlanProgress{
//		WeekTimes:            0,
//		WeekDistance:         0,
//		WeekCompleteProgress: "",
//	}
//
//	res, err := Engine.Table("plan_progress").Cols("week_distance", "week_times", "week_complet_progress").Update(&zeroWeek)
//	if err != nil {
//		fmt.Printf("清空学生每周数据错误：%v", err)
//		return
//	}
//	println("res:", res)
//
//}
