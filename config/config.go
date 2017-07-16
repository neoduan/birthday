package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	defaultConfigFileName = "config.yaml"
	defaultConfigFilePath = "./"
	defaultEnvPrefix      = "ENV"
)

// TODO: add others
//这里的envCfg结构体，会在Parse时候，被导入环境变量.
//这样每个模块通过环境变量就可以拿到关心的信息了.
var envCfg = struct {
	Prj_Name string `yaml:"_Prj_Name"`

	Log_Level string `yaml:"_Log_Level"`

	Trace_Collect_Type string `yaml:"_Trace_Collect_Type"`
	Trace_Collect_Addr string `yaml:"_Trace_Collect_Addr"`
	Trace_Sample_Rate  string `yaml:"_Trace_Sample_Rate"`
}{}

//配置文件默认是在当前路径下的config.yaml
//当然，为了方便，我们可以设置ENV_CONFIG_FILE环境变量来加载一个配置文件
func Parse(cfg interface{}) {
	var (
		envConfigFile  string
		configFilePath string = defaultConfigFilePath
		configFileName string = defaultConfigFileName
	)

	if envConfigFile = os.Getenv("ENV_CONFIG_FILE"); envConfigFile != "" {
		index := strings.LastIndexByte(envConfigFile, '/')
		configFilePath = envConfigFile[:index+1]
		configFileName = envConfigFile[index+1:]
	}

	configure := NewConfigure(configFilePath, configFileName)
	configure.LoadData()
	configure.ParseCfg(cfg)

	//将解析之后的envCfg导出到环境变量
	{
		configure.ParseCfg(&envCfg)
		export(envCfg)
	}
	return
}

type Configure struct {
	filePath string
	fileName string
	content  []byte
}

func NewConfigure(filePath, fileName string) *Configure {
	return &Configure{
		filePath: filePath,
		fileName: fileName,
	}
}

func (this *Configure) LoadData() {
	var (
		file = this.filePath + this.fileName
		err  error
		data []byte
	)

	data, err = ioutil.ReadFile(file)
	if err != nil {
		log.Panicf("[config] read file[%s] failed, errmsg:[%s].\n", file, err)
	}

	this.content = data
	return
}

func (this *Configure) ParseCfg(cfg interface{}) {
	err := yaml.Unmarshal(this.content, cfg)
	if err != nil {
		log.Panicf("[config] parse failed, errmsg:[%s].\n", err)
	}
}

func export(cfg interface{}) {
	var (
		envKey, envVal string
		err            error
	)

	t := reflect.TypeOf(cfg)
	v := reflect.ValueOf(cfg)

	for i := 0; i < t.NumField(); i++ {
		envKey = fmt.Sprintf("%s_%s", defaultEnvPrefix, strings.ToUpper(t.Field(i).Name))
		envVal = v.Field(i).String()

		if err = os.Setenv(envKey, envVal); err != nil {
			log.Panicf("[config] export failed, errmsg:[%s].\n", err)
		}
	}
}
