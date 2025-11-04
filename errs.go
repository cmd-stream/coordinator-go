package coordinator

import "errors"

var (
	ErrCmdDelayed = errors.New("cmd delayed")
	ErrCmdBlocked = errors.New("cmd blocked")
)
