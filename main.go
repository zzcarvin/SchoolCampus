package main

import (
	"Campus/internal"
	"fmt"
)

func main() {

	//初始化数据库连接
	err := internal.Init()
	if err != nil {
		fmt.Printf("初始化失败：%v", err.Error())
		return
	}

	//启动Web服务
	internal.WebServe()
}

//func CreateDir(dir string) (bool, error) {
//	_, err := os.Stat(dir)
//
//	if err == nil {
//		//directory exists
//		return true, nil
//	}
//
//	err2 := os.MkdirAll(dir, 0755)
//	if err2 != nil {
//		return false, err2
//	}
//
//	return true, nil
//
//	unc main() {
//		res2, err := fileManager.CreateDir("/LOG/PATH")    //创建文件夹
//		if res2 == false {
//			panic(err)
//		}
//		file, _ := os.OpenFile("/LOG/PATH/error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)    //打开日志文件，不存在则创建
//		defer file.Close()
//
//		log.SetOutput(file)    //设置输出流
//		log.SetPrefix("[Error]")    //日志前缀
//		log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)    //日志输出样式
//		        log.Println("Hi file")
//	}
