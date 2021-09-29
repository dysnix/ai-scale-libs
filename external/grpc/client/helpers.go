package client

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	"github.com/dysnix/ai-scale-libs/external/configs"
	"github.com/dysnix/ai-scale-libs/external/grpc/zstd_compressor"
	_ "github.com/dysnix/ai-scale-libs/external/grpc/zstd_compressor"
)

const (
	DefaultMaxMsgSize = 2 << 20 // 2Mb
)

func SetGrpcClientOptions(conf *configs.GRPC, internalInterceptors ...grpc.UnaryClientInterceptor) (options []grpc.DialOption, err error) {
	unaryClientInterceptors := make([]grpc.UnaryClientInterceptor, 0)

	options = []grpc.DialOption{
		//grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(zstd_compressor.Name)),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                conf.Keepalive.Time,
				Timeout:             conf.Keepalive.Timeout,
				PermitWithoutStream: conf.Keepalive.EnforcementPolicy.PermitWithoutStream,
			},
		),
	}

	if conf.Conn.ReadBufferSize > 0 {
		options = append(options, grpc.WithReadBufferSize(int(conf.Conn.ReadBufferSize)))
	}

	if conf.Conn.WriteBufferSize > 0 {
		options = append(options, grpc.WithWriteBufferSize(int(conf.Conn.WriteBufferSize)))
	}

	if conf.Conn.MaxMessageSize > 0 {
		options = append(options, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(conf.Conn.MaxMessageSize))))
	} else {
		options = append(options, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(DefaultMaxMsgSize)))
	}

	if conf.Conn.Insecure {
		options = append(options, grpc.WithInsecure())
	}

	// TODO: implement all needed interceptors...

	unaryClientInterceptors = append(unaryClientInterceptors,
		PanicClientInterceptor(func(ctx context.Context, err error, params ...interface{}) error {
			//TODO:? can be any other logic...
			return status.Errorf(codes.Unknown, "panic triggered: %v", err)
		}))

	unaryClientInterceptors = append(unaryClientInterceptors, internalInterceptors...)

	options = append(options,
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			unaryClientInterceptors...,
		)),
	)

	return options, err
}
