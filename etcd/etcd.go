package etcd

import (
	"context"
	"sync"
	"time"

	"github.com/archer-plus/util/logx"

	"github.com/go-kit/kit/sd/etcdv3"
)

var client etcdv3.Client
var svc etcdv3.Service
var err error
var once sync.Once

// Init 初始化etcd,
func Init(address []string, username, password string, args ...int) {
	once.Do(func() {
		timeout := 3
		keepalive := 3
		l := len(args)
		if l > 0 {
			if args[0] != 0 {
				timeout = args[0]
			}
		}
		if l > 1 {
			if args[1] != 0 {
				timeout = args[1]
			}
		}
		client, err = etcdv3.NewClient(context.Background(), address, etcdv3.ClientOptions{
			DialTimeout:   time.Second * time.Duration(timeout),
			DialKeepAlive: time.Second * time.Duration(keepalive),
			Username:      username,
			Password:      password,
		})
		if err != nil {
			logx.Sugar.Panicf("etcd error %v", err)
		}
	})
}

// Register 注册服务
func Register(serviceName, address string) {
	if serviceName != "" {
		if serviceName[0] != '/' {
			serviceName = "/" + serviceName
		}
	}
	idx := len(serviceName)
	if serviceName[idx-1] != '/' {
		serviceName += "/"
	}
	svc = etcdv3.Service{
		Key:   serviceName + address,
		Value: address,
	}
	registrar := etcdv3.NewRegistrar(client, svc, logx.New())

	registrar.Register()
	logx.Sugar.Info("etcd 注册成功")
}

// Client 获取etcd客户端
func Client() *etcdv3.Client {
	return &client
}

// Deregister 反注册
func Deregister() {
	if client != nil {
		client.Deregister(svc)
	}
}
