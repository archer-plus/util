package mongox

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/archer-plus/util/config"
	"github.com/archer-plus/util/logx"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ErrorDataExisted = "E11000"                        // 数据已存在
	ErrorDataIsEmpty = "mongo: no documents in result" // 没有查询到数据
)

// Config mongo配置
type Config struct {
	Database   string `mapstructure:"db"`
	Conn       string `mapstructure:"conn"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	AuthSource string `mapstructure:"auth_source"`
	Timeout    int    `mapstructure:"timeout"`
	MaxPool    int    `mapstructure:"max_pool"`
}

var instance *mongo.Client
var conf = Config{}
var once sync.Once

// Init 初始化mongo数据库
func Init() {
	RegisterMongoDB()
}

// NewID 生成ObjectID
func NewID() primitive.ObjectID {
	return primitive.NewObjectID()
}

// RegisterMongoDB 注册mongodb
func RegisterMongoDB() {
	err := config.UnmarshalKey("mongo", &conf)
	if err != nil {
		logx.Sugar.Warnf("未发现mongo配置信息, %v", err)
		return
	}
	if conf.Conn == "" {
		logx.Sugar.Errorf("mongo配置信息错误, %v", "mongo", err)
		return
	}
	once.Do(func() {
		if conf.Timeout <= 0 {
			conf.Timeout = 5
		}
		if conf.MaxPool <= 0 {
			conf.MaxPool = 10
		}
		clientOptions := options.Client().ApplyURI(conf.Conn).SetConnectTimeout(time.Duration(conf.Timeout) * time.Second).SetMaxPoolSize(uint64(conf.MaxPool))
		clientOptions.SetAuth(options.Credential{
			AuthMechanism: "SCRAM-SHA-1",
			AuthSource:    conf.Database,
			Username:      conf.Username,
			Password:      conf.Password,
		})
		instance, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			logx.Sugar.Errorf("创建连接mongo失败, %v", "mongo", err)
			return
		}
		err = instance.Ping(context.TODO(), nil)
		if err != nil {
			logx.Sugar.Errorf("连接mongo失败, %v", "mongo", err)
			return
		}
		fmt.Println("mongo 注册成功...")
	})
}

// New 根据数据库名称，集合名称获取集合操作对象；参数顺序database_name,collection_name;
func New(key ...string) (*mongo.Collection, error) {
	l := len(key)
	var collection *mongo.Collection
	if l == 0 {
		return nil, errors.New("请输入数据库名称及集合名称")
	}
	if l == 1 {
		collection = instance.Database(conf.Database).Collection(key[0])
	} else {
		collection = instance.Database(key[0]).Collection(key[1])
	}
	return collection, nil
}

// IndexExisted 检查索引是否存在
func IndexExisted(coll, index string) bool {
	c, err := New(coll)
	if err != nil {
		panic(err)
	}
	curs, err := c.Indexes().List(context.TODO(), nil)
	if err != nil {
		return true
	}
	for curs.Next(context.TODO()) {
		curr, _ := strconv.Unquote(curs.Current.Lookup("name").String())
		if index == curr {
			return true
		}
	}
	if curs.Err() != nil {
		return true
	}
	err = curs.Close(context.TODO())
	if err != nil {
		return true
	}
	return false
}

// Close 关闭Mongodb连接
func Close() {
	if instance != nil {
		err := instance.Disconnect(context.TODO())
		if err != nil {
			logx.Sugar.Fatal(err)
		}
	}
}
