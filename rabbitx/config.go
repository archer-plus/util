package rabbitx

// Config 消息队列配置信息
type Config struct {
	Address    string `mapstructure:"address"`     // rabbitmq连接地址
	VHost      string `mapstructure:"vhost"`       // 虚拟路径
	Heartbeat  int    `mapstructure:"heartbeat"`   // 心跳间隔时间
	Exchange   string `mapstructure:"exchange"`    // 交换机名称
	Queue      string `mapstructure:"queue"`       // 队列名称
	RoutingKey string `mapstructure:"routing_key"` // 路由
}
