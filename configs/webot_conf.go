package configs

import (
	"fmt"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

type WebotConf struct {
	Name   string             `yaml:"name"`
	DbList map[string]*DbConf `yaml:"db"`
	Baidu  Baidu              `yaml:"baidu"`
	Upload Upload             `yaml:"upload"`
}

type Upload struct {
	Path string `yaml:"path"`
}

type Baidu struct {
	Ak string `yaml:"ak"`
	Sk string `yaml:"sk"`
}

type DbConf struct {
	Dsn string `yaml:"dsn"`
}

var (
	webotConfig *WebotConf
	once        sync.Once
)

func GetConf() *WebotConf {
	initWebotConf()
	return webotConfig
}

// initWebotConf 初始化配置
func initWebotConf() {
	once.Do(func() {
		yamlFile, err := ioutil.ReadFile("./conf/webot.yml")
		if err != nil {
			fmt.Println(err.Error())
		}
		var confTmp WebotConf
		err = yaml.Unmarshal(yamlFile, &confTmp)
		if err != nil {
			fmt.Println(err.Error())
		}
		webotConfig = &confTmp
	})
}
