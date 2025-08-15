//go:build linux

package diskmon

import (
	"bufio"
	"os"
	"strings"
	"time"
)

type linuxReader struct {
	prevR, prevW uint64
	prevStamp    time.Time
}

func New() Reader { return &linuxReader{} }

func (r *linuxReader) Snapshot(include, exclude []string) (Snapshot, error) {
	inc := set(include)
	exc := set(exclude)
	f, err := os.Open("/proc/diskstats")
	if err != nil {
		return Snapshot{}, err
	}
	defer f.Close()
	var readSectors, writeSectors uint64
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fs := strings.Fields(sc.Text())
		if len(fs) < 14 {
			continue
		}
		dev := fs[2]
		if exc[dev] {
			continue
		}
		if len(inc) > 0 && !inc[dev] {
			continue
		}
		// sectors read = fs[5], sectors written = fs[9]
		readSectors += parseU64(fs[5])
		writeSectors += parseU64(fs[9])
	}
	// 512 bytes per sector
	curR := readSectors * 512
	curW := writeSectors * 512
	now := time.Now()
	if r.prevStamp.IsZero() {
		r.prevR, r.prevW, r.prevStamp = curR, curW, now
		return Snapshot{Stamp: now, PerSec: true}, nil
	}
	dt := now.Sub(r.prevStamp).Seconds()
	rbps := perSec(curR, r.prevR, dt)
	wbps := perSec(curW, r.prevW, dt)
	r.prevR, r.prevW, r.prevStamp = curR, curW, now
	return Snapshot{Stamp: now, ReadBytes: rbps, WriteBytes: wbps, PerSec: true}, nil
}

func set(ss []string) map[string]bool {
	m := map[string]bool{}
	for _, s := range ss {
		m[s] = true
	}
	return m
}
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
func perSec(cur, prev uint64, dt float64) uint64 {
	if dt <= 0 || cur < prev {
		return 0
	}
	return uint64(float64(cur-prev) / dt)
}
