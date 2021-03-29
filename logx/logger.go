package logx

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/archer-plus/util/config"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var instance *Logger
var once sync.Once
var core zapcore.Core
var logger *zap.Logger

// Config 日志配置信息
type Config struct {
	FileName   string `mapstructure:"file_name"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	Level      string `mapstructure:"level"`
	SizeIncise bool   `mapstructrue:"size_incise"`
}

// Sugar 快捷方式
var Sugar *zap.SugaredLogger

// Logger 封装日志对象
type Logger struct {
	Sugar *zap.SugaredLogger
}

// New 创建日志对象，单例
func New() *Logger {
	return instance
}

// Log 实现 gokit 日志接口
func (l *Logger) Log(info ...interface{}) error {
	//l.Sugar.Info(info)
	fmt.Println(info)
	return nil
}

// Handle grpc error handler
func (l *Logger) Handle(ctx context.Context, err error) {
	l.Sugar.Warn(err)
}

func encoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		FunctionKey:    zapcore.OmitKey,
	}
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func getLevel(level string) zapcore.Level {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	}
	return zapcore.DebugLevel
}

// Init 初始化日志
func Init() {
	once.Do(func() {
		instance = &Logger{}
		conf := Config{}
		err := config.UnmarshalKey("log", &conf)
		if err != nil {
			return
		}
		if conf.MaxSize == 0 {
			conf.MaxSize = 100
		}
		if strings.TrimSpace(conf.FileName) == "" {
			conf.FileName = "logs/log"
		} else {
			idx := strings.LastIndex(conf.FileName, "/")
			if idx+1 <= utf8.RuneCountInString(conf.FileName) {
				suf := conf.FileName[idx+1:]
				if strings.TrimSpace(suf) == "" {
					conf.FileName += "log"
				}
			}
		}
		if strings.TrimSpace(conf.Level) == "" {
			conf.Level = "debug"
		}
		if conf.SizeIncise {
			w := zapcore.AddSync(&lumberjack.Logger{
				Filename:   conf.FileName,
				MaxSize:    conf.MaxSize,
				MaxBackups: 5,
				MaxAge:     conf.MaxAge,
				Compress:   false,
			})
			core = zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoderConfig()),
				zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), w),
				getLevel(conf.Level),
			)
			logger = zap.New(core, zap.AddCaller())
			instance.Sugar = logger.Sugar()
			Sugar = instance.Sugar
		} else {
			logLevel := getLevel(conf.Level)
			infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl < zapcore.WarnLevel && lvl >= logLevel
			})
			warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.WarnLevel && lvl <= logLevel
			})
			infoWriter := getWriter(conf.FileName)
			warnWriter := getWriter(conf.FileName + ".error")
			encoder := zapcore.NewConsoleEncoder(encoderConfig())
			core = zapcore.NewTee(
				zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
				zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
				zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), getLevel(conf.Level)),
			)
			logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
			instance.Sugar = logger.Sugar()
			Sugar = instance.Sugar
		}
		fmt.Println("初始化日志成功...")
	})

}

func getWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		filename+".%Y%m%d",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(err)
	}
	return hook
}
