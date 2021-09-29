package server

type Server interface {
	Start() <-chan error
	Stop() error
}
