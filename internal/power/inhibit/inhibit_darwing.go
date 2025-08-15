//go:build darwin

package inhibit

import "os/exec"

type impl struct{ cmd *exec.Cmd }

func newImpl() Inhibitor { return &impl{} }

func (i *impl) Acquire(reason string) error {
	if i.cmd != nil {
		return nil
	}
	i.cmd = exec.Command("caffeinate", "-dimsu", "-w", "0") // бесконечно
	return i.cmd.Start()
}
func (i *impl) Release() error {
	if i.cmd == nil {
		return nil
	}
	err := i.cmd.Process.Kill()
	i.cmd = nil
	return err
}
