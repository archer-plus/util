package mqtx

// Config mqtt配置
type Config struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	ClientID string `mapstructure:"client_id"`
	TLS      string `mapstructure:"tls"`
}
