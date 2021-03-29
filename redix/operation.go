package redix

import (
	"fmt"
	"time"
)

type empty struct {
}

const (
	// EXPIRE 过期时间
	EXPIRE = "expire"
	// NX 如果不存在，则 SET
	NX = "nx"
	// XX 如果存在，则 SET
	XX = "xx"
)

// Operation 操作配置
type Operation struct {
	Name  string
	Value interface{}
}

// Operations 操作数组
type Operations []*Operation

// Find 根据名称返回结果
func (c Operations) Find(name string) *Result {
	for _, attr := range c {
		if attr.Name == name {
			return NewResult(attr.Value, nil)
		}
	}
	return NewResult(nil, fmt.Errorf("operation found error: %s", name))
}

// WithExpire 具备超时
func WithExpire(t int) *Operation {
	return &Operation{Name: EXPIRE, Value: time.Duration(t) * time.Second}
}

// WithNX 具备不存在写入
func WithNX() *Operation {
	return &Operation{Name: NX, Value: empty{}}
}

// WithXX 具备存在写入
func WithXX() *Operation {
	return &Operation{Name: XX, Value: empty{}}
}
