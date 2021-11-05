package server

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	grpcC "github.com/dysnix/ai-scale-libs/external/grpc"
	pb "github.com/dysnix/ai-scale-proto/external/proto/services"
)

const (
	clusterID = "cluster_id"
)

func InjectClientMetadataInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		st := reflect.TypeOf(req)
		_, ok := st.MethodByName("GetHeader")
		if ok {

			header := pb.Header{}

			md, ok := metadata.FromIncomingContext(ctx)
			if ok {
				for key, val := range md {
					if strings.Contains(key, grpcC.ClusterIDKey) && len(val) > 0 {
						header.ClusterId = val[0]
						break
					}
				}

				if len(header.GetClusterId()) == 0 {
					for key, val := range md {
						if strings.Contains(key, grpcC.NameKey) && len(val) > 0 {
							header.ClusterId = val[0]
							break
						}
					}
				}
			}

			var b interface{} = header
			field := reflect.New(reflect.TypeOf(b))
			field.Elem().Set(reflect.ValueOf(b))
			reflect.ValueOf(req).Elem().FieldByName("Header").Set(field)
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
