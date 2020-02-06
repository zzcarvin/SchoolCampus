package internal

import (
	"Campus/configs"
	"Campus/internal/lib"
)

func Init() error {
	var err error

	// 加载配置
	err = configs.Load()
	if err != nil {
		return err
	}

	//初始化redis
	err = lib.RedisInit()
	if err != nil {
		return err
	}

	// 初始化xorm
	err = lib.XormInit()
	if err != nil {
		return err
	}

	// todo zaplog日志初始化, 需要优化
	lib.ZapCoreInit()

	//定时任务
	lib.NewCron()

	//go func() {
	//	for {
	//		//planrecordid出队列
	//		key := "finishRunJob"
	//
	//		run.GetJob(key)
	//
	//	}
	//}()

	return nil
}
