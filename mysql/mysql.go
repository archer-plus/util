package mysql

import (
	"fmt"
	"time"

	"github.com/archer-plus/util/config"
	"github.com/archer-plus/util/logx"

	// mysql 驱动
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Config mysql 配置信息
type Config struct {
	Name     string `mapstructure:"name"`     // 连接名称
	Host     string `mapstructure:"host"`     // 主机地址
	Port     int    `mapstructure:"port"`     // 连接端口号
	DB       string `mapstructure:"db"`       // 数据库名称
	User     string `mapstructure:"user"`     // 用户名
	Password string `mapstructure:"password"` // 密码
	MaxIDLE  int    `mapstructure:"max_idle"` // 最大空闲数
	MaxOpen  int    `mapstructure:"max_open"` // 最大连接数
	Timeout  int    `mapstructure:"timeout"`  // 连接超时
}

var instance map[string]*sqlx.DB = make(map[string]*sqlx.DB)

// Init 初始化mysql
func Init() {
	slice := make([]*Config, 0)
	err := config.UnmarshalKey("mysql", &slice)
	if len(slice) == 0 {
		logx.Sugar.Warnf("未发现mysql配置信息, %v", err)
		return
	}
	// 注册数据库
	for i, config := range slice {
		if i == 0 && config.Name == "" {
			RegisterDataBase("default", "mysql", dataSource(config), config.MaxIDLE, config.MaxOpen)
		} else {
			RegisterDataBase(config.Name, "mysql", dataSource(config), config.MaxIDLE, config.MaxOpen)
		}
	}
}

func dataSource(conf *Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%ds&charset=utf8mb4&loc=%s&parseTime=true",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DB, conf.Timeout, time.Local.String())
}

// RegisterDataBase 注册mysql数据库
func RegisterDataBase(aliasName, driverName, dataSource string, params ...int) {
	if aliasName == "" {
		logx.Sugar.Error("数据库别名为空")
		return
	}
	db, err := sqlx.Open(driverName, dataSource)
	if err != nil {
		panic(err.Error())
	}
	for i, v := range params {
		switch i {
		case 0:
			db.SetMaxIdleConns(v)
		case 1:
			db.SetMaxOpenConns(v)
		}
	}
	instance[aliasName] = db
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	logx.Sugar.Info("数据库 " + aliasName + " 注册成功...")
}

// ShowDataBase 展示数据库信息
func ShowDataBase() {
	logx.Sugar.Info("database: %v \n", instance)
}

// DB 根据名称获取数据库连接
func DB(key ...string) *sqlx.DB {
	if len(key) == 0 {
		return instance["default"]
	}
	return instance[key[0]]
}
