package server

import (
	pb "github.com/dysnix/ai-scale-proto/external/proto/health"
)

type Server interface {
	Start() <-chan error
	Stop() error
}

type Health interface {
	Server
	pb.HealthServer
}
