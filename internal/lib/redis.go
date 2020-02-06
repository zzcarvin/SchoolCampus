package lib

import (
	"Campus/configs"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

//全局的Redis数据库连接池
var RedisPool *redis.Pool

//Redis连接池初始化
func RedisInit() error {
	if RedisPool != nil {
		return fmt.Errorf("Redis已经初始化")
	}

	//获取配置
	cfg := configs.Conf.Redis

	//创建连接池
	RedisPool = &redis.Pool{
		MaxIdle:     30,
		MaxActive:   0,
		IdleTimeout: 3000 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(cfg.Network, cfg.Address)
			if err != nil {
				fmt.Println(" redis dial err:", err)
				return nil, err
			}
			//鉴权
			if cfg.Auth != "" {
				c.Do("AUTH", cfg.Auth)
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

//Redis连接池关闭
func RedisClose() error {
	err := RedisPool.Close()
	if err != nil {
		//TODO log error
	}
	return err
}

//获取Redis连接，需要调用 conn.Close()来释放连接
func GetRedisConn() redis.Conn {
	if RedisPool == nil {
		//TODO log error
		log.Print("redisPool is nil")
		return nil
	}
	return RedisPool.Get()
}


