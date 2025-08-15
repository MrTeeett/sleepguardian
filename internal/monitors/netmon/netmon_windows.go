//go:build windows

package netmon

import (
	"github.com/shirou/gopsutil/v3/net"
	"time"
)

type windowsReader struct {
	prev      Snapshot
	prevStamp time.Time
}

func New() Reader { return &windowsReader{} }

func (r *windowsReader) Snapshot(include, excludeIF []string) (Snapshot, error) {
	inc := make(map[string]bool)
	exc := make(map[string]bool)
	for _, s := range include {
		inc[s] = true
	}
	for _, s := range excludeIF {
		exc[s] = true
	}
	counters, err := net.IOCounters(true)
	if err != nil {
		return Snapshot{}, err
	}
	var rx, tx uint64
	for _, c := range counters {
		name := c.Name
		if exc[name] {
			continue
		}
		if !sumIF(inc, name) {
			continue
		}
		rx += c.BytesRecv
		tx += c.BytesSent
	}
	now := time.Now()
	cur := Snapshot{Stamp: now, BytesIn: rx, BytesOut: tx}
	if r.prevStamp.IsZero() {
		r.prev = cur
		r.prevStamp = now
		return Snapshot{Stamp: now, PerSec: true}, nil
	}
	dt := now.Sub(r.prevStamp).Seconds()
	inBps := perSec(cur.BytesIn, r.prev.BytesIn, dt)
	outBps := perSec(cur.BytesOut, r.prev.BytesOut, dt)
	r.prev = cur
	r.prevStamp = now
	return Snapshot{Stamp: now, BytesIn: inBps, BytesOut: outBps, PerSec: true}, nil
}

func perSec(cur, prev uint64, dt float64) uint64 {
	if dt <= 0 || cur < prev {
		return 0
	}
	return uint64(float64(cur-prev) / dt)
}

func parseU64(s string) uint64 {
	var v uint64
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		v = v*10 + uint64(c-'0')
	}
	return v
}
