package configs

type SignalStopper interface {
	Stop()
}

type SignalStopperWithErr interface {
	Stop() error
}

type SignalCloser interface {
	Close()
}

type SignalCloserWithErr interface {
	Close() error
}

type SingleUseGetter interface {
	SingleEnabled() bool
}
