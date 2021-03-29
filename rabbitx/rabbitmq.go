package rabbitx

import (
	"errors"
	"time"

	"github.com/archer-plus/util/logx"

	"github.com/streadway/amqp"
)

const (
	reconnectDelay    = 5 * time.Second // 断线重连时间
	reconnetDetectDur = 5 * time.Second
)

var (
	errNotConnected  = errors.New("不能连接到AMQP服务")
	errAlreadyClosed = errors.New("已经关闭: 不能连接到AMQP服务")
)

// Apply 应用方法
type Apply func(ch *amqp.Channel) error

// RabbitMQ 消息队列实例
type RabbitMQ struct {
	Address     string      // 连接地址
	Config      amqp.Config // 连接配置
	apply       Apply
	connection  *amqp.Connection
	channel     *amqp.Channel
	done        chan bool
	changeConn  chan struct{}
	chanNotify  chan *amqp.Error // 通道错误通知
	connNotify  chan *amqp.Error // 连接错误通知
	isConnected bool             // 是否连接
	hasConsumer bool             // 是否存在消费者
}

// 连接RabbitMQ
func (r *RabbitMQ) connect() (bool, error) {
	conn, err := amqp.DialConfig(r.Address, r.Config)
	if err != nil {
		return false, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return false, err
	}
	if err := r.apply(ch); err != nil {
		return false, err
	}
	r.isConnected = true
	r.changeConnection(conn, ch)
	return true, nil
}

// 监听RabbitMQ 通道状态
func (r *RabbitMQ) changeConnection(connection *amqp.Connection, channel *amqp.Channel) {
	r.connection = connection
	r.connNotify = make(chan *amqp.Error, 1)
	r.connection.NotifyClose(r.connNotify)

	r.channel = channel
	r.chanNotify = make(chan *amqp.Error, 1)
	r.channel.NotifyClose(r.chanNotify)

	if r.hasConsumer {
		r.changeConn <- struct{}{}
	}
	logx.Sugar.Info("rabbitmq 连接成功")
}

// Close 关闭连接
func (r *RabbitMQ) Close() error {
	if !r.isConnected {
		return errAlreadyClosed
	}
	err := r.channel.Close()
	if err != nil {
		return err
	}
	err = r.connection.Close()
	if err != nil {
		return err
	}
	close(r.done)
	r.isConnected = false
	return nil
}

// 断线重连
func (r *RabbitMQ) handleReconnect() {
	for {
		if !r.isConnected {
			logx.Sugar.Warn("尝试重连")
			var (
				connected = false
				err       error
			)

			for i := 0; !connected; i++ {
				if connected, err = r.connect(); err != nil {
					logx.Sugar.Errorf("重连失败: %v \n", err)
				}
				if !connected {
					logx.Sugar.Warnf("重新连接... %d \n", i)
				}
				time.Sleep(reconnectDelay)
			}

			select {
			case <-r.done:
				return
			case err := <-r.chanNotify:
				logx.Sugar.Errorf("通道关闭通知: %v \n", err)
				r.isConnected = false
			case err := <-r.connNotify:
				logx.Sugar.Errorf("连接关闭通知：%v \n", err)
				r.isConnected = false
			}
			time.Sleep(reconnetDetectDur)
		}
	}
}

// Channel 获取通道
func (r *RabbitMQ) Channel(timeout time.Duration) (*amqp.Channel, error) {
	timer := time.NewTimer(timeout)
	for !r.isConnected {
		select {
		case <-timer.C:
			return nil, errNotConnected
		default:
			time.Sleep(100 * time.Microsecond)
		}
	}
	return r.channel, nil
}

// New 创建RabbitMQ
func New(addr string, cfg amqp.Config, f Apply) *RabbitMQ {
	r := &RabbitMQ{
		Address:    addr,
		Config:     cfg,
		apply:      f,
		changeConn: make(chan struct{}, 1),
	}
	go r.handleReconnect()
	return r
}
