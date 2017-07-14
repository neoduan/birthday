package in_out_logger


import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"fmt"
)


func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Printf("[DEBUG] tag:%d, msg:%+v.\n", 1, "")
		ctx = context.WithValue(ctx, "test_key_1", "test_value_1")
		resp, err := handler(ctx, req)
		fmt.Printf("[DEBUG] tag:%d, msg:%+v.\n", 2, "")


		extCtx := grpc_zap.Extract(ctx)

		ccc:= ctx.Value("test_key_1")
		fmt.Printf("[DEBUG] tag:%d, msg:%+v.\n", 10000, ccc )

		extCtx.Debug("debg for logging request and response",
			zap.Reflect("req" , req),
			zap.Reflect("rsp" , resp),
			zap.Reflect("err" , err),
		)

		return resp, err
	}
}

