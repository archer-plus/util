package util

import (
	"hash/crc32"
	"strings"
	"sync"

	"github.com/bwmarrin/snowflake"
	uuid "github.com/satori/go.uuid"
)

// NewUUID 创建UUID唯一编号
func NewUUID() string {
	u := uuid.NewV4()
	return strings.ReplaceAll(u.String(), "-", "")
}

// NewSequenceID 创建数字唯一编号
func NewSequenceID() int64 {
	id := int64(1)
	u := NewUUID()
	v := int(crc32.ChecksumIEEE([]byte(u)))
	if v >= 0 {
		id += int64(v)
	}
	if -v >= 0 {
		id += int64(-v)
	}
	return id
}

var (
	node *snowflake.Node
	err  error
	once sync.Once
)

// SnowIDInit 初始化雪花ID，n代表节点数，0～1023，不可以重复
func SnowIDInit(n int) error {
	once.Do(func() {
		snowflake.Epoch = 1601510400000
		node, err = snowflake.NewNode(int64(n))
		if err != nil {
			panic(err)
		}
	})
	return nil
}

// NewSnowID 新建一个雪花ID
func NewSnowID() snowflake.ID {
	return node.Generate()
}
