package log

import (
	"os"
	"os/signal"
	"syscall"

	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/context"
)

const (
	defaultLogPath  = "/data/project/"
	defaultPrjName  = "unknown"
	defaultLogLevel = "info"
)

type (
	Fields map[string]interface{}
	Entry  struct {
		data Fields
		kit  *zap.Logger
	}
)

var Kit *zap.Logger = nil

func Stop() {
	if Kit != nil {
		Kit.Sync()
	}
}

func WithField(fields Fields) *Entry {
	return &Entry{data: fields, kit: Kit}

}

func WithoutField() *Entry {

	return &Entry{kit: Kit}

}

func WithTrace(ctx context.Context) *Entry {
	return &Entry{kit: grpc_zap.Extract(ctx)}
}

func (this *Entry) Debug(msg string, fields ...Fields) {
	this.kit.Debug(msg, this.zapField(fields...)...)
}

func (this *Entry) Info(msg string, fields ...Fields) {
	this.kit.Info(msg, this.zapField(fields...)...)
}

func (this *Entry) Warn(msg string, fields ...Fields) {
	this.kit.Warn(msg, this.zapField(fields...)...)
}

func (this *Entry) Error(msg string, fields ...Fields) {
	this.kit.Error(msg, this.zapField(fields...)...)
}

func (this *Entry) Panic(msg string, fields ...Fields) {
	this.kit.Panic(msg, this.zapField(fields...)...)
}

func (this *Entry) Fatal(msg string, fields ...Fields) {
	this.kit.Fatal(msg, this.zapField(fields...)...)
}

func (this *Entry) zapField(fields ...Fields) []zapcore.Field {
	var (
		zapLogSlice      []zapcore.Field
		key              string
		value            interface{}
		zapLogSliceIndex = 0
		fieldsSize       = len(fields)
	)

	if fieldsSize > 0 {
		zapLogSlice = make([]zapcore.Field, len(this.data)+len(fields[0]))
		for key, value = range this.data {
			zapLogSlice[zapLogSliceIndex] = zap.Any(key, value)
			zapLogSliceIndex++
		}

		for key, value = range fields[0] {
			zapLogSlice[zapLogSliceIndex] = zap.Any(key, value)
			zapLogSliceIndex++
		}
	} else {
		zapLogSlice = make([]zapcore.Field, len(this.data))
		for key, value = range this.data {
			zapLogSlice[zapLogSliceIndex] = zap.Any(key, value)
			zapLogSliceIndex++
		}
	}

	return zapLogSlice
}

func init() {
	var (
		envLogPath  = defaultLogPath //log路径使用默认的,不从环境变量里取了
		envPrjName  = defaultPrjName
		envLogLevel = defaultLogLevel
	)

	if name := os.Getenv("ENV_PRJ_NAME"); name != "" {
		envPrjName = name
	}

	if level := os.Getenv("ENV_LOG_LEVEL"); level != "" {
		envLogLevel = level
	}

	helper := NewLogHelper(envPrjName, envLogPath)
	helper.SetLogLevel(envLogLevel)
	helper.Execute()

	Kit = helper.GetLogKit()

	go func() {
		for {
			c := make(chan os.Signal)
			signal.Notify(c, syscall.SIGUSR1)
			<-c
			fmt.Printf("[DEBUG] tag:%d, msg:%+v.\n", 111111, "swithch")
			helper.LogLevelSwitch()
		}
	}()
}
