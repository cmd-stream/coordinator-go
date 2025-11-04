package coordinator

import (
	man "github.com/cmd-stream/coordinator-go/manager"
	"github.com/cmd-stream/core-go"
	csrv "github.com/cmd-stream/core-go/server"
	"github.com/cmd-stream/delegate-go"
	dsrv "github.com/cmd-stream/delegate-go/server"
	"github.com/cmd-stream/handler-go"
	"github.com/cmd-stream/transport-go"
)

// Options defines the configuration settings for initializing a server.
//
// These options are composed of modular components that configure different
// layers of the server, including transport, handler logic, delegate behavior,
// and base server setup.
type Options[T any] struct {
	Info      delegate.ServerInfo
	Core      []csrv.SetOption
	Delegate  []dsrv.SetOption
	Handler   []handler.SetOption
	Transport []transport.SetOption
	Manager   []man.SetOption
	// ExecutionPolicy ExecutionPolicyOptions

	MaxCmds              int
	ImmediateResultProxy ProxyFactory[T]
	DelayedCmdProxy      core.Proxy
}

type SetOption[T any] func(o *Options[T])

// WithServerInfo sets the ServerInfo for the server.
//
// ServerInfo helps the client identify a compatible server.
func WithServerInfo[T any](info delegate.ServerInfo) SetOption[T] {
	return func(o *Options[T]) { o.Info = info }
}

// WithCore applies core-level configuration options.
func WithCore[T any](ops ...csrv.SetOption) SetOption[T] {
	return func(o *Options[T]) { o.Core = ops }
}

// WithDelegate applies delegate-specific options.
//
// These options customize the behavior of the server delegate.
func WithDelegate[T any](ops ...dsrv.SetOption) SetOption[T] {
	return func(o *Options[T]) { o.Delegate = ops }
}

// WithHandler sets options for the connection handler.
//
// These options customize the behavior of the connection handler.
func WithHandler[T any](ops ...handler.SetOption) SetOption[T] {
	return func(o *Options[T]) { o.Handler = ops }
}

// WithTransport applies transport-specific options.
//
// These options configure the transport layer.
func WithTransport[T any](ops ...transport.SetOption) SetOption[T] {
	return func(o *Options[T]) { o.Transport = ops }
}

func WithManager[T any](ops ...man.SetOption) SetOption[T] {
	return func(o *Options[T]) { o.Manager = ops }
}

// // func WithRetryInterval(interval time.Duration) SetOption {
// // 	return func(o *Options[T]) { o.RetryInterval = interval }
// // }
//
// func WithFastRetryInterval[T any](interval time.Duration) SetOption[T] {
// 	return func(o *Options[T]) { o.FastRetryInterval = interval }
// }
//
// func WithSlowRetryInterval[T any](interval time.Duration) SetOption[T] {
// 	return func(o *Options[T]) { o.SlowRetryInterval = interval }
// }
//
// func WithJitterRange[T any](interval time.Duration) SetOption[T] {
// 	return func(o *Options[T]) { o.JitterRange = interval }
// }

// func WithStorage[T any](storage Storage) SetOption[T] {
// 	return func(o *Options[T]) { o.Storage = storage }
// }

func Apply[T any](ops []SetOption[T], o *Options[T]) {
	for i := range ops {
		if ops[i] != nil {
			ops[i](o)
		}
	}
}
