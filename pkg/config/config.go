package config

import (
	"strings"

	"github.com/spf13/viper"
)

type config struct {
	AppName    string       // 应用名称
	configFile string       // 配置文件名
	configObj  *viper.Viper // 配置文件对象
}

/*
 * 获取配置项
 */
func (config *config) Config() *viper.Viper {
	return config.configObj
}

/*
 * 从配置文件加载配置
 */
func (config *config) ReadConfigFromFile() error {
	config.configObj.SetConfigFile(config.configFile)
	return config.configObj.ReadInConfig()
}

/*
 * 更新配置到配置文件
 */
func (config *config) UpdateConfigToFile() error {
	return config.configObj.WriteConfig()
}

/*
 * 内部接口：读取环境配置项目
 */
func (config *config) readViperConfig() {
	config.configObj.SetEnvPrefix(strings.ToUpper(config.AppName))
	config.configObj.AutomaticEnv()
}
