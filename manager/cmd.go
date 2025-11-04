package manager

import "github.com/cmd-stream/core-go"

type Cmd[T any] interface {
	core.Cmd[T]
	SetWorkflow(workflow *Workflow[T]) Cmd[T]
	Workflow() *Workflow[T]
}
