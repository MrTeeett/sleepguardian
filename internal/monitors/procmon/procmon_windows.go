//go:build windows

package procmon

import (
	"os"
	"path/filepath"

	"github.com/shirou/gopsutil/v3/process"
)

type windowsReader struct{}

func New() Reader { return &windowsReader{} }

func (w *windowsReader) ByNamesOrPaths(names, paths []string) ([]ProcStat, error) {
	nameSet := set(names)
	pathSet := set(paths)
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}
	var out []ProcStat
	for _, p := range procs {
		name, _ := p.Name()
		exe, _ := p.Exe()
		if len(nameSet) > 0 && !nameSet[name] {
			if len(pathSet) == 0 || !pathSet[exe] {
				continue
			}
		}
		if len(pathSet) > 0 && !pathSet[exe] {
			if len(nameSet) == 0 || !nameSet[name] {
				continue
			}
		}
		io, _ := p.IOCounters()
		var r, wbytes uint64
		if io != nil {
			r = io.ReadBytes
			wbytes = io.WriteBytes
		}
		conns, _ := p.Connections()
		hasNet := len(conns) > 0
		out = append(out, ProcStat{
			Name:       name,
			Path:       exe,
			PID:        int(p.Pid),
			ReadBytes:  r,
			WriteBytes: wbytes,
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
