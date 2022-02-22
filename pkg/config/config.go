package config

import (
	"gohub/pkg/helpers"
	"os"

	"github.com/spf13/cast"
	viperlib "github.com/spf13/viper"
)

var viper *viperlib.Viper

// 动态加载配置信息
type ConfigFunc func() map[string]interface{}

// 先加载到此数组, loadconfig 在动态生成配置信息
var ConfigFuncs map[string]ConfigFunc

func init() {
	// 初始化 Viper 库
	viper = viperlib.New()
	// 配置类型, 支持 "json", "toml", "yaml", "env", "dotenv",
	viper.SetConfigType("env")
	// 环境变量配置文件查找的路径, 相对于 main.go
	viper.AddConfigPath(".")
	// 设置环境变量前缀, 用以区分 go 的环境变量
	viper.SetEnvPrefix("appenv")
	// 读取环境变量(支持 flags)
	viper.AutomaticEnv()

	ConfigFuncs = make(map[string]ConfigFunc)
}

// 初始化配置信息, 完成对环境变量已经 config 信息的加载
func InitConfig(env string) {
	// 加载环境变量
	loadEnv(env)
	// 注册配置信息
	loadConfig()
}

func loadConfig() {
	for name, fn := range ConfigFuncs {
		viper.Set(name, fn())
	}
}

func loadEnv(envSuffix string) {
	// 默认加载 .env 文件, 如果有传参 --env=name 的话, 加载 .env.name 文件
	envPath := ".env"

	if len(envSuffix) > 0 {
		filepath := ".env." + envSuffix
		if _, err := os.Stat(filepath); err == nil {
			envPath = filepath
		}
	}

	// 加载 env
	viper.SetConfigName(envPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// 监控 .env 文件, 变更时重新加载
	viper.WatchConfig()
}

// Env 读取环境变量, 支持默认值
func Env(envName string, defaultValue ...interface{}) interface{} {
	if len(defaultValue) > 0 {
		return internalGet(envName, defaultValue[0])
	}

	return internalGet(envName)
}

// 新增配置项
func Add(name string, configFn ConfigFunc) {
	ConfigFuncs[name] = configFn
}

// 获取配置项
func Get(path string, defaultValue ...interface{}) string {
	return GetString(path, defaultValue...)
}

func internalGet(path string, defaultValue ...interface{}) interface{} {
	if !viper.IsSet(path) || helpers.Empty(viper.Get(path)) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return viper.Get(path)
}

// 获取 String 类型的配置
func GetString(path string, defaultValue ...interface{}) string {
	return cast.ToString(internalGet(path, defaultValue...))
}

// 获取 Int 类型的配置
func GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(internalGet(path, defaultValue...))
}

func GetFloat(path string, defaultValue ...interface{}) float64 {
	return cast.ToFloat64(internalGet(path, defaultValue...))
}

func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(internalGet(path, defaultValue...))
}

// GetUint 获取 Uint 类型的配置信息
func GetUint(path string, defaultValue ...interface{}) uint {
	return cast.ToUint(internalGet(path, defaultValue...))
}

// GetBool 获取 Bool 类型的配置信息
func GetBool(path string, defaultValue ...interface{}) bool {
	return cast.ToBool(internalGet(path, defaultValue...))
}

// GetStringMapString 获取结构数据
func GetStringMapString(path string) map[string]string {
	return viper.GetStringMapString(path)
}
