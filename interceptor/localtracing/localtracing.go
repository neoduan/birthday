package localtracing

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"reflect"
)

const traceFlag = "TraceId"

func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateOptions(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var (
			extCtxTags        = grpc_ctxtags.Extract(ctx)
			extCtxVals        = extCtxTags.Values()
			traceId    string = ""
		)

		if openTraceId, exists := extCtxVals[grpc_opentracing.TagTraceId]; !exists {
			traceId = o.idGenerateFunc()
			extCtxTags.Set(grpc_opentracing.TagTraceId, traceId)
		} else {
			if id, ok := openTraceId.(string); ok {
				traceId = id
			}
		}

		resp, err := handler(ctx, req)

		if val := reflect.ValueOf(resp).Elem().FieldByName(traceFlag); val.CanSet(){
			val.SetString(traceId)
		}

		return resp, err
	}
}
