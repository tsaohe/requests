package requests

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

/*
*******************************************************************************************
 */

type Config struct {
	System   SystemConfig   `mapstructure:"system" json:"system" yaml:"system"`
	Wechat   WechatConfig   `yaml:"wechat" mapstructure:"wechat" json:"wechat"`
	BoKa     BoKaConfig     `yaml:"bokaApi" mapstructure:"bokaApi" json:"bokaApi"`
	Constant ConstantConfig `yaml:"constant" json:"constant" mapstructure:"constant"`
}

type SystemConfig struct {
	Name   string `mapstructure:"name" json:"name" yaml:"name"`
	Port   string `mapstructure:"port" json:"port" yaml:"port"`
	Mode   string `mapstructure:"mode" json:"mode" yaml:"mode"`
	Env    string `mapstructure:"env" json:"env" yaml:"env"`
	Domain string `mapstructure:"domain" json:"domain" yaml:"domain"`
}

type WechatConfig struct {
	AppID       string `yaml:"appId" mapstructure:"appId" json:"appId"`
	AppSecret   string `yaml:"appSecret" mapstructure:"appSecret" json:"appSecret"`
	Token       string `yaml:"token" mapstructure:"token" json:"token"`
	AesKey      string `yaml:"aesKey" mapstructure:"aesKey" json:"aesKey"`
	AdminOpenid string `yaml:"adminOpenid" mapstructure:"adminOpenid" json:"adminOpenid"`
}

type BoKaConfig struct {
	Sec      int64  `mapstructure:"sec" json:"sec" yaml:"sec"`
	CustId   string `json:"custId" yaml:"custId" mapstructure:"custId" `
	CompId   string `json:"compId" yaml:"compId" mapstructure:"compId" `
	UserName string `json:"userName" yaml:"userName" mapstructure:"userName" `
	PassWord string `json:"passWord" yaml:"passWord" mapstructure:"passWord" `
	Source   string `json:"source" yaml:"source" mapstructure:"source" `
}

type ConstantConfig struct {
	ExceptUserList []string `json:"exceptUserList" yaml:"exceptUserList,flow" mapstructure:"exceptUserList"`
}

/*
*******************************************************************************************
 */

func InitConf(path ...string) {

	var config string
	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "choose conf file.")
		flag.Parse()
		if config == "" { // 优先级: 命令行 > 环境变量 > 默认值
			if configEnv := os.Getenv("CONFFILE"); configEnv == "" {
				config = ConfigFile
				log.Printf("您正在使用config的默认值,config的路径为%v\n", config)
			} else {
				config = configEnv
				log.Printf("您正在使用configEnv环境变量,config的路径为%v\n", config)
			}
		} else {
			log.Printf("您正在使用命令行的-c参数传递的值,config的路径为%v\n", config)
		}
	} else {
		config = path[0]
		log.Printf("您正在使用InitConf传递的值,config的路径为%v\n", config)
	}
	v := viper.New()
	v.SetConfigFile(config)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error conf file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		log.Println("conf file changed:", e.Name)
		if err = v.Unmarshal(&CONFIG); err != nil {
			log.Println(err)
		}
	})
	if err = v.Unmarshal(&CONFIG); err != nil {
		log.Println(err)
	}
	VP = v
}
