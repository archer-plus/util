package captchx

import (
	"context"
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/archer-plus/util/config"

	"github.com/go-redis/redis/v8"
	"github.com/mojocn/base64Captcha"
)

var once sync.Once
var customStore *RedisStore

var defaultConfig = &ConfigBody{
	Id:          base64Captcha.RandomId(),
	CaptchaType: "string",
	DriverString: &base64Captcha.DriverString{
		Height:          60,
		Width:           240,
		NoiseCount:      10,
		ShowLineOptions: 0,
		Length:          6,
		Source:          "1234567890abcdefghijklmnopquestuvwxyz",
		BgColor: &color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		},
	},
}

type RedisStore struct {
	client *redis.Client
}

type ConfigBody struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

// GenerateCaptcha 生成图形码
func GenerateCaptcha(param *ConfigBody) (id, b64s string, err error) {
	if param == nil {
		param = defaultConfig
	}
	var driver base64Captcha.Driver
	switch param.CaptchaType {
	case "audio":
		driver = param.DriverAudio
	case "string":
		driver = param.DriverString.ConvertFonts()
	case "math":
		driver = param.DriverMath.ConvertFonts()
	case "chinese":
		driver = param.DriverChinese.ConvertFonts()
	default:
		driver = param.DriverDigit
	}
	c := base64Captcha.NewCaptcha(driver, NewRedisStore())
	id, b64s, err = c.Generate()
	return
}

func VerifyCaptcha(param *ConfigBody) bool {
	stroe := NewRedisStore()
	return stroe.Verify(param.Id, param.VerifyValue, true)
}

func (c *RedisStore) Set(id string, value string) {
	err := c.client.Set(context.TODO(), id, value, time.Minute*10).Err()
	if err != nil {
		fmt.Errorf("%v", err)
	}
}

func (c *RedisStore) Get(id string, clear bool) string {
	val, err := c.client.Get(context.TODO(), id).Result()
	if err != nil {
		fmt.Errorf("%v", err)
		return ""
	}
	if clear {
		err := c.client.Del(context.TODO(), id).Err()
		if err != nil {
			fmt.Errorf("%v", err)
			return ""
		}
	}
	return val
}

func (c *RedisStore) Verify(id, answer string, clear bool) bool {
	v := c.Get(id, clear)
	return v == answer
}

//  NewRedisStore 使用前需要初始化redis
func NewRedisStore() *RedisStore {
	once.Do(func() {
		customStore = &RedisStore{client: redis.NewClient(&redis.Options{
			Addr:     config.GetString("captcha.redis.conn"),
			Password: config.GetString("captcha.redis.password"),
			DB:       config.GetInt("captcha.redis.db"),
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
			IdleTimeout:     time.Second * time.Duration(config.GetInt("captcha.redis.timeout")),
		})}
	})
	return customStore
}
