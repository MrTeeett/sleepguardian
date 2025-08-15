package procmon

type ProcStat struct {
	Name       string
	Path       string
	PID        int
	ReadBytes  uint64
	WriteBytes uint64
	HasNetSock bool
}
type Reader interface {
	ByNamesOrPaths(names, paths []string) ([]ProcStat, error)
}
