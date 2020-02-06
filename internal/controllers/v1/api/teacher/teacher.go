package teacher

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/go-xorm/builder"
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type TeacherInfo struct {
	ID           string
	realName     string
	gender       int
	departmentId int
}

type DepartmentInfo struct {
	Id   int
	Name string
	//departmentType int
}

type Column struct {
	Name string
	Tag  string
}

type UpdateModel struct {
	Id             int    `json:"code"`
	Name           string `json:"name"`
	Gender         string `json:"gender"`
	DepartmentName string `json:"department_name"`
	Class          string `json:"class"`
}

type UpdateTeacher struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Gender       string `json:"gender"`
	DepartmentId string `json:"department_id"`
	//Class string `json:"-"`
}

type updTeaMap struct {
	Id   int
	Code int
}

type InsertTeacherClass struct {
	TeacherId int `json:"teacher_id" xorm:"teacher_id"`
	ClassId   int `json:"class_id" xorm:"class_id"`
}

func create(ctx iris.Context) {
	teacher := models.Teacher{}

	//解析department
	err := ctx.ReadJSON(&teacher)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	//插入数据
	res, err := lib.Engine.Table("teacher").Insert(&teacher)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	lib.NewResponseOK(res)
	ctx.JSON(lib.NewResponseOK(teacher))

}

func remove(ctx iris.Context) {
	id := ctx.Params().GetUint64Default("id", 0)
	teacher := models.Teacher{}
	session := lib.Engine.NewSession()
	defer session.Close()

	err1 := session.Begin()
	if err1 != nil {
		ctx.JSON(lib.NewResponseFail(1, "事务开启失败"))
		println("事务开启失败")
		return
	}
	affected, err2 := session.Table("teacher").ID(id).Delete(&teacher)
	if err2 != nil {
		session.Rollback()
		ctx.JSON(lib.NewResponseFail(1, "删除老师失败"))
		return
	}
	err4 := session.Commit()
	if err4 != nil {
		panic(err4.Error())
	}
	ctx.JSON(lib.NewResponseOK(affected))
}

func update(ctx iris.Context) {
	id := ctx.Params().GetUint64Default("id", 0)

	teacher := models.Teacher{}

	//解析department
	err := ctx.ReadJSON(&teacher)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//TODO 验证数据有效性

	//插入数据
	_, err2 := lib.Engine.Table("teacher").ID(id).Update(teacher)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	ctx.JSON(lib.NewResponseOK(teacher))
}

func get(ctx iris.Context) {

	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	teacher := models.TeacherAllInfos{}
	//根据id查询
	b, err := lib.Engine.Table("teacher").
		Join("INNER", "department", "department.id=teacher.department_id").
		Where("teacher.id=?", id).
		Cols("teacher.id", "teacher.code", "teacher.name", "department.name", "teacher.department_id", "teacher.gender", "teacher.create_at", "teacher.cellphone").
		Get(&teacher)
	// SELECT `student`.`id`, `student`.`code`, `student`.`name`, `classes`.`name`, `department`.`name`, `student`.`gender`, `student`.`create_at`, `student`.`cellphone`
	// FROM `student` INNER JOIN classes ON classes.id=student.class_id INNER JOIN department ON department.id=student.department_id WHERE (student.id=?) LIMIT 1 []interface {}{0x1}

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该用户"))
		return
	}

	ctx.JSON(lib.NewResponseOK(teacher))
}

func search(ctx iris.Context) {

	//创建查询Session指针
	query := lib.Engine.Table("teacher")

	//字段查询
	if ctx.URLParamExists("name") {
		query.And(builder.Like{"name", ctx.URLParam("name")})
	}
	if ctx.URLParamExists("departmentid") {
		query.And(builder.Like{"department_id", ctx.URLParam("departmentid")})
	}
	if ctx.URLParamExists("code") {
		query.And(builder.Like{"code", ctx.URLParam("code")})
	}
	if ctx.URLParamExists(("userid")) {
		query.And(builder.Like{"user_id", ctx.URLParam("userid")})
	}
	if ctx.URLParamExists("cellphone") {
		query.And(builder.Like{"cellphone", ctx.URLParam("cellphone")})
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
	var dataList []models.TeacherAllInfos
	counts, err := query.Join("INNER", "department", "department.id=teacher.department_id").
		Cols("teacher.id", "teacher.code", "teacher.name", "department.name", "teacher.department_id", "teacher.gender", "teacher.create_at", "teacher.cellphone").
		FindAndCount(&dataList)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	var pageModel = models.Page{
		List:  dataList,
		Total: counts,
	}

	ctx.JSON(lib.NewResponseOK(pageModel))
}

func excel(ctx iris.Context) {
	// Get the file from the request.
	_, info, err := ctx.FormFile("file")
	ext := path.Ext(info.Filename)
	fileName := uuid.NewV4().String() + ext
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(lib.NewResponseFail(0, "文件上传失败"))
		return
	}
	lib.SaveUploadedFile(info, "./excel", fileName)
	errfile := os.Remove(fileName) //删除文件test.txt
	if errfile != nil {
		//如果删除失败则输出 file remove Error!
		fmt.Println("file remove Error!")
		//输出错误详细信息
		fmt.Printf("%s", err)
	}

	f, err := excelize.OpenFile("./excel/" + fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	departments := make(map[string]int)
	//查院系
	var departmentList []DepartmentInfo
	lib.Engine.Table("department").Cols("id", "name").Find(&departmentList)
	if len(departmentList) != 0 {
		for _, value := range departmentList {
			departments[value.Name] = value.Id
			println("")
			fmt.Printf("院系id:%d,院系名称：%v", value.Id, departments[value.Name])
		}
	}

	updateNum := 0
	createAt := time.Now().Format("2006-01-02 15:04:05")
	sqlStr := "insert into teacher(code,name,gender,department_id,create_at) values"
	rows := f.GetRows("Sheet1")

	IdList := make(map[string]int)
	rowIndex := make(map[string]int) //列名对应的索引，这样导入的时候不用担心列名顺序的问题
	for _, row := range rows {
		updateNum++
		if updateNum == 1 {
			//第一行根据列名导入
			var columns = []string{"教工号", "姓名", "性别", "学院"}
			rowIndex = lib.GetRowIndex(row, columns)
			continue
		}
		//教工号重复直接返回
		if IdList[row[rowIndex["教工号"]]] == 1 {
			fmt.Println("教工号：" + row[rowIndex["教工号"]])
			ctx.JSON(lib.NewResponseFail(1, "教工号重复："+row[rowIndex["教工号"]]))
			return
		}
		var info = TeacherInfo{
			ID:       row[rowIndex["教工号"]],
			realName: row[rowIndex["姓名"]],
		}
		if row[rowIndex["性别"]] == "男" {
			info.gender = 1
		} else {
			info.gender = 2
		}
		info.departmentId = departments[row[rowIndex["学院"]]]
		IdList[info.ID] = 1
		code, _ := strconv.Atoi(info.ID)
		sqlStr += fmt.Sprintf("(%d,'%s',%d,%d,'%s'),", code, info.realName, info.gender, info.departmentId, createAt)
	}

	//TODO: 如果和库里比对，有重复的教工号要提示
	sqlStr = strings.TrimRight(sqlStr, ",")
	fmt.Printf("添加教工sql %s", sqlStr)
	num := strconv.Itoa(updateNum - 1)
	str := "导入教工成功,共添加" + num + "条教工"
	fmt.Println(str)
	_, err = lib.Engine.Query(sqlStr)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(lib.NewResponseFail(1, "导入失败"))
		return
	}
	ctx.JSON(lib.NewResponseOK(str))
	return
}

//批量更新体育老师信息，包括修改老师信息，修改老师的体育班信息
//获取excel，更新teacher表，删除旧的teacher_class表信息，添加新的teacher_class信息
func manyUpdateExcel(ctx iris.Context) {

	//1.获取excel文件内容
	_, info, err := ctx.FormFile("file")
	ext := path.Ext(info.Filename)
	fileName := uuid.NewV4().String() + ext
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(lib.NewResponseFail(0, "文件上传失败"))
		return
	}
	lib.SaveUploadedFile(info, "./excel", fileName)
	errfile := os.Remove(fileName) //删除文件test.txt
	if errfile != nil {
		//如果删除失败则输出 file remove Error!
		fmt.Println("file remove Error!")
		//输出错误详细信息
		fmt.Printf("%s", err)
	}

	f, err := excelize.OpenFile("./excel/" + fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	departments := make(map[string]int)
	//查院系
	var departmentList []DepartmentInfo
	lib.Engine.Table("department").Cols("id", "name").Find(&departmentList)
	if len(departmentList) != 0 {
		for _, value := range departmentList {
			departments[value.Name] = value.Id
			println("")
			fmt.Printf("院系id:%d,院系名称：%v", value.Id, departments[value.Name])
		}
	}

	classes := make(map[string]models.Classes)
	//查班级
	var classList []models.Classes
	lib.Engine.Table("classes").Find(&classList)
	if len(departmentList) != 0 {
		for _, value := range classList {
			classes[value.Name] = value
			println("")
			fmt.Printf("班级id:%d,班级名称：%v", value.Id, classes[value.Name])
		}
	}

	//查教师
	updateTeachers := make(map[string]models.Teacher)
	var teacherList []models.Teacher
	lib.Engine.Table("teacher").Find(&teacherList)
	if len(departmentList) != 0 {
		for _, value := range teacherList {
			updateTeachers[value.Code] = value
			println("")
			fmt.Printf("教师id:%d,教师名称：%v", value.Id, value.Name)
		}
	}

	updateNum := 0
	createAt := time.Now().Format("2006-01-02 15:04:05")
	rows := f.GetRows("Sheet1")

	//2.将excel表数据整理成golang数据
	//IdList := make(map[string]int)
	rowIndex := make(map[string]int) //列名对应的索引，这样导入的时候不用担心列名顺序的问题
	println("row len:", len(rows))
	if len(rows) == 0 {
		println("教师信息为空。")
		ctx.JSON(lib.NewResponseFail(1, "更新失败,教师信息为空。"))
		return
	}
	teachers := make([]UpdateTeacher, len(rows))

	//插入语句
	sqlTeaClaInsertStr := "insert into teacher_class(teacher_id,department_id,class_id,create_at) values"

	for indexRow, row := range rows {
		updateNum++
		if updateNum == 1 {
			//第一行根据列名导入
			var columns = []string{"教工号", "姓名", "性别", "院系", "体育班"}
			rowIndex = lib.GetRowIndex(row, columns)
			continue
		}
		teachers[indexRow] = UpdateTeacher{
			Code: row[rowIndex["教工号"]],
			Name: row[rowIndex["姓名"]],
		}
		if row[rowIndex["性别"]] == "男" {
			teachers[indexRow].Gender = "1"
		} else {
			teachers[indexRow].Gender = "2"
		}

		teachers[indexRow].DepartmentId = strconv.Itoa(departments[row[rowIndex["院系"]]])

		fmt.Printf("(%s,'%s',%s,%s,'%s'),", teachers[indexRow].Code, teachers[indexRow].Name, teachers[indexRow].Gender, teachers[indexRow].DepartmentId, createAt)
		println("")
		println("教工号：", row[rowIndex["教师"]], "", teachers[indexRow].Name)
		if updateTeachers[teachers[indexRow].Code].Id == 0 {
			println("文件上传失败，教工号：" + row[rowIndex["教师"]] + "不存在，请先添加教师")
			ctx.JSON(lib.NewResponseFail(0, "文件上传失败，教工号："+row[rowIndex["教师"]]+"不存在，请先添加教师"))
			return
		}

		println("院系名称：", row[rowIndex["院系"]], "院系id:", departments[row[rowIndex["院系"]]])
		//判断要修改的院系是否存在
		if departments[row[rowIndex["院系"]]] == 0 {
			println("文件上传失败，院系：" + row[rowIndex["院系"]] + "不存在，请先添加院系")
			ctx.JSON(lib.NewResponseFail(0, "文件上传失败，院系："+row[rowIndex["院系"]]+"不存在，请先添加院系"))
			return
		}

		//教师id
		teachId := updateTeachers[teachers[indexRow].Code].Id

		//班级id
		classId := classes[classList[indexRow].Name].Id

		//该班级的院系id
		departId := classes[classList[indexRow].Name].DepartmentId

		//teacher_id int,department_id int,class_id int,create_at string
		sqlTeaClaInsertStr += fmt.Sprintf("(%d,'%d',%d,'%s'),", teachId, departId, classId, createAt)

	}
	println(sqlTeaClaInsertStr)
	//3.更新teacher表
	teachersInfo := teachers[1:len(teachers)]
	sqlUpTeacherStr := lib.UpdateBatch("teacher", teachersInfo)

	// 创建 Session 对象
	sess := lib.Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err = sess.Begin(); err != nil {
		fmt.Println(err)
		ctx.JSON(lib.NewResponseFail(1, "更新失败"))
		return
	}

	//更新teacher表
	//_, err = lib.Engine.Query(sqlUpTeacherStr)
	//if err != nil {
	//	fmt.Println(err)
	//	ctx.JSON(lib.NewResponseFail(1, "更新失败"))
	//	return
	//}

	//4.更新teacher_class表

	//4.1删除旧的教师_班级对应记录
	//获取所有教师工号
	println("获取所有教师工号")
	delTeaClaWhere := ""
	for _, vv := range teachersInfo {
		t := reflect.ValueOf(vv)
		delTeaClaWhere += fmt.Sprintf("'%s', ", GetReferenceValue(t, "Code"))
	}
	println(delTeaClaWhere)
	delSqlStr := fmt.Sprintf("DELETE teacher_class FROM teacher_class INNER JOIN teacher ON teacher.id=teacher_class.teacher_id WHERE teacher.code IN(%s)", strings.TrimRight(delTeaClaWhere, ", "))
	println(delSqlStr)
	//_, err = lib.Engine.Query(delSqlStr)
	//if err != nil {
	//	fmt.Println(err)
	//	ctx.JSON(lib.NewResponseFail(1, "删除失败"))
	//	return
	//}

	//4.2添加新的教师_班级对应记录
	inserTeaCla := strings.TrimRight(sqlTeaClaInsertStr, ", ")
	//_, err = lib.Engine.Query(inserTeaCla)
	//if err != nil {
	//	fmt.Println(err)
	//	ctx.JSON(lib.NewResponseFail(1, "插入失败"))
	//	return
	//}

	if _, err = sess.Query(sqlUpTeacherStr); err != nil {
		sess.Rollback()
		fmt.Println(err)
		ctx.JSON(lib.NewResponseFail(1, "更新失败"))
		return
	} else if _, err = sess.Query(delSqlStr); err != nil {
		sess.Rollback()
		fmt.Println(err)
		ctx.JSON(lib.NewResponseFail(1, "更新失败"))
		return
	} else if _, err = sess.Query(inserTeaCla); err != nil {
		sess.Rollback()
		fmt.Println(err)
		ctx.JSON(lib.NewResponseFail(1, "更新失败"))
		return
	}

	// 完成事务
	sess.Commit()

	ctx.JSON(lib.NewResponseOK("修改成功"))
	return

}

//通用更新传入参数的结构，结构的所有字段数据类型应为string
func UpdateBatch(tableName string, models []UpdateTeacher) string {
	var columns []Column
	var sqlStr string
	var whereIn string
	//从数据集中取第一行，获取列名
	columns = GetFieldAttr(models[0])
	referenceTag := columns[0].Tag
	referenceName := columns[0].Name
	if tableName != "" {
		sqlStr += "UPDATE " + tableName + " SET "
		for _, column := range columns {
			if column.Tag != referenceTag {
				sqlStr += column.Tag + " = CASE "
				for _, vv := range models {
					t := reflect.ValueOf(vv)
					referenceValue := GetReferenceValue(t, column.Name)
					sqlStr += fmt.Sprintf("WHEN %s = '%s' THEN '%s' ", column.Tag, referenceValue, t.FieldByName(column.Name))
				}
				sqlStr += "ELSE " + column.Tag + " END, "
			}
		}
		for _, vv := range models {
			t := reflect.ValueOf(vv)
			whereIn += fmt.Sprintf("'%s', ", GetReferenceValue(t, referenceName))
		}
		sqlStr = fmt.Sprintf("%s WHERE %s IN(%s)", strings.TrimRight(sqlStr, ", "), referenceTag, strings.TrimRight(whereIn, ", "))

	}

	fmt.Println(sqlStr)
	return sqlStr
}

//获取数据，为了主键不是Int类型时通用，全部返回string
func GetReferenceValue(t reflect.Value, referenceColumn string) interface{} {

	var referenceValue interface{}

	if t.FieldByName(referenceColumn).Type().String() == "int" {
		referenceValue = Int642str(t.FieldByName(referenceColumn).Int())
	} else {
		referenceValue = t.FieldByName(referenceColumn)
	}
	return referenceValue
}

func GetFieldAttr(structName interface{}) []Column {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	var columns []Column
	for i := 0; i < fieldNum; i++ {
		if t.Field(i).Name != "" {
			var column = Column{t.Field(i).Name, t.Field(i).Tag.Get("json")}
			columns = append(columns, column)
		}
	}
	return columns
}

func Int642str(t int64) string {
	str := strconv.FormatInt(int64(t), 10)
	return str
}
