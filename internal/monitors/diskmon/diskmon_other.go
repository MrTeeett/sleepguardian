//go:build !linux && !windows

package diskmon

import "time"

type otherReader struct{}

func New() Reader { return &otherReader{} }
func (o *otherReader) Snapshot(_, _ []string) (Snapshot, error) {
	return Snapshot{Stamp: time.Now(), PerSec: true}, nil
}
