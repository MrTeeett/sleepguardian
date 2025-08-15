package diskmon

import "time"

type Snapshot struct {
	Stamp      time.Time
	ReadBytes  uint64
	WriteBytes uint64
	PerSec     bool
}
type Reader interface {
	Snapshot(include, exclude []string) (Snapshot, error)
}
