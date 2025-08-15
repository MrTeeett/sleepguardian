//go:build linux

package procmon

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type linuxReader struct{}

func New() Reader { return &linuxReader{} }

func (l *linuxReader) ByNamesOrPaths(names, paths []string) ([]ProcStat, error) {
	nameSet := set(names)
	pathSet := set(paths)
	var out []ProcStat
	d, err := os.ReadDir("/proc")
	if err != nil {
		return out, err
	}
	for _, e := range d {
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}
		base := "/proc/" + e.Name()
		comm, _ := os.ReadFile(base + "/comm")
		name := strings.TrimSpace(string(comm))
		if !nameSet[name] && len(pathSet) > 0 {
			// можно сравнить exe symlink
			exe, _ := os.Readlink(base + "/exe")
			if !pathSet[exe] {
				continue
			}
		} else if len(nameSet) > 0 && !nameSet[name] {
			continue
		}
		var r, w uint64
		if f, err := os.Open(base + "/io"); err == nil {
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				line := sc.Text()
				if strings.HasPrefix(line, "read_bytes:") {
					r += parseU64(strings.TrimSpace(strings.TrimPrefix(line, "read_bytes:")))
				}
				if strings.HasPrefix(line, "write_bytes:") {
					w += parseU64(strings.TrimSpace(strings.TrimPrefix(line, "write_bytes:")))
				}
			}
			f.Close()
		}
		hasNet := fileNonEmpty(base+"/net/tcp") || fileNonEmpty(base+"/net/tcp6") || fileNonEmpty(base+"/net/udp")
		out = append(out, ProcStat{
			Name:      name,
			Path:      mustReadlink(base + "/exe"),
			PID:       pid,
			ReadBytes: r, WriteBytes: w,
			HasNetSock: hasNet,
		})
	}
	return out, nil
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
func fileNonEmpty(p string) bool   { fi, err := os.Stat(p); return err == nil && fi.Size() > 0 }
func mustReadlink(p string) string { s, _ := os.Readlink(p); return filepath.Clean(s) }
