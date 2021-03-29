package rabbitx

import (
	"time"

	"github.com/streadway/amqp"
)

// Produce 生产消息
func (r *RabbitMQ) Produce(exchange, key string, mandatory, immediate bool, data []byte) error {
	ch, err := r.Channel(5 * time.Second)
	if err != nil {
		return err
	}
	return ch.Publish(exchange, key, mandatory, immediate, amqp.Publishing{
		ContentType: "text/plain",
		Body:        data,
	})
}
