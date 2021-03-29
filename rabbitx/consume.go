package rabbitx

import (
	"time"

	"github.com/archer-plus/util/logx"

	"github.com/streadway/amqp"
)

// ConsumeHandler 消费者监听
type ConsumeHandler func(d amqp.Delivery) error

// Consume 消费方法
func (r *RabbitMQ) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table, handler ConsumeHandler) {
	var (
		ch       *amqp.Channel
		delivery <-chan amqp.Delivery
		err      error
	)

	r.hasConsumer = true
	for {
		select {
		case <-r.changeConn:
			if ch, err = r.Channel(5 * time.Second); err != nil {
				logx.Sugar.Errorf("获取通道错误：%v", err)
				break
			}
			if delivery, err = ch.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args); err != nil {
				logx.Sugar.Errorf("订阅失败：%v", err)
				break
			}
		default:
			if !r.isConnected || delivery == nil {
				time.Sleep(1 * time.Second)
				break
			}

			for d := range delivery {
				if err := handler(d); err != nil {
					logx.Sugar.Errorf("订阅消息错误：%v -- %v", d, err)
				}
			}
		}
	}
}
