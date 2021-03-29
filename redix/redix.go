package redix

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/archer-plus/util"

	"github.com/archer-plus/util/logx"

	"github.com/archer-plus/util/config"

	"github.com/go-redis/redis/v8"
)

// Config Redis配置信息
type Config struct {
	Name     string `mapstructure:"name"`
	Conn     string `mapstructure:"conn"`
	Password string `mapstructure:"password"`
	Timeout  int    `mapstructure:"timeout"`
	DB       int    `mapstructure:"db"`
}

// RedisClient Redis客户端
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

// Init 初始化Redis
func Init() {
	RegisterRedis()
}

var instance map[string]*RedisClient = make(map[string]*RedisClient)

// RegisterRedis 注册
func RegisterRedis() {
	configSlice := make([]*Config, 0)
	err := config.UnmarshalKey("redis", &configSlice)
	if err != nil {
		logx.Sugar.Errorf("解析配置失败 key:[%s], %v", "redis", err)
		return
	}
	if len(configSlice) == 0 {
		logx.Sugar.Errorf("未发现redis配置信息, %v", "redis", err)
		return
	}
	for _, config := range configSlice {
		c := new(RedisClient)
		c.ctx = context.Background()
		if config.Timeout == 0 {
			config.Timeout = 60
		}
		c.client = redis.NewClient(&redis.Options{
			Addr:     config.Conn,
			Password: config.Password,
			DB:       config.DB,
			//连接池容量及闲置连接数量
			PoolSize:     15,              // 连接池数量
			MinIdleConns: 10,              //好比最小连接数
			DialTimeout:  5 * time.Second, //连接建立超时时间
			ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
			WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
			PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

			//闲置连接检查包括IdleTimeout，MaxConnAge
			IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
			MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

			//命令执行失败时的重试策略
			MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
			MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
			MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
			IdleTimeout:     time.Second * time.Duration(config.Timeout),
		})
		// 测试是否连接成功
		err = c.client.Ping(c.ctx).Err()
		if err != nil {
			panic("redis 连接失败")
		}
		instance[config.Name] = c
	}
	fmt.Println("redis 注册成功...")
}

// DB 根据名称获取数据库连接
func DB(name string) *RedisClient {
	v, ok := instance[name]
	if !ok {
		panic("get " + name + " redis db failed")
	}
	return v
}

// Set 将字符串值 value 关联到 key
func (c *RedisClient) Set(key string, value interface{}, ops ...*Operation) *Result {
	exp := Operations(ops).Find(EXPIRE).Result(time.Second * 0).(time.Duration)
	nx := Operations(ops).Find(NX).Result(nil)

	if nx != nil {
		return NewResult(c.client.SetNX(c.ctx, key, value, exp).Result())
	}
	xx := Operations(ops).Find(XX).Result(nil)
	if xx != nil {
		return NewResult(c.client.SetXX(c.ctx, key, value, exp).Result())
	}
	return NewResult(c.client.Set(c.ctx, key, value, exp).Result())
}

// Get 返回与键 key 相关联的字符串值
func (c *RedisClient) Get(key string) *Result {
	return NewResult(c.client.Get(c.ctx, key).Result())
}

// Keys 返回key列表
func (c *RedisClient) Keys(arg string) *Result {
	return NewResult(c.client.Keys(c.ctx, arg).Result())
}

// MGet 返回包含了所有给定键的值的列表
func (c *RedisClient) MGet(key string) *MResult {
	return NewMResult(c.client.MGet(c.ctx, key).Result())
}

// Del 删除key
func (c *RedisClient) Del(key string) *Result {
	return NewResult(c.client.Del(c.ctx, key).Result())
}

// Inrc 数字加一
func (c *RedisClient) Inrc(key string) *Result {
	return NewResult(c.client.Incr(c.ctx, key).Result())
}

// Lock 分布式锁
func (c *RedisClient) Lock(key string, acquire, timeout time.Duration) (string, error) {
	code := util.NewUUID()
	endTime := time.Now().Add(acquire).UnixNano()
	for time.Now().UnixNano() <= endTime {
		if success, err := c.client.SetNX(c.ctx, key, code, timeout).Result(); err != nil {
			return "", err
		} else if success {
			return code, nil
		} else if c.client.TTL(c.ctx, key).Val() == -1 {
			c.client.Expire(c.ctx, key, timeout)
		}
		time.Sleep(time.Millisecond)
	}
	return "", errors.New("lock timeout")
}

// UnLock 释放分布式锁
func (c *RedisClient) UnLock(key, code string) bool {
	txf := func(tx *redis.Tx) error {
		if v, err := tx.Get(c.ctx, key).Result(); err != nil && err != redis.Nil {
			return err
		} else if v == code {
			_, err := tx.Pipelined(c.ctx, func(pipe redis.Pipeliner) error {
				pipe.Del(c.ctx, key)
				return nil
			})
			return err
		}
		return nil
	}

	for {
		if err := c.client.Watch(c.ctx, txf, key); err == nil {
			return true
		} else if err == redis.TxFailedErr {
			logx.Sugar.Warnf("watch key is modified,retry to release lock. err: %s\n", err.Error())
			logx.Sugar.Warnf("key: %s,code: %s\n", key, code)
		} else {
			return false
		}
	}
}
