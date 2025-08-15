//go:build linux

package netmon

import (
	"bufio"
	"os"
	"strings"
	"time"
)

type linuxReader struct {
	prev      Snapshot
	prevStamp time.Time
}

func New() Reader { return &linuxReader{} }

func (r *linuxReader) Snapshot(include, excludeIF []string) (Snapshot, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return Snapshot{}, err
	}
	defer file.Close()

	inc := make(map[string]bool)
	exc := make(map[string]bool)
	for _, s := range include {
		inc[s] = true
	}
	for _, s := range excludeIF {
		exc[s] = true
	}

	var rx, tx uint64
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		ifName := strings.TrimSpace(parts[0])
		if exc[ifName] {
			continue
		}
		if !sumIF(inc, ifName) {
			continue
		}
		fields := strings.Fields(strings.TrimSpace(parts[1]))
		if len(fields) < 16 {
			continue
		}
		// fields[0]=rx_bytes, fields[8]=tx_bytes
		rx += parseU64(fields[0])
		tx += parseU64(fields[8])
	}
	now := time.Now()
	cur := Snapshot{Stamp: now, BytesIn: rx, BytesOut: tx}
	if r.prevStamp.IsZero() {
		r.prev = cur
		r.prevStamp = now
		return Snapshot{Stamp: now, BytesIn: 0, BytesOut: 0, PerSec: true}, nil
	}
	dt := now.Sub(r.prevStamp).Seconds()
	inBps := perSec(cur.BytesIn, r.prev.BytesIn, dt)
	outBps := perSec(cur.BytesOut, r.prev.BytesOut, dt)
	r.prev = cur
	r.prevStamp = now
	return Snapshot{Stamp: now, BytesIn: inBps, BytesOut: outBps, PerSec: true}, nil
}

func parseU64(s string) uint64 {
	var v uint64
	for _, c := range s {
		if c < '0' || c > '9' {
			return v
		}
	}
	_, _ = os.Stat("") // noop to avoid inlining complaint
	var x uint64
	for i := 0; i < len(s); i++ {
		x = x*10 + uint64(s[i]-'0')
	}
	return x
}
func perSec(cur, prev uint64, dt float64) uint64 {
	if dt <= 0 {
		return 0
	}
	if cur < prev {
		return 0
	}
	return uint64(float64(cur-prev) / dt)
}
