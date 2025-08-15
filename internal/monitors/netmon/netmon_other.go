//go:build !linux

package netmon

import "time"

type otherReader struct{}

func New() Reader { return &otherReader{} }
func (o *otherReader) Snapshot(_, _ []string) (Snapshot, error) {
	return Snapshot{Stamp: time.Now(), BytesIn: 0, BytesOut: 0, PerSec: true}, nil
}
