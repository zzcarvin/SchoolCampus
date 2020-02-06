package configs

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

/**
 * 全局配置及默认值
 */
var Conf = struct {
	Ver      string `yaml:"ver"`
	Web      web
	Jwt      jwt
	Database database
	School   school
	Limit    limit
	Redis    redis
}{
	Ver: "1.0",
	Web: web{
		Addr:   ":8080",
		Static: "dist",
	},
	Jwt: jwt{
		Key: "wuxitaihu",
	},
	Database: database{
		Driver: "mysql",
	},
	Redis: redis{
		Network: "tcp",
	},
}

type web struct {
	Addr   string `yaml:"addr"`
	Static string `yaml:"s http://localhost:8085tatic"`
}

type jwt struct {
	Key           string `yaml:"key"`
	Expire        string `yaml:"expire"`
	RefreshExpire uint64 `yaml:"refreshexpire"`
}

type school struct {
	Name string `json:"name"`
}

/**
 * 数据库配置
 */
type database struct {
	Driver     string `yaml:"driver"`
	Conn       string `yaml:"conn"`
	MaxOpen    int    `yaml:"max_open"`
	MaxIdle    int    `yaml:"max_idle"`
	Debug      bool   `yaml:"debug"`
	CoreDBName string `yaml:"core_db_name"`
}

type limit struct {
	MinPace      int `yaml:"minPace"`
	BoyMinPace   int `yaml:"BoyminPace"`
	GirlMinPace  int `yaml:"GirlminPace"`
	MaxPace      int `yaml:"maxPace"`
	BoyMaxPace   int `yaml:"BoyMaxPace"`
	GirlMaxPace  int `yaml:"GirlMaxPace"`
	MinFrequency int `yaml:"minFrequency"`
	MaxFrequency int `yaml:"maxFrequency"`
}

/*
加载配置
*/
func Load() error {
	yamlFile, err := ioutil.ReadFile("./conf.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err #%v ", err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}
	//TODO 验证配置文件有效性

	return nil
}

/**
 * Redis缓存配置
 */
type redis struct {
	Network string `yaml:"network"`
	Address string `yaml:"address"`
	Auth    string `yaml:"auth"`
}
