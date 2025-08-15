//go:build windows

package diskmon

import (
	"github.com/shirou/gopsutil/v3/disk"
	"time"
)

type windowsReader struct {
	prevR, prevW uint64
	prevStamp    time.Time
}

func New() Reader { return &windowsReader{} }

func (r *windowsReader) Snapshot(include, exclude []string) (Snapshot, error) {
	inc := set(include)
	exc := set(exclude)
	counters, err := disk.IOCounters()
	if err != nil {
		return Snapshot{}, err
	}
	var readBytes, writeBytes uint64
	for name, c := range counters {
		if exc[name] {
			continue
		}
		if len(inc) > 0 && !inc[name] {
			continue
		}
		readBytes += c.ReadBytes
		writeBytes += c.WriteBytes
	}
	now := time.Now()
	if r.prevStamp.IsZero() {
		r.prevR, r.prevW, r.prevStamp = readBytes, writeBytes, now
		return Snapshot{Stamp: now, PerSec: true}, nil
	}
	dt := now.Sub(r.prevStamp).Seconds()
	rbps := perSec(readBytes, r.prevR, dt)
	wbps := perSec(writeBytes, r.prevW, dt)
	r.prevR, r.prevW, r.prevStamp = readBytes, writeBytes, now
	return Snapshot{Stamp: now, ReadBytes: rbps, WriteBytes: wbps, PerSec: true}, nil
}

func set(ss []string) map[string]bool {
	m := map[string]bool{}
	for _, s := range ss {
		m[s] = true
	}
	return m
}

func perSec(cur, prev uint64, dt float64) uint64 {
	if dt <= 0 || cur < prev {
		return 0
	}
	return uint64(float64(cur-prev) / dt)
}

// parseU64 is kept for test parity with linux implementation
func parseU64(s string) uint64 {
	var x uint64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		x = x*10 + uint64(c-'0')
	}
	return x
}
