package log

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/neoduan/birthday/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"sync"
)

type logHelper struct {
	prjName   string      //工程的名称
	logPath   string      //日志路径名
	logName   string      //日志文件名
	logFile   string      //日志文件
	logPrefix string      //日志前缀
	logKit    *zap.Logger //zap日志对象
	config    zap.Config
	logLevel  zapcore.Level //zap日志级别
	sync.Mutex
}

func NewLogHelper(prj, prefix string) *logHelper {
	return &logHelper{
		prjName:   prj,
		logPrefix: prefix,
	}
}

func (this *logHelper) SetLogLevel(level string) {
	switch level {
	case "debug":
		this.logLevel = zapcore.DebugLevel
	case "info":
		this.logLevel = zapcore.InfoLevel
	case "warn":
		this.logLevel = zapcore.WarnLevel
	case "error":
		this.logLevel = zapcore.ErrorLevel
	default:
		log.Panicf("[logger] level[%s] is illegal.", level)
	}

	return
}

func (this *logHelper) GetLogKit() *zap.Logger {
	return this.logKit
}

func (this *logHelper) Execute() {
	if exists := this.islogFileExist(); !exists {
		this.createLogFile()
	}

	this.loadKit()
}

func (this *logHelper) Cancel() {
	if Kit != nil {
		Kit.Sync()
	}
}

func (this *logHelper) LogLevelSwitch() {
	this.Lock() //锁过程,保证kit加载不竞争
	defer this.Unlock()

	if this.logLevel == zapcore.DebugLevel {
		this.logLevel = zapcore.InfoLevel
	} else {
		this.logLevel = zapcore.DebugLevel
	}

	this.config.Level.SetLevel(this.logLevel)
}

func (this *logHelper) loadKit() {
	cfg := zap.NewProductionConfig()

	cfg.OutputPaths = []string{this.logFile}
	cfg.Level = zap.NewAtomicLevelAt(this.logLevel)
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		log.Panicf("[logger] config build failed, err:%s.\n", err)
	}

	grpc_zap.ReplaceGrpcLogger(logger)
	this.logKit = logger
	this.config = cfg
	return
}

//func (this *logHelper) work() {
//	exists, err := this.islogFileExist()
//	if err != nil {
//		log.Printf("[logger] whether logFile[%s] exists or not, err:%s.\n", this.logFile, err)
//		return
//	}
//
//	if !exists {
//		this.createLogFile()
//
//		this.Lock()
//		this.loadKit()
//		this.Unlock()
//	}
//}

func (this *logHelper) islogFileExist() bool {
	var (
		err    error
		exists bool
	)

	logPath := fmt.Sprintf("%s/%s/", this.logPrefix, this.prjName)
	logName := fmt.Sprintf("%s.log", this.prjName)
	logFile := logPath + logName

	exists, err = utils.PathExists(logFile)
	if err != nil {
		log.Panicf("[logger] logFile[%s] exists, err:%s.\n", logFile, err)
	}

	this.logPath = logPath
	this.logName = logName
	this.logFile = logFile

	return exists
}

func (this *logHelper) createLogFile() {
	err := os.MkdirAll(this.logPath, 0755)
	if err != nil {
		log.Panicf("[logger] create logPath[%s] failed, err:%s.\n", this.logPath, err)
	}
	_, err = os.Create(this.logFile)
	if err != nil {
		log.Panicf("[logger] create logFile[%s] failed, err:%s.\n", this.logFile, err)
	}
}
