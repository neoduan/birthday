package opentracing

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/neoduan/birthday/assembly"
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"strconv"
)

var _ assembly.IAssembly = &Assembly{}

var (
	defaultTraceCollectType = "http"
	defaultTraceCollectAddr = "127.0.0.1:9411"
	defaultTraceSampleRate  = "1.0" // <1.0 不采样  >=1.0全采样
)

type Assembly struct {
	prjName string
	tracer  opentracing.Tracer
}

func (this *Assembly) Setup() {
	this.setupUSI()
	this.setupSSI()

	this.setupUCI()
	this.setupSCI()
}

func (this *Assembly) Unload() {
	//nothing to do
}

func (this *Assembly) setupUSI() {
	assembly.USI = append(assembly.USI,
		grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(this.tracer)),
	)

}

func (this *Assembly) setupSSI() {
	assembly.SSI = append(assembly.SSI,
		grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(this.tracer)),
	)
}

func (this *Assembly) setupUCI() {
	assembly.UCI = append(assembly.UCI,
		grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(this.tracer)),
	)
}

func (this *Assembly) setupSCI() {
	assembly.SCI = append(assembly.SCI,
		grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(this.tracer)),
	)
}

func New(prjName string) *Assembly {
	var (
		tracer opentracing.Tracer
		err    error
	)

	if tracer, err = newTracer(prjName); err != nil {
		tracer = opentracing.NoopTracer{}
		log.Printf("[trace] init failed, err:%s.\n", err)
	}

	return &Assembly{
		prjName: prjName,
		tracer:  tracer,
	}
}

func newTracer(prjName string) (opentracing.Tracer, error) {
	var (
		traceCollectType = defaultTraceCollectType
		traceCollectAddr = defaultTraceCollectAddr
		traceSampleRate  = defaultTraceSampleRate
		collector        zipkin.Collector
		tracer           opentracing.Tracer
		err              error
	)

	if collectType := os.Getenv("ENV_TRACE_COLLECT_TYPE"); collectType != "" {
		traceCollectType = collectType
	}

	if collectAddr := os.Getenv("ENV_TRACE_COLLECT_ADDR"); collectAddr != "" {
		traceCollectAddr = collectAddr
	}

	if sampleRate := os.Getenv("ENV_TRACE_SAMPLE_RATE"); sampleRate != "" {
		traceSampleRate = sampleRate
	}

	switch traceCollectType {
	case "http":
		collector, err = zipkin.NewHTTPCollector(fmt.Sprintf("http://%s/api/v1/spans", traceCollectAddr))
	case "kafka":
		collector, err = zipkin.NewKafkaCollector(strings.Split(traceCollectAddr, ","))
	default:
		log.Panicf("[trace] collector type[%s] is illegal.\n", traceCollectType)
	}

	if err != nil {
		return tracer, err
	}

	rate, _ := strconv.ParseFloat(traceSampleRate, 32)
	tracer, err = zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, "0.0.0.0:0", prjName),
		zipkin.ClientServerSameSpan(true),
		zipkin.TraceID128Bit(true),
		zipkin.WithSampler(zipkin.NewCountingSampler(rate)),
	)

	return tracer, err
}
