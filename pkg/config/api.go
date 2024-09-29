package config

import (
	"github.com/spf13/viper"
)

/*
 * 创建配置文件
 * 参数:
 *	AppName: 当 AppName 不为空时，将会从环境变量中读取配置
 *	configFile: 如果使用配置文件需要提供文件并执行 ReadConfigFromFile
 *	defaultConfig: 提供一个设置默认配置的回调函数
 * 返回:
 *	配置对象
 */
func New(AppName string,
	configFile string,
	defaultConfig func(*viper.Viper)) Config {

	// 一定要用New
	c := viper.New()
	cfg := config{
		configObj: c,
	}
	cfg.configFile = configFile
	cfg.AppName = AppName

	// 读取环境配置
	cfg.readViperConfig()

	// 读取默认配置
	if defaultConfig != nil {
		defaultConfig(cfg.configObj)
	}

	return &cfg
}

type Config interface {

	/*
	 * 获取配置项
	 */
	Config() *viper.Viper

	/*
	 * 从配置文件加载配置
	 */
	ReadConfigFromFile() error

	/*
	 * 更新配置到配置文件
	 */
	UpdateConfigToFile() error
}
