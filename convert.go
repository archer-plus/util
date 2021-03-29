package util

import (
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// CheckDate 检查是否日期格式
func CheckDate(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	if err != nil {
		return false
	}
	return true
}

// CheckTime 检查是否时间格式
func CheckTime(s string) bool {
	_, err := ParseTime(s)
	if err != nil {
		return false
	}
	return true
}

// ParseTime 字符串转时间，默认格式 2006-01-02 15:04:05
func ParseTime(s string, layout ...string) (time.Time, error) {
	local, _ := time.LoadLocation("Local")
	if len(layout) > 0 {
		l := layout[0]
		t, err := time.ParseInLocation(l, s, local)
		return t, err
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, local)
	return t, err
}

// TimeFormat 将时间格式化为字符串
func TimeFormat(t time.Time, layout ...string) string {
	s := ""
	if t.IsZero() {
		return ""
	}
	if len(layout) == 0 {
		s = t.Format("2006-01-02 15:04:05")
	} else {
		s = t.Format(layout[0])
	}
	return s
}

// ItoSQLLike Any转Like字符串'%...%'
func ItoSQLLike(i interface{}) string {
	res := ItoS(i)
	return "'%" + res + "%'"
}

func transSQL(str string) string {
	str = strings.ReplaceAll(str, "\\", "\\\\")
	str = strings.ReplaceAll(str, "'", "\\'")
	return str
}

// ItoSQL Any转SQL字符串，字符带''，数字不带''
func ItoSQL(i interface{}) string {
	if i == nil {
		return "''"
	}
	res := "''"
	switch v := i.(type) {
	case []uint8:
		res = "'" + transSQL(string(v)) + "'"
	case *[]uint8:
		res = "'" + transSQL(string(*v)) + "'"
	case []int32:
		res = "'" + transSQL(string(v)) + "'"
	case *[]int32:
		res = "'" + transSQL(string(*v)) + "'"
	case string:
		res = "'" + transSQL(v) + "'"
	case *string:
		res = "'" + transSQL(*v) + "'"
	case int8:
		res = strconv.Itoa(int(v))
	case *int8:
		res = strconv.Itoa(int(*v))
	case uint8:
		res = strconv.Itoa(int(v))
	case *uint8:
		res = strconv.Itoa(int(*v))
	case int:
		res = strconv.Itoa(v)
	case *int:
		res = strconv.Itoa(*v)
	case uint:
		res = strconv.FormatUint(uint64(v), 10)
	case *uint:
		res = strconv.FormatUint(uint64(*v), 10)
	case int32:
		res = strconv.Itoa(int(v))
	case *int32:
		res = strconv.Itoa(int(*v))
	case uint32:
		res = strconv.FormatUint(uint64(v), 10)
	case *uint32:
		res = strconv.FormatUint(uint64(*v), 10)
	case int64:
		res = strconv.FormatInt(v, 10)
	case *int64:
		res = strconv.FormatInt(*v, 10)
	case uint64:
		res = strconv.FormatUint(v, 10)
	case *uint64:
		res = strconv.FormatUint(*v, 10)
	case float32:
		res = strconv.FormatFloat(float64(v), 'f', 2, 64)
	case *float32:
		res = strconv.FormatFloat(float64(*v), 'f', 2, 64)
	case float64:
		res = strconv.FormatFloat(v, 'f', 2, 64)
	case *float64:
		res = strconv.FormatFloat(*v, 'f', 2, 64)
	case bool:
		res = strconv.FormatBool(v)
	case *bool:
		res = strconv.FormatBool(*v)
	case time.Time:
		res = "'" + TimeFormat(v) + "'"
	case *time.Time:
		res = "'" + TimeFormat(*v) + "'"
	case struct{}, map[string]interface{}, interface{}, []interface{}:
		tmp, _ := jsoniter.Marshal(v)
		res = "'" + transSQL(string(tmp)) + "'"
	}
	return res
}

// ItoS Any转字符串
func ItoS(i interface{}, val ...string) string {
	if i == nil {
		return ""
	}
	res := ""
	switch v := i.(type) {
	case []uint8:
		res = string(v)
	case *[]uint8:
		res = string(*v)
	case []int32:
		res = string(v)
	case *[]int32:
		res = string(*v)
	case string:
		res = v
	case int8:
		res = strconv.Itoa(int(v))
	case *int8:
		res = strconv.Itoa(int(*v))
	case uint8:
		res = strconv.Itoa(int(v))
	case *uint8:
		res = strconv.Itoa(int(*v))
	case int:
		res = strconv.Itoa(v)
	case *int:
		res = strconv.Itoa(*v)
	case uint:
		res = strconv.FormatUint(uint64(v), 10)
	case *uint:
		res = strconv.FormatUint(uint64(*v), 10)
	case int32:
		res = strconv.Itoa(int(v))
	case *int32:
		res = strconv.Itoa(int(*v))
	case uint32:
		res = strconv.FormatUint(uint64(v), 10)
	case *uint32:
		res = strconv.FormatUint(uint64(*v), 10)
	case int64:
		res = strconv.FormatInt(v, 10)
	case *int64:
		res = strconv.FormatInt(*v, 10)
	case uint64:
		res = strconv.FormatUint(v, 10)
	case *uint64:
		res = strconv.FormatUint(*v, 10)
	case float32:
		res = strconv.FormatFloat(float64(v), 'f', 2, 64)
	case *float32:
		res = strconv.FormatFloat(float64(*v), 'f', 2, 64)
	case float64:
		res = strconv.FormatFloat(v, 'f', 2, 64)
	case *float64:
		res = strconv.FormatFloat(*v, 'f', 2, 64)
	case bool:
		res = strconv.FormatBool(v)
	case *bool:
		res = strconv.FormatBool(*v)
	case time.Time:
		res = TimeFormat(v)
	case *time.Time:
		res = TimeFormat(*v)
	case struct{}, map[string]interface{}, interface{}, []interface{}:
		tmp, _ := jsoniter.Marshal(v)
		res = string(tmp)
	}
	if res == "" && len(val) > 0 {
		return strings.TrimSpace(val[0])
	}
	return strings.TrimSpace(res)
}

// MtoS 获取map字符串值
func MtoS(src map[string]interface{}, key string, val ...string) string {
	res := ""
	if src != nil {
		if value, ok := src[key]; ok {
			res = ItoS(value, val...)
		}
	}
	return res
}

// ItoB Any转布尔值
func ItoB(i interface{}) bool {
	res := false
	switch v := i.(type) {
	case string:
		res, _ = strconv.ParseBool(v)
	case *string:
		res, _ = strconv.ParseBool(*v)
	case int8:
		res = v != int8(0)
	case *int8:
		res = *v != int8(0)
	case uint8:
		res = v != uint8(0)
	case *uint8:
		res = *v != uint8(0)
	case int:
		res = v != 0
	case *int:
		res = *v != 0
	case uint:
		res = v != uint(0)
	case *uint:
		res = *v != uint(0)
	case int32:
		res = v != int32(0)
	case *int32:
		res = *v != int32(0)
	case uint32:
		res = v != uint32(0)
	case *uint32:
		res = *v != uint32(0)
	case int64:
		res = v != int64(0)
	case *int64:
		res = *v != int64(0)
	case uint64:
		res = v != uint64(0)
	case *uint64:
		res = *v != uint64(0)
	case float32:
		res = v != float32(0)
	case *float32:
		res = *v != float32(0)
	case float64:
		res = v != float64(0)
	case *float64:
		res = *v != float64(0)
	case bool:
		res = v
	}
	return res
}

// MtoB 获取map布尔值
func MtoB(src map[string]interface{}, key string) bool {
	res := false
	if src != nil {
		if value, ok := src[key]; ok {
			res = ItoB(value)
		}
	}
	return res
}

// ItoF64 Any转浮点数
func ItoF64(i interface{}, val ...float64) float64 {
	res := float64(0)
	switch v := i.(type) {
	case int8:
		res = float64(v)
	case *int8:
		res = float64(*v)
	case uint8:
		res = float64(v)
	case *uint8:
		res = float64(*v)
	case int:
		res = float64(v)
	case *int:
		res = float64(*v)
	case uint:
		res = float64(v)
	case *uint:
		res = float64(*v)
	case int32:
		res = float64(v)
	case *int32:
		res = float64(*v)
	case uint32:
		res = float64(v)
	case *uint32:
		res = float64(*v)
	case int64:
		res = float64(v)
	case *int64:
		res = float64(*v)
	case uint64:
		res = float64(v)
	case *uint64:
		res = float64(*v)
	case float32:
		res = float64(v)
	case *float32:
		res = float64(*v)
	case float64:
		res = v
	case *float64:
		res = *v
	case string:
		res, _ = strconv.ParseFloat(v, 64)
	case *string:
		res, _ = strconv.ParseFloat(*v, 64)
	case bool, *bool:
		if true {
			res = float64(1)
		} else {
			res = float64(0)
		}
	}
	if res == float64(0) && len(val) > 0 {
		return val[0]
	}
	return res
}

// MtoF64 获取map浮点数值
func MtoF64(src map[string]interface{}, key string, val ...float64) float64 {
	res := float64(0)
	if src != nil {
		if value, ok := src[key]; ok {
			res = ItoF64(value, val...)
		}
	}
	return res
}

// ItoI64 Any转整数
func ItoI64(i interface{}, val ...int64) int64 {
	res := int64(0)
	switch v := i.(type) {
	case int8:
		res = int64(v)
	case *int8:
		res = int64(*v)
	case uint8:
		res = int64(v)
	case *uint8:
		res = int64(*v)
	case int:
		res = int64(v)
	case *int:
		res = int64(*v)
	case uint:
		res = int64(v)
	case *uint:
		res = int64(*v)
	case int32:
		res = int64(v)
	case *int32:
		res = int64(*v)
	case uint32:
		res = int64(v)
	case *uint32:
		res = int64(*v)
	case int64:
		res = v
	case *int64:
		res = *v
	case uint64:
		res = int64(v)
	case *uint64:
		res = int64(*v)
	case float32:
		res = int64(v)
	case *float32:
		res = int64(*v)
	case float64:
		res = int64(v)
	case *float64:
		res = int64(*v)
	case string:
		tmp, _ := strconv.Atoi(v)
		res = int64(tmp)
	case *string:
		tmp, _ := strconv.Atoi(*v)
		res = int64(tmp)
	case bool, *bool:
		if true {
			res = int64(1)
		} else {
			res = int64(0)
		}
	}
	if res == int64(0) && len(val) > 0 {
		return val[0]
	}
	return res
}

// MtoI64 获取map整数值64位
func MtoI64(src map[string]interface{}, key string, val ...int64) int64 {
	res := int64(0)
	if src != nil {
		if value, ok := src[key]; ok {
			res = ItoI64(value, val...)
		}
	}
	return res
}

// ItoI Any转整数
func ItoI(i interface{}, val ...int) int {
	res := 0
	switch v := i.(type) {
	case int8:
		res = int(v)
	case *int8:
		res = int(*v)
	case uint8:
		res = int(v)
	case *uint8:
		res = int(*v)
	case int:
		res = v
	case *int:
		res = *v
	case uint:
		res = int(v)
	case *uint:
		res = int(*v)
	case int32:
		res = int(v)
	case *int32:
		res = int(*v)
	case uint32:
		res = int(v)
	case *uint32:
		res = int(*v)
	case int64:
		res = int(v)
	case *int64:
		res = int(*v)
	case uint64:
		res = int(v)
	case *uint64:
		res = int(*v)
	case float32:
		res = int(v)
	case *float32:
		res = int(*v)
	case float64:
		res = int(v)
	case *float64:
		res = int(*v)
	case string:
		res, _ = strconv.Atoi(v)
	case *string:
		res, _ = strconv.Atoi(*v)
	case bool, *bool:
		if true {
			res = 1
		} else {
			res = 0
		}
	}
	if res == 0 && len(val) > 0 {
		return val[0]
	}
	return res
}

// MtoI 获取map整数值
func MtoI(src map[string]interface{}, key string, val ...int) int {
	res := 0
	if src != nil {
		if value, ok := src[key]; ok {
			res = ItoI(value, val...)
		}
	}
	return res
}
