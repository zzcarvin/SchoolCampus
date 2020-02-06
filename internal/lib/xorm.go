package lib

import (
	"Campus/configs"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var Engine *xorm.Engine

func XormInit() error {
	if Engine != nil {
		return fmt.Errorf("Xorm已经初始化")
	}

	//获取配置
	cfg := configs.Conf.Database
	println("mysql连接信息:", cfg.Conn)
	var err error

	Engine, err = xorm.NewEngine(cfg.Driver, cfg.Conn)
	if err != nil {
		fmt.Printf("xorm初始化失败：%v", err.Error())
		return err
	}

	maxOpen := 1000
	maxIdle := 300
	debug := false
	if cfg.MaxIdle != 0 {
		maxIdle = cfg.MaxIdle
	}
	if cfg.MaxOpen != 0 {
		maxOpen = cfg.MaxOpen
	}
	if cfg.Debug != false {
		debug = cfg.Debug
	}
	fmt.Println("最大打开数据库链接", maxOpen)
	Engine.SetMaxIdleConns(maxOpen)
	Engine.SetMaxOpenConns(maxIdle)

	//打印调试信息
	if debug {
		Engine.ShowSQL(true)
		Engine.Logger().SetLevel(core.LOG_DEBUG)
	}

	//err = Engine.Sync2(new(models.PlanProgress), new(models.Plan), new(models.PlanRecord),new(models.PlanRoute))
	//if err != nil {
	//	print(err)
	//	return err
	//}

	//这个只执行一次，平时不要执行，这个要是每次都执行，相当于每次都重新创建数据库表，只要有不同的地方，就会改数据库，不推荐这样同步数据库，如果要改数据库结构，直接去数据库改
	return err
}

func XormClose() error {
	err := Engine.Close()
	if err != nil {
		//TODO log error
		fmt.Printf("%v", err.Error())
	}
	return err
}
