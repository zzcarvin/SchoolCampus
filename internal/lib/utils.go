package lib

import (
	"Campus/internal/models"
	"fmt"
	"log"
	"math"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type weekContent struct {
	Sequence int     `json:"sequence"`
	Distance int     `json:"distance"`
	Count    int     `json:"count"`
	Duration int     `json:"_"`
	Steps    int     `json:"_"`
	Pace     float32 `json:"pace"`
	Status   int     `json:"status"`
}

//对外使用的
type SimpleWeekContent struct {
	Sequence int `json:"sequence"`
	Distance int `json:"distance"`
	Times    int `json:"times"`
	Status   int `json:"status"`
}

//生成随机数字串
func RandomNumber(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(10) + 0x30 //'0'
		bytes[i] = byte(b)
	}
	return string(bytes)
	//return fmt.Sprintf("%0"+ strconv.Itoa(len)+"v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

//生成随机数字串
func RandomUpperString(len int) string {
	var str = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(36)
		bytes[i] = byte(str[b])
	}
	return string(bytes)
}

//生成随机数字串
func RandomLowerString(len int) string {
	var str = []byte("abcdefghijklmnopqrstuvwxyz0123456789")
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(36)
		bytes[i] = byte(str[b])
	}
	return string(bytes)
}

//检查手机号合法性
func IsValidPhone(phone string) bool {
	reg := `^[1](([3][0-9])|([4][5-9])|([5][0-3,5-9])|([6][5,6])|([7][0-8])|([8][0-9])|([9][189]))[0-9]{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(phone)
}

//步频转换成字符串
func FomPace(fl64Frequency float64) string {
	flminute := int(fl64Frequency / 60)
	flsecond := fl64Frequency - float64(flminute*60)
	println(flsecond)
	minute := strconv.Itoa(flminute)
	second := strconv.FormatFloat(flsecond, 'f', 0, 64)
	strFrequency := minute + "'" + second + "''"

	return strFrequency

}

func FormFloat(value float64) (float64, error) {
	flvalue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	println("form float:")
	fmt.Printf("%f", flvalue)
	println("-----------")
	return flvalue, nil

	//pow10_n := math.Pow10(num)
	//return math.Trunc((value+0.5/pow10_n)*pow10_n) / pow10_n
}

func FloatRoundingToFloat(fl float64, num int) float64 {
	str := FloatRounding(fl, num)
	return Str2float64(str)
}

//保留小数点后几位
func FloatRounding(fl float64, num int) string {
	var t float64
	var data string

	if num < 0 {
		num = -num
	}
	if fl == 0 {
		return "0"
	}
	f := math.Pow10(num)
	fl += 0.0000000001
	x := fl * f
	if x >= 0.0 {
		t = math.Ceil(x)
		if (t - x) > 0.5 {
			t -= 1.0
		}
	} else {
		t = math.Ceil(-x)
		if (t + x) > 0.5 {
			t -= 1.0
		}
		t = -t
	}
	x = t / f
	data = strconv.FormatFloat(x, 'f', num, 64)
	if num > 0 {
		data = strings.TrimRight(data, "0")
		data = strings.TrimRight(data, ".")
	}
	return data
}

func Str2float64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

//格式化时间，00:00:00
func ConvertTime(seconds int) string {

	hour := seconds / 3600
	minute := (seconds - hour*3600) / 60
	second := seconds - hour*3600 - minute*60
	strHour := ""
	strMinute := ""
	strSecond := ""
	if hour < 10 {
		strHour = fmt.Sprintf("0%d", hour)
	} else {
		strHour = fmt.Sprintf("%d", hour)
	}
	if minute < 10 {
		strMinute = fmt.Sprintf("0%d", minute)
	} else {
		strMinute = fmt.Sprintf("%d", minute)
	}
	if second < 10 {
		strSecond = fmt.Sprintf("0%d", second)
	} else {
		strSecond = fmt.Sprintf("%d", second)
	}

	convertTime := strHour + ":" + strMinute + ":" + strSecond
	return convertTime
}

func UnixToFormTime(timeStamp int64) string {

	t := time.Unix(int64(timeStamp), 0)
	//返回string
	dateStr := t.Format("2006/01/02 15:04:05")
	return dateStr
}

func Int2str(t int) string {
	str := strconv.Itoa(t)
	return str
}

func Str2int(str string) int {
	t, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return t
}

func Int642str(t int64) string {
	str := strconv.FormatInt(int64(t), 10)
	return str
}

func GetReferenceValue(t reflect.Value, referenceColumn string) interface{} {
	var referenceValue interface{}
	if t.FieldByName(referenceColumn).Type().String() == "int" {
		referenceValue = Int642str(t.FieldByName(referenceColumn).Int())
	} else {
		referenceValue = t.FieldByName(referenceColumn)
	}
	return referenceValue
}

func GetFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

//根据开始时间和结束时间获取中间的秒数
func GetSecond(startTime time.Time, endTime time.Time) int64 {

	var duration int64
	duration = endTime.Unix() - startTime.Unix()
	return duration
}

//获取当前是计划的第几周
func GetWeekSequence(planId int, studentId int, nowTime time.Time) int {

	weekSequence := 0

	var plan models.Plan
	res, err := Engine.Table("plan").Where("id=?", planId).Get(&plan)
	if err != nil {
		println("")
		fmt.Printf("%v", err)
		return 0
	}
	if res == false {
		return 0
	}

	//获取学生信息
	student := models.Student{}
	resStudent, err := Engine.Table("student").Where("id=?", studentId).Get(&student)
	if err != nil {
		fmt.Printf("%v", err)
		return 0
	}
	if resStudent == false {
		return 0
	}
	//timeLayout := "2006-01-02 15:04:05"
	planBegin := plan.DateBegin.Unix()
	//第一周的开始时间和结束时间（时间戳） 结束时间：%v,firstWeekEndDay.Format("2006-01-02 15:04:05"),
	weekNumber := int(plan.DateBegin.Weekday())
	if weekNumber == 0 { //周日，转成7
		weekNumber = 7
	}
	//firstWeekEndDay := planBegin + int64(((7-int(plan.DateBegin.Weekday()))+1)*3600*24)
	firstWeekEndDay := planBegin + int64(7-weekNumber+1)*3600*24
	if plan.DateBegin.Weekday() == 0 { //周日，是0
		firstWeekEndDay = planBegin + 3600*24
	}

	println("计划开始时周几：", plan.DateBegin.Weekday(), int(plan.DateBegin.Weekday()))
	fmt.Printf("第一周开始时间:%v，时间戳：%v,结束时间戳：%v", plan.DateBegin, planBegin, firstWeekEndDay)

	//最后一周的开始时间和结束时间（时间戳）
	planEnd := plan.DateEnd.Unix()
	lastWeekDays := int(plan.DateEnd.Weekday()) - 1
	//if int(plan.DateEnd.Weekday()) == 0 {
	//	lastWeekDays = 6
	//}
	if int(plan.DateEnd.Weekday()) == 0 {
		lastWeekDays = 7
	}
	lastWeekEndDay := planEnd - int64(lastWeekDays*3600*24)
	println("")
	fmt.Printf("最后一周开始时间戳：%v,最后一周结束时间戳:%v", lastWeekEndDay, planEnd)

	//每周的结束时间

	//总周数
	Weeks := (lastWeekEndDay-firstWeekEndDay)/604800 + 2
	fmt.Printf("取整 int:%d", (lastWeekEndDay-firstWeekEndDay)/604800)
	floWeeks := float32(lastWeekEndDay-firstWeekEndDay) / 604800
	fmt.Printf("floweeks: %f", floWeeks)
	if floWeeks <= 1 {
		Weeks = 1
	}
	fmt.Printf("中间周数:%v,取余：%v", Weeks, (lastWeekEndDay-firstWeekEndDay)%604800)

	weekUnixs := make([]int, 0)
	weekUnixs = append(weekUnixs, int(planBegin))

	//分成三种情况，一周，两周，三周和三周以上
	//一周：firstWeekEndDay==lastWeekEndDay

	//两周:secondWeekEndDay==lastWeekEndDay
	println("")
	//planWeeks := 0
	println("第一周结束时间戳：", firstWeekEndDay, "form:", UnixToFormTime(firstWeekEndDay), "最后一周结束时间戳：", planEnd, "form:", UnixToFormTime(planEnd))
	if firstWeekEndDay >= planEnd { //一周
		weekUnixs = append(weekUnixs, int(planEnd))
		//planWeeks = 1
	} else if (firstWeekEndDay + 604800) >= planEnd { //二周
		weekUnixs = append(weekUnixs, int(firstWeekEndDay))
		weekUnixs = append(weekUnixs, int(planEnd))
		//planWeeks = 2
	} else { //三周和三周以上
		weekUnixs = append(weekUnixs, int(firstWeekEndDay))
		//planWeeks = 3
		//println("开始,长度：", len(weekUnixs))
		//for index, value := range weekUnixs {
		//	println("index:", index, "value:", value, time.Unix(int64(value), 0).Format(timeLayout))
		//}

		//println("---------")
		//
		//println("中间开始第一周第一天：", time.Unix(int64(firstWeekEndDay), 0).Format(timeLayout), "最后一周最后一天", time.Unix(int64(lastWeekEndDay), 0).Format(timeLayout), "周数：", int((lastWeekEndDay-firstWeekEndDay)/604800))

		lastWeekStartDay := lastWeekEndDay + 1
		for i := 0; i < int((lastWeekStartDay-firstWeekEndDay)/604800); i++ {
			weekUnixs = append(weekUnixs, int(firstWeekEndDay)+604800*(i+1))
		}
		//println("---------")
		//println("中间,长度：", len(weekUnixs))
		//for index, value := range weekUnixs {
		//	println("index:", index, "value:", value, value, time.Unix(int64(value), 0).Format(timeLayout))
		//}

		//println("---------")
		weekUnixs = append(weekUnixs, int(planEnd))
		//println("最后,长度：", len(weekUnixs))

		//for index, value := range weekUnixs {
		//	println("index:", index, "value:", value, value, time.Unix(int64(value), 0).Format(timeLayout))
		//}
	}

	//println("本计划的周数：", planWeeks)
	//for index, value := range weekUnixs {
	//	println("index:", index, "value:", value, value, time.Unix(int64(value), 0).Format(timeLayout))
	//}

	for i := 0; i < len(weekUnixs); i++ {
		if int(nowTime.Unix()) >= weekUnixs[i] && int(nowTime.Unix()) <= weekUnixs[i+1] {
			weekSequence = i + 1
		}
	}
	return weekSequence
}

//获取当天零时时间戳
func ZeroTimeToUnix(presentTime time.Time) int64 {
	timeStr := presentTime.Format("2006-01-02")
	fmt.Println("timeStr:", timeStr)
	t, _ := time.Parse("2006-01-02", timeStr)
	timeNumber := t.Unix()
	return timeNumber
}

/**
获取本周周一的日期
*/
//func GetFirstDateOfWeek(present time.Time) (weekMonday string) {
//	now := present
//
//	offset := int(time.Monday - now.Weekday())
//	if offset > 0 {
//		offset = -6
//	}
//
//	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
//	weekMonday = weekStartDate.Format("2006-01-02")
//	return
//}
func GetFirstDateOfWeek(present time.Time) (weekMonday time.Time) {
	now := present

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	//weekMonday = weekStartDate.Format("2006-01-02")
	return weekStartDate
}

//获取一个学生的某个计划每周的完成状态
func GetWeeksRecords(studentId int, plan models.Plan) [][]models.PlanRecord {
	planId := plan.Id
	var nilReocrds [][]models.PlanRecord

	//获取学生该计划的所有运动记录
	var allRecords []models.PlanRecord
	allRecordErr := Engine.Table("plan_record").Where("status=1").And("student_id=?", studentId).And("plan_id=?", planId).Find(&allRecords)
	if allRecordErr != nil {
		fmt.Printf("查询所有记录错误：%v", allRecordErr)
		return nilReocrds
	}

	//println("打印所有运动记录：**********************")
	//for index, value := range allRecords {
	//	println("")
	//	fmt.Printf("index%d:%v", index, value)
	//}

	//timeLayout := "2006-01-02 15:04:05"
	planBegin := plan.DateBegin.Unix()
	//第一周的开始时间和结束时间（时间戳） 结束时间：%v,firstWeekEndDay.Format("2006-01-02 15:04:05"),
	weekNumber := int(plan.DateBegin.Weekday())
	if weekNumber == 0 { //周日，转成7
		weekNumber = 7
	}
	//firstWeekEndDay := planBegin + int64(((7-int(plan.DateBegin.Weekday()))+1)*3600*24)
	firstWeekEndDay := planBegin + int64(7-weekNumber+1)*3600*24
	if plan.DateBegin.Weekday() == 0 { //周日，是0
		firstWeekEndDay = planBegin + 3600*24
	}

	println("计划开始时周几：", plan.DateBegin.Weekday(), int(plan.DateBegin.Weekday()))
	fmt.Printf("第一周开始时间:%v，时间戳：%v,结束时间戳：%v", plan.DateBegin, planBegin, firstWeekEndDay)

	//最后一周的开始时间和结束时间（时间戳）
	planEnd := plan.DateEnd.Unix()
	lastWeekDays := int(plan.DateEnd.Weekday()) - 1
	//if int(plan.DateEnd.Weekday()) == 0 {
	//	lastWeekDays = 6
	//}
	if int(plan.DateEnd.Weekday()) == 0 {
		lastWeekDays = 7
	}
	lastWeekEndDay := planEnd - int64(lastWeekDays*3600*24)
	println("")
	fmt.Printf("最后一周开始时间戳：%v,最后一周结束时间戳:%v", lastWeekEndDay, planEnd)

	//每周的结束时间

	//总周数
	Weeks := (lastWeekEndDay-firstWeekEndDay)/604800 + 2
	fmt.Printf("取整 int:%d", (lastWeekEndDay-firstWeekEndDay)/604800)
	floWeeks := float32(lastWeekEndDay-firstWeekEndDay) / 604800
	fmt.Printf("floweeks: %f", floWeeks)
	if floWeeks <= 1 {
		Weeks = 1
	}
	fmt.Printf("中间周数:%v,取余：%v", Weeks, (lastWeekEndDay-firstWeekEndDay)%604800)

	weekUnixs := make([]int, 0)
	weekUnixs = append(weekUnixs, int(planBegin))

	//分成三种情况，一周，两周，三周和三周以上
	//一周：firstWeekEndDay==lastWeekEndDay

	//两周:secondWeekEndDay==lastWeekEndDay
	println("")
	//planWeeks := 0
	println("第一周结束时间戳：", firstWeekEndDay, "form:", UnixToFormTime(firstWeekEndDay), "最后一周结束时间戳：", planEnd, "form:", UnixToFormTime(planEnd))
	if firstWeekEndDay >= planEnd { //一周
		weekUnixs = append(weekUnixs, int(planEnd))
		//planWeeks = 1
	} else if (firstWeekEndDay + 604800) >= planEnd { //二周
		weekUnixs = append(weekUnixs, int(firstWeekEndDay))
		weekUnixs = append(weekUnixs, int(planEnd))
		//planWeeks = 2
	} else { //三周和三周以上
		weekUnixs = append(weekUnixs, int(firstWeekEndDay))
		//planWeeks = 3
		//println("开始,长度：", len(weekUnixs))
		//for index, value := range weekUnixs {
		//	println("index:", index, "value:", value, time.Unix(int64(value), 0).Format(timeLayout))
		//}

		//println("---------")

		//println("中间开始第一周第一天：", time.Unix(int64(firstWeekEndDay), 0).Format(timeLayout), "最后一周最后一天", time.Unix(int64(lastWeekEndDay), 0).Format(timeLayout), "周数：", int((lastWeekEndDay-firstWeekEndDay)/604800))

		lastWeekStartDay := lastWeekEndDay + 1
		for i := 0; i < int((lastWeekStartDay-firstWeekEndDay)/604800); i++ {
			weekUnixs = append(weekUnixs, int(firstWeekEndDay)+604800*(i+1))
		}
		//println("---------")
		//println("中间,长度：", len(weekUnixs))
		//for index, value := range weekUnixs {
		//	println("index:", index, "value:", value, value, time.Unix(int64(value), 0).Format(timeLayout))
		//}

		//println("---------")
		weekUnixs = append(weekUnixs, int(planEnd))
		//println("最后,长度：", len(weekUnixs))

		//for index, value := range weekUnixs {
		//	println("index:", index, "value:", value, value, time.Unix(int64(value), 0).Format(timeLayout))
		//}
	}

	//println("本计划的周数：", planWeeks)
	//for index, value := range weekUnixs {
	//	println("index:", index, "value:", value, value, time.Unix(int64(value), 0).Format(timeLayout))
	//}

	//获取每周的运动记录----开始
	weeksRecord := make([][]models.PlanRecord, len(weekUnixs))
	//for index, value := range weekUnixs {
	//	println("index:", index, "value:", value, value, time.Unix(int64(value), 0).Format(timeLayout))
	//
	//}

	weekContents := make([]weekContent, 0)
	i3 := 0
	for i := 0; i < len(weekUnixs); i++ {

		newWeekContent := weekContent{}

		println("")
		//fmt.Printf("本周的开始时间戳%v,结束时间戳：%v",time.Unix(int64(weekUnixs[i]), 0).Format(timeLayout),time.Unix(int64(weekUnixs[i+1]), 0).Format(timeLayout))
		for i2 := i3; i2 < len(allRecords); i2++ {
			//println("i2:", i2, allRecords[i2].CreateAt.Unix(), "weekUnix i:", i, weekUnixs[i])
			if int(allRecords[i2].CreateAt.Unix()) >= weekUnixs[i] && int(allRecords[i2].CreateAt.Unix()) <= weekUnixs[i+1] {
				weeksRecord[i] = append(weeksRecord[i], allRecords[i2])
				i3++
				//println("当前record index:", i3)
			} else {
				break
			}
		}

		if newWeekContent.Steps != 0 && newWeekContent.Duration != 0 {
			newWeekContent.Pace = float32(newWeekContent.Steps) / (float32(newWeekContent.Duration) / 60) //步频,修改精度问题
		}
		weekContents = append(weekContents, newWeekContent)

		//println("")
		//fmt.Printf("周：%d,记录数量：%d", i, len(weeksRecord[i]))

	}
	//获取每周的运动记录----结束

	//返回所有周的运动记录

	return weeksRecord
}

//获取本月每周的开始时间和结束时间
//获取当前时间在时间段内的第几周
//思路：首先用根据第一周的结束时间，判断时间段一共多少周。
// 一周，两周直接累加。三周以上将时间段分成所有三段，第一周，中间所有周，最后一周。并转成时间戳。（用于比较来判断第几周）
func CurrentMonthEveryWeekLimits(startTime time.Time, endTime time.Time, nowTime time.Time) []string {

	planBegin := startTime.Unix()
	//第一周的开始时间和结束时间（时间戳） 结束时间：%v,firstWeekEndDay.Format("2006-01-02 15:04:05"),
	weekNumber := int(startTime.Weekday())
	if weekNumber == 0 { //周日，转成7
		weekNumber = 7
	}

	firstWeekEndDay := planBegin + int64(7-weekNumber+1)*3600*24
	if startTime.Weekday() == 0 { //周日，是0
		firstWeekEndDay = planBegin + 3600*24
	}

	//最后一周的开始时间和结束时间（时间戳）
	planEnd := endTime.Unix()
	lastWeekDays := int(endTime.Weekday()) - 1
	if int(endTime.Weekday()) == 0 {
		lastWeekDays = 7
	}
	lastWeekEndDay := planEnd - int64(lastWeekDays*3600*24)
	//每周的结束时间

	//总周数
	Weeks := (lastWeekEndDay-firstWeekEndDay)/604800 + 2
	floWeeks := float32(lastWeekEndDay-firstWeekEndDay) / 604800
	if floWeeks <= 1 {
		Weeks = 1
	}
	println(Weeks)
	weekUnixs := make([]int, 0)
	weekUnixs = append(weekUnixs, int(planBegin))

	//分成三种情况，一周，两周，三周和三周以上
	if firstWeekEndDay >= planEnd { //一周
		weekUnixs = append(weekUnixs, int(planEnd))
		//planWeeks = 1
	} else if (firstWeekEndDay + 604800) >= planEnd { //二周
		weekUnixs = append(weekUnixs, int(firstWeekEndDay))
		weekUnixs = append(weekUnixs, int(planEnd))
		//planWeeks = 2
	} else { //三周和三周以上
		weekUnixs = append(weekUnixs, int(firstWeekEndDay))
		lastWeekStartDay := lastWeekEndDay + 1
		for i := 0; i < int((lastWeekStartDay-firstWeekEndDay)/604800); i++ {
			weekUnixs = append(weekUnixs, int(firstWeekEndDay)+604800*(i+1))
		}
		weekUnixs = append(weekUnixs, int(planEnd))
	}

	//转化成时间字符串
	weeksLimitsStr := make([]string, len(weekUnixs))
	for i := 0; i < len(weekUnixs); i++ {

		weeksLimitsStr[i] = time.Unix(int64(weekUnixs[i]), 0).Format("2006-01-02")
		fmt.Println(weeksLimitsStr[i]) //2018-07-11 15:10:19

	}

	return weeksLimitsStr
}

//获取某个月的最后一天
func GetMonthLastDay(current time.Time) time.Time {

	currentYear, currentMonth, _ := current.Date()
	currentLocation := current.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return lastOfMonth
}

//获取某个月第一天
func GetMonthFirstDay(current time.Time) time.Time {
	currentYear, currentMonth, _ := current.Date()
	currentLocation := current.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)

	return firstOfMonth
}
