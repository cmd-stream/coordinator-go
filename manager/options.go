package manager

import "time"

type Options struct {
	FastRetryInterval time.Duration
	SlowRetryInterval time.Duration
	JitterRange       time.Duration
	Storage           Storage
}

type SetOption func(o *Options)

func WithFastRetryInterval(interval time.Duration) SetOption {
	return func(o *Options) { o.FastRetryInterval = interval }
}

func WithSlowRetryInterval(interval time.Duration) SetOption {
	return func(o *Options) { o.SlowRetryInterval = interval }
}

func WithJitterRange(interval time.Duration) SetOption {
	return func(o *Options) { o.JitterRange = interval }
}

func WithStorage[T any](storage Storage) SetOption {
	return func(o *Options) { o.Storage = storage }
}

func Apply(ops []SetOption, o *Options) {
	for _, op := range ops {
		op(o)
	}
}
