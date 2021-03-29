package mqtx

import (
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// DecodeMessage 解码消息
func DecodeMessage(payload []byte) (*map[string]interface{}, error) {
	message := make(map[string]interface{}, 0)
	decoder := jsoniter.NewDecoder(strings.NewReader(string(payload)))
	decoder.UseNumber()
	if err := decoder.Decode(&message); err != nil {
		return nil, err
	}
	return &message, nil
}
