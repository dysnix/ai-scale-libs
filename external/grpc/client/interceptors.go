package client

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/dysnix/ai-scale-libs/external/configs"
)

const (
	clusterID = "cluster_id"
)

func InjectClientMetadataInterceptor(conf configs.Client) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) (err error) {
		return invoker(metadata.AppendToOutgoingContext(ctx, clusterID, conf.ClusterID), method, req, reply, cc, opts...)
	}
}

func PanicClientInterceptor(handler func(ctx context.Context, err error, params ...interface{}) error, params ...interface{}) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) (err error) {

		defer func() {
			if r := recover(); r != nil {
				switch errType := r.(type) {
				case error:
					err = handler(ctx, errType, params...)
				case string:
					err = handler(ctx, errors.New(errType), params...)
				}
			}
		}()

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
