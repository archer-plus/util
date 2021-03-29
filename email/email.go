package email

import (
	"fmt"
	"net/smtp"

	"github.com/archer-plus/util/config"
	"github.com/archer-plus/util/logx"
)

// Config email服务器配置
type Config struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Nickname string `mapstructure:"nickname"`
}

var conf = Config{}

// Init 初始化邮件配置
func Init() {
	err := config.UnmarshalKey("smtp", &conf)
	if err != nil {
		logx.Sugar.Warnf("未发现mongo配置信息, %v", err)
		return
	}
}

// Send 发送邮件
func Send(address []string, subject string, body string) (err error) {
	auth := smtp.PlainAuth("", conf.User, conf.Password, conf.Host)
	contentType := "Content-Type: text/html; charset=UTF-8"
	for _, v := range address {
		s := fmt.Sprintf("To:%s\r\nFrom:%s<%s>\r\nSubject:%s\r\n%s\r\n\r\n%s",
			v, conf.Nickname, conf.User, subject, contentType, body)
		msg := []byte(s)
		addr := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
		err = smtp.SendMail(addr, auth, conf.User, []string{v}, msg)
		if err != nil {
			return err
		}
	}
	return nil
}
