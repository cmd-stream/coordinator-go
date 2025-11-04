package manager

import "time"

type CmdEnvelope struct {
	At time.Time
	Bs []byte
}
