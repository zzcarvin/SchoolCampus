package lib

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type Column struct {
	Name string
	Tag  string
}

func GetRowIndex(rows []string, columns []string) map[string]int {
	rowIndex := make(map[string]int)
	for kk, vv := range rows {
		for _, cv := range columns {
			if vv == cv {
				rowIndex[cv] = kk
				break
			}
		}
	}
	return rowIndex
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func SaveUploadedFile(fh *multipart.FileHeader, destDirectory string, fileName string) (int64, error) {
	//execl := path.Join(pwd, "execl")
	//println("\n*********************\n\n当前execl文件夹目录为\n\n***********************", execl)
	//execlDir := path.Join(pwd, "execl", fname)
	//println("\n*********************\n\n当前execl表格目录为\n\n***********************", execlDir)
	exist, err := PathExists(destDirectory)
	if err != nil {
		return 0, err
	}
	if !exist {
		// 创建文件夹
		err := os.Mkdir(destDirectory, 0777)
		if err != nil {
			return 0, err
		}
	}
	src, err := fh.Open()
	if err != nil {
		return 0, err
	}
	defer src.Close()
	filePath := filepath.Join(destDirectory, fileName)
	out, err := os.OpenFile(filePath,
		os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	defer out.Close()

	return io.Copy(out, src)
}

func Execl() {
	xlsx, err := excelize.OpenFile("./lib/中文系中文1班4000人.xlsx")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get sheet index.
	index := xlsx.GetSheetIndex("Sheet2")
	// Get all the rows in a sheet.
	rows := xlsx.GetRows("sheet" + strconv.Itoa(index))
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}

//通用更新传入参数的结构，结构的所有字段数据类型应为string
func NewUpdateBatch(tableName string, models interface{}) string {

	//将任意类型转成[]interface{}
	modelsArr, _ := CreateAnyTypeSlice(models)

	var columns []Column
	var sqlStr string
	var whereIn string
	//从数据集中取第一行，获取列名
	columns = GetFieldAttr(modelsArr[0])
	referenceTag := columns[0].Tag
	referenceName := columns[0].Name
	if tableName != "" {
		sqlStr += "UPDATE " + tableName + " SET "
		for _, column := range columns {
			if column.Tag != referenceTag {
				sqlStr += column.Tag + " =  "
				for _, vv := range modelsArr {
					t := reflect.ValueOf(vv)
					//referenceValue := GetReferenceValue(t, column.Name)
					sqlStr += fmt.Sprintf("WHEN '%s' THEN '%s' ", GetReferenceValue(t, referenceName), t.FieldByName(column.Name))
				}
				sqlStr += " , "
			}
		}
		for _, vv := range modelsArr {
			t := reflect.ValueOf(vv)
			whereIn += fmt.Sprintf("'%s', ", GetReferenceValue(t, referenceName))
		}
		sqlStr = fmt.Sprintf("%s WHERE %s IN(%s)", strings.TrimRight(sqlStr, ", "), referenceTag, strings.TrimRight(whereIn, ", "))

	}

	fmt.Println(sqlStr)
	return sqlStr
}

//获取数据，为了主键不是Int类型时通用，全部返回string
//func GetReferenceValue(t reflect.Value, referenceColumn string) interface{} {
//
//	var referenceValue interface{}
//
//	if t.FieldByName(referenceColumn).Type().String() == "int" {
//		referenceValue = Int642str(t.FieldByName(referenceColumn).Int())
//	} else {
//		referenceValue = t.FieldByName(referenceColumn)
//	}
//	return referenceValue
//}

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

//将任意类型转成[]interface{}
func CreateAnyTypeSlice(slice interface{}) ([]interface{}, bool) {
	val, ok := isSlice(slice)

	if !ok {
		return nil, false
	}

	sliceLen := val.Len()

	out := make([]interface{}, sliceLen)

	for i := 0; i < sliceLen; i++ {
		out[i] = val.Index(i).Interface()
	}

	return out, true
}

func isSlice(arg interface{}) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)

	if val.Kind() == reflect.Slice {
		ok = true
	}

	return
}

//mysql语句例子如下：
//UPDATE student SET name = CASE code WHEN '1314' THEN '孙德冬' ELSE name END, gender = CASE code WHEN '1314' THEN '1' ELSE gender END, class_id = CASE code WHEN '1314' THEN '99' ELSE class_id END, department_id
//= CASE code WHEN '1314' THEN '8' ELSE department_id END, plan_id = CASE code WHEN '1314' THEN '0' ELSE plan_id END WHERE code IN('1314')
func UpdateBatch(tableName string, models interface{}) string {

	//将任意类型转成[]interface{}
	modelsArr, _ := CreateAnyTypeSlice(models)

	var columns []Column
	var sqlStr string
	var whereIn string
	//从数据集中取第一行，获取列名
	columns = GetFieldAttr(modelsArr[0])
	referenceTag := columns[0].Tag   //json字段名
	referenceName := columns[0].Name //结构体字段名
	println("columns[0].Tag:" + referenceTag + ",columns[0].Name:" + referenceName)
	if tableName != "" {
		sqlStr += "UPDATE " + tableName + " SET "
		for _, column := range columns {
			println("")
			fmt.Printf("column:%s", column)
			if column.Tag != referenceTag { //唯一标识符不更改
				sqlStr += column.Tag + " = CASE " + referenceTag + " "
				for _, vv := range modelsArr {
					t := reflect.ValueOf(vv)
					sqlStr += fmt.Sprintf("WHEN '%s' THEN '%s' ", GetReferenceValue(t, referenceName), t.FieldByName(column.Name))
				}
				sqlStr += "ELSE " + column.Tag + " END, "
			}
		}

		for _, vv := range modelsArr {
			t := reflect.ValueOf(vv)
			whereIn += fmt.Sprintf("'%s', ", GetReferenceValue(t, referenceName))
		}
		sqlStr = fmt.Sprintf("%s WHERE %s IN(%s)", strings.TrimRight(sqlStr, ", "), referenceTag, strings.TrimRight(whereIn, ", "))

	}

	fmt.Println(sqlStr)
	return sqlStr
}
