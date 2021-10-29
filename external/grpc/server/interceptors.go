package server

import (
	"context"
	"errors"
	"reflect"

	"google.golang.org/grpc"

	pb "github.com/dysnix/ai-scale-proto/external/proto/services"
)

func InjectClientMetadataInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		st := reflect.TypeOf(req)
		_, ok := st.MethodByName("GetHeader")
		if ok {
			var b interface{} = pb.Header{ClusterId: "bsc-1"}
			field := reflect.New(reflect.TypeOf(b))
			field.Elem().Set(reflect.ValueOf(b))
			reflect.ValueOf(&req).Elem().FieldByName("Header").Set(field)
		}

		resp, err = handler(ctx, req)

		return resp, err
	}
}

func PanicServerInterceptor(panicHandler func(ctx context.Context, err error, params ...interface{}) error, params ...interface{}) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		panicked := true

		defer func() {
			if r := recover(); r != nil || panicked {
				switch errBody := r.(type) {
				case error:
					err = panicHandler(ctx, errBody, params...)
				case string:
					err = panicHandler(ctx, errors.New(errBody), params...)
				}
			}
		}()

		resp, err = handler(ctx, req)

		panicked = false
		return resp, err
	}
}
