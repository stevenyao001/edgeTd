package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App      appConf        `mapstructure:"app"`
	Log      logConf        `mapstructure:"log"`
	Mqtt     []MqttConf     `mapstructure:"mqtt"`
	TDengine []TdengineConf `mapstructure:"tdengine"`
}

type appConf struct {
	AppId     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
	Env       string `mapstructure:"env"`
}

type logConf struct {
	MainPath string `mapstructure:"main_path"`
}

type MqttConf struct {
	InsName  string `mapstructure:"ins_name"`
	Addr     string `mapstructure:"addr"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	ClientId string `mapstructure:"client_id"`
}

// tdengine配置
type TdengineConf struct {
	InsName     string `mapstructure:"ins_name"`
	Fqdn        string `mapstructure:"fqdn"`
	Port        string `mapstructure:"port"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	DbName      string `mapstructure:"db_name"`
}

func InitConf(filePath string, conf interface{}) {
	//设置配置文件类型
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		panic("init config fail :" + err.Error())
	}

	if err := viper.Unmarshal(conf); err != nil {
		panic("resolution config fail :" + err.Error())
	}
}
