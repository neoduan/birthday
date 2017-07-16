package log

/*

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/neoduan/birthday/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const defaultLogCheckInterval = 1

type logHelper struct {
	prjName   string        //工程的名称
	logPath   string        //日志路径名
	logName   string        //日志文件名
	logFile   string        //日志文件
	logPrefix string        //日志前缀
	interval  int           //日志检测间隔
	hasInit   bool          //zap是否已初始化
	logKit    *zap.Logger   //zap日志对象
	logLevel  zapcore.Level //zap日志级别
	sync.Mutex
}

func NewLogHelper(prj, prefix string) *logHelper {
	return &logHelper{
		prjName:   prj,
		logPrefix: prefix,
		hasInit:   false,
		interval:  defaultLogCheckInterval,
	}
}

func (this *logHelper) SetLogInvl(interval int) {
	this.interval = interval
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

//func (this *logHelper) GetLogKit() *zap.Logger {
//	return this.logKit
//}

func (this *logHelper) Execute() {
	if this.hasInit {
		return
	}

	if exists, _ := this.islogFileExist(); !exists {
		this.createLogFile()
	}

	this.loadKit()
	this.hasInit = true

	go this.background()
}

func (this *logHelper) Cancel() {
	if Kit != nil {
		Kit.Sync()
	}
}

func (this *logHelper) LogLevelSwitch() {
	this.Lock() //锁过程,保证kit加载不竞争
	defer this.Unlock()

	fmt.Printf("[DEBUG] tag:%d, msg:%+v.\n", 1111110, this.logLevel)
	if this.logLevel == zapcore.DebugLevel {
		this.logLevel = zapcore.InfoLevel
	} else {
		this.logLevel = zapcore.DebugLevel
	}
	fmt.Printf("[DEBUG] tag:%d, msg:%+v.\n", 11111100001, this.logLevel)

	fmt.Printf("[DEBUG] tag:%d, msg:%+v.\n", 8989, Kit)
	this.loadKit()
	fmt.Printf("[DEBUG] tag:%d, msg:%+v.\n", 899999, Kit)
}

func (this *logHelper) loadKit() {
	cfg := zap.NewProductionConfig()

	cfg.OutputPaths = []string{this.logFile}

	cfg.Level = zap.NewAtomicLevelAt(this.logLevel)
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		if this.hasInit == false {
			log.Panicf("[logger] config build failed, err:%s.\n", err)
		} else {
			return
		}
	}

	this.Cancel()
	grpc_zap.ReplaceGrpcLogger(logger)
	//this.logKit = logger
	Kit = logger

	return
}

func (this *logHelper) background() {
	tick := time.Tick(time.Duration(this.interval) * time.Second)
	for {
		select {
		case <-tick:
			this.work()
		}
	}
}

func (this *logHelper) work() {
	exists, err := this.islogFileExist()
	if err != nil {
		log.Printf("[logger] whether logFile[%s] exists or not, err:%s.\n", this.logFile, err)
		return
	}

	if !exists {
		this.createLogFile()

		this.Lock()
		this.loadKit()
		this.Unlock()
	}
}

func (this *logHelper) islogFileExist() (bool, error) {
	var (
		err                    error
		exists                 bool
		now                    = time.Now()
		year, month, day, hour = now.Year(), now.Month(), now.Day(), now.Hour()
	)

	logPath := fmt.Sprintf("%s/%s/%d%02d/%02d/", this.logPrefix, this.prjName, year, month, day)
	logName := fmt.Sprintf("%02d.log", hour)
	logFile := logPath + logName

	exists, err = utils.PathExists(logFile)
	if err != nil && !this.hasInit {
		log.Panicf("[logger] whether logFile[%s] exists or not, err:%s.\n", logFile, err)
	}

	if err == nil {
		this.logPath = logPath
		this.logName = logName
		this.logFile = logFile
	}

	return exists, err
}

func (this *logHelper) createLogFile() error {
	err := os.MkdirAll(this.logPath, 0755)
	if err != nil && !this.hasInit {
		log.Panicf("[logger] create logPath[%s] failed, err:%s.\n", this.logPath, err)
	}
	_, err = os.Create(this.logFile)
	if err != nil && !this.hasInit {
		log.Panicf("[logger] create logFile[%s] failed, err:%s.\n", this.logFile, err)
	}

	return err

}
*/
