//go:build !linux && !windows

package procmon

type otherReader struct{}

func New() Reader                                                       { return &otherReader{} }
func (o *otherReader) ByNamesOrPaths(_, _ []string) ([]ProcStat, error) { return nil, nil }
