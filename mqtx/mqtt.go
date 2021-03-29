package mqtx

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/archer-plus/util/config"
	"github.com/archer-plus/util/logx"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var instance map[string]*Client = make(map[string]*Client)

// QoS 消息发送策略类型
type QoS int32

const (
	// QoS0 最多一次
	QoS0 QoS = 0
	// QoS1 至少一次
	QoS1 QoS = 1
	// QoS2 只有一次
	QoS2 QoS = 2
)

// Client mqtt客户端
type Client struct {
	client    mqtt.Client
	options   *mqtt.ClientOptions
	locker    *sync.Mutex
	observer  func(c *Client, payload []byte)
	onConnect func(c *Client)            // 连接成功调用
	onLost    func(c *Client, err error) // 关闭链接回调
}

// Options 连接选项
func (c *Client) Options() *mqtt.ClientOptions {
	return c.options
}

// GetClientID 获取客户端ID
func (c *Client) GetClientID() string {
	return c.options.ClientID
}

// Connect 连接
func (c *Client) Connect() error {
	return c.ensureConnected()
}

// ensureConnected 确保连接
func (c *Client) ensureConnected() error {
	if !c.client.IsConnected() {
		c.locker.Lock()
		defer c.locker.Unlock()
		if !c.client.IsConnected() {
			if token := c.client.Connect(); token.Wait() && token.Error() != nil {
				fmt.Println("链接错误: ", token.Error())
				return token.Error()
			}
		}
	}
	return nil
}

// Publish 发布消息,topic主题，data发送数据，qos发送等级
func (c *Client) Publish(topic string, data interface{}, qos QoS) error {
	if err := c.ensureConnected(); err != nil {
		return err
	}
	token := c.client.Publish(topic, byte(qos), false, data)
	if err := token.Error(); err != nil {
		return err
	}
	if !token.WaitTimeout(time.Second * 10) {
		return errors.New("MQTT 发送超时")
	}
	return nil
}

func (c *Client) LostHandler(lost func(c *Client, err error)) {
	c.onLost = lost
}

// Subscribe 订阅消息
func (c *Client) Subscribe(qos QoS, observer func(c *Client, payload []byte), topics ...string) error {
	if len(topics) == 0 {
		return errors.New("主题为空")
	}
	if observer == nil {
		return errors.New("处理函数为空")
	}
	if c.observer != nil {
		return errors.New("已存在订阅客户窜，需要重新订阅")
	}
	c.observer = observer
	filters := make(map[string]byte)
	for _, topic := range topics {
		filters[topic] = byte(qos)
	}
	c.client.SubscribeMultiple(filters, c.messageHandler)
	return nil
}

// Unsubscribe 取消订阅
func (c *Client) Unsubscribe(topics ...string) {
	c.observer = nil
	c.client.Unsubscribe(topics...)
}

func (c *Client) messageHandler(client mqtt.Client, msg mqtt.Message) {
	if c.observer == nil {
		fmt.Println("没有找到消息处理函数")
		return
	}
	c.observer(c, msg.Payload())
}

// Init 初始化mqtt
func Init() {
	slice := make([]*Config, 0)
	err := config.UnmarshalKey("mqtt", &slice)
	if err != nil || len(slice) == 0 {
		logx.Sugar.Warnf("未发现mqtt配置信息,  %v", err)
		return
	}
	for i, config := range slice {
		aliasName := config.Name
		if i == 0 && config.Name == "" {
			aliasName = "default"
		}
		RegisterMQTT(aliasName, config)
	}
}

// RegisterMQTT 注册mqtt
func RegisterMQTT(aliasName string, conf *Config) {
	if aliasName == "" {
		logx.Sugar.Error("MQTT别名为空")
		return
	}
	cli := NewClient(conf)
	if err := cli.Connect(); err != nil {
		panic(err)
	}
	instance[aliasName] = cli
}

func NewClient(conf *Config) *Client {
	protocol := "tcp"
	if conf.TLS != "" {
		protocol = "tls"
	}
	host := fmt.Sprintf("%s://%s:%d", protocol, conf.Host, conf.Port)
	fmt.Println(host, "password: ", conf.Password, "username: ", conf.User, "client id: ", conf.ClientID)
	options := mqtt.NewClientOptions().
		AddBroker(host).
		SetUsername(conf.User).
		SetPassword(conf.Password).
		SetClientID(conf.ClientID).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetKeepAlive(120 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetWriteTimeout(10 * time.Second).
		SetOnConnectHandler(func(client mqtt.Client) {
			logx.Sugar.Infof("MQTT 连接成功！client: %s - %s", conf.Name, conf.ClientID)
			reader := client.OptionsReader()
			for _, conn := range instance {
				reader2 := conn.client.OptionsReader()
				if reader.ClientID() == reader2.ClientID() {
					if conn.onConnect != nil {
						conn.onConnect(conn)
						return
					}
				}
			}
		}).
		SetConnectionLostHandler(func(client mqtt.Client, err error) {
			logx.Sugar.Infof("MQTT 关闭连接! error: %v", err)
			reader := client.OptionsReader()
			for _, conn := range instance {
				reader2 := conn.client.OptionsReader()
				if reader.ClientID() == reader2.ClientID() {
					if conn.onLost != nil {
						conn.onLost(conn, err)
						return
					}
				}
			}
		})
	if conf.TLS != "" {
		tlsConfig := NewTLSConfig(conf.TLS)
		options = options.SetTLSConfig(tlsConfig)
	}
	cli := mqtt.NewClient(options)
	return &Client{
		client:  cli,
		options: options,
		locker:  &sync.Mutex{},
	}
}

func NewTLSConfig(f string) *tls.Config {
	if f != "" {
		certPool := x509.NewCertPool()
		pemCerts, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Printf("TLS file error: %v", err)
			return nil
		}
		certPool.AppendCertsFromPEM(pemCerts)
		return &tls.Config{
			RootCAs:            certPool,
			InsecureSkipVerify: true,
			ClientAuth:         tls.NoClientCert,
			ClientCAs:          nil,
		}
	}
	return nil
}

// MQ 根据名称获取mq客户端
func MQ(name ...string) *Client {
	if len(name) == 0 {
		return instance["default"]
	}
	return instance[name[0]]
}
