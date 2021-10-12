package configs

type SignalStopper interface {
	Stop()
}

type SignalCloser interface {
	Close()
}

type SignalCloserWithErr interface {
	Close() error
}
