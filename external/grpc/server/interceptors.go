package server

import (
	"context"
	"errors"

	"google.golang.org/grpc"
)

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
