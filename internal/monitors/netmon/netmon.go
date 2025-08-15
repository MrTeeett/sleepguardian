package netmon

import "time"

type Snapshot struct {
	Stamp    time.Time
	BytesIn  uint64
	BytesOut uint64
	PerSec   bool
}
type Reader interface {
	Snapshot(include, excludeIF []string) (Snapshot, error)
}

func sumIF(ifn map[string]bool, name string) bool {
	if len(ifn) == 0 {
		return true
	}
	return ifn[name]
}
