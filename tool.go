package util

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	crand "crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
)

func HSHA256(message string, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(message))
	return h.Sum(nil)
}

func HSHA1(message string, secret []byte) []byte {
	h := hmac.New(sha1.New, secret)
	h.Write([]byte(message))
	return h.Sum(nil)
}

func HMD5(message string, secret []byte) []byte {
	h := hmac.New(md5.New, secret)
	h.Write([]byte(message))
	return h.Sum(nil)
}

// NewMD516 MD5 16位 小写
func NewMD516(s string) string {
	return NewMD5(s)[8:24]
}

// NewMD5 MD5 32位 小写
func NewMD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	cipher := h.Sum(nil)
	return fmt.Sprintf("%x", cipher)
}

// NewBytesMD5 字节流MD5
func NewBytesMD5(data []byte) string {
	h := md5.New()
	h.Write(data)
	cipher := h.Sum(nil)
	return fmt.Sprintf("%x", cipher)
}

// Base64 编码
func Base64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// MilliSecondTimestamp 毫秒时间戳
func MilliSecondTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

// IntMax 最大int值
const IntMax = int(^uint(0) >> 1)

// RandomInt64MinMax 生成指定范围内的随机数字
func RandomInt64MinMax(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

// RandomInt64 生成指定位数的随机数字
func RandomInt64(b ...int) int64 {
	l := 5
	if len(b) > 0 && b[0] != 0 {
		l = b[0]
	}
	w := math.Pow(10, float64(l))
	n := int64(0)
	if w > math.MaxFloat64 {
		n = int64(math.MaxInt64)
	} else {
		n = int64(w)
	}
	result, _ := crand.Int(crand.Reader, big.NewInt(n))
	if l > len(result.String()) {
		return RandomInt64(l)
	}
	return result.Int64()
}

// MixMobile 混淆手机号
func MixMobile(mobile string) string {
	var phone string
	chars := strings.Split(mobile, "")
	if len(chars) < 7 {
		return mobile
	}
	for i := 0; i < len(chars); i++ {
		if i > 2 && i < 7 {
			phone += "*"
		} else {
			phone += chars[i]
		}
	}
	return phone
}

// ClientIP 获取客户端IP
func ClientIP(r *http.Request, c *websocket.Conn) string {
	ip := ""
	if r != nil {
		ip = r.Header.Get("X-Forwarded-For")
		ip = strings.TrimSpace(strings.Split(ip, ",")[0])
		if ip == "" {
			ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
		}
		if ip == "" && r.Host != "" {
			idx := strings.Index(r.Host, ":")
			ip = r.Host[:idx]
			if ip == "localhost" {
				ip = "127.0.0.1"
			}
		}
		if ip != "" {
			return ip
		}
	}
	if c != nil {
		if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.RemoteAddr().String())); err == nil {
			return ip
		}
		if ip == "" {
			return c.RemoteAddr().String()
		}
	}
	return ip
}

// Exists 判断路径文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return false
}

// IsDir 判断路径是否是文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断是否是文件
func IsFile(path string) bool {
	return !IsDir(path)
}

//HTTPGet get 请求
func HTTPGet(uri string) ([]byte, error) {
	response, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

//HTTPPost post 请求
func HTTPPost(uri string, data string) ([]byte, error) {
	body := bytes.NewBuffer([]byte(data))
	response, err := http.Post(uri, "", body)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

// GinContext 获取gin上下文
func GinContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("_gin.context_")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func ConvertOctalUtf8(in string) string {
	s := []byte(in)
	reg := regexp.MustCompile(`\\[0-7]{3}`)

	out := reg.ReplaceAllFunc(s,
		func(b []byte) []byte {
			i, _ := strconv.ParseInt(string(b[1:]), 8, 0)
			return []byte{byte(i)}
		})
	return string(out)
}

func IsHTTPS(ctx *gin.Context) bool {
	if ctx.GetHeader("X-Forwarded-Proto") == "https" || ctx.Request.TLS != nil {
		return true
	}
	return false
}

func IsMobile(userAgent string) bool {
	if len(userAgent) == 0 {
		return false
	}
	isMobile := false
	mobileKeys := []string{"Mobile", "Android", "Silk/", "Kindle", "BlackBerry", "Opera Mini", "Opera Mobi"}

	for i := 0; i < len(mobileKeys); i++ {
		if strings.Contains(userAgent, mobileKeys[i]) {
			isMobile = true
			break
		}
	}
	return isMobile
}
