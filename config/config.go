package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Init 初始化配置文件(参数顺序:name,path,type)，默认配置文件为config.yaml
func Init(args ...string) {
	workPath, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	configPath := "."
	configName := "config"
	configType := "yml"
	if len(args) > 0 {
		configName, configPath, configType = resolvePath(args[0])
	}
	configPath = filepath.Join(workPath, configPath)
	configType = strings.ToLower(configType)
	if configType != "yaml" && configType != "yml" && configType != "json" && configType != "toml" && configType != "properties" {
		panic("不支持的配置文件格式")
	}
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType(configType)
	err = viper.ReadInConfig()
	if err != nil {
		panic("加载配置文件错误: " + err.Error())
	}
	fmt.Println("配置初始化成功...")
}

func resolvePath(arg string) (fName, fPath, fType string) {
	fName = "conf"
	fPath = "."
	fType = "yml"
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return
	}
	//兼容windows
	arg = strings.ReplaceAll(arg, "\\\\", "/")
	index1 := strings.LastIndex(arg, "/")
	index2 := strings.LastIndex(arg, ".")
	if index1 != -1 {
		fPath = arg[0:index1]
	}
	if index2 != -1 {
		fName = arg[index1+1 : index2]
		fType = arg[index2+1:]
	}
	return
}

// Get 获取所有类型的配置
func Get(key string) interface{} {
	return viper.Get(key)
}

// GetString 获取字符串类型的配置
func GetString(key string) string {
	return viper.GetString(key)
}

// GetBool 获取布尔类型的配置
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetInt 获取整数类型的配置
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetStringSlice 获取字符串数组类型的配置
func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

// GetStringMap 获取map接口配置
func GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

// AllSettings 所有配置
func AllSettings() map[string]interface{} {
	return viper.AllSettings()
}

// UnmarshalKey 根据key解析配置到指定的结构中
func UnmarshalKey(key string, config interface{}) error {
	err := viper.UnmarshalKey(key, config)
	if err != nil {
		fmt.Printf("解析配置失败 key:[%s], %v \n", key, err)
	}
	return err
}

// Unmarshal 解析整个配置到指定的结构中
func Unmarshal(config interface{}) error {
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("解析配置失败: %v \n", err)
	}
	return err
}
