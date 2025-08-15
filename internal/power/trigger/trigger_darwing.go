//go:build darwin

package trigger

import "os/exec"

type impl struct{}

func newImpl() Trigger { return &impl{} }

func (i *impl) Suspend(_ string) error { return exec.Command("pmset", "sleepnow").Run() }
func (i *impl) Hibernate(_ string) error {
	return exec.Command("pmset", "hibernatemode", "25", ";", "pmset", "sleepnow").Run()
}
func (i *impl) SystemPreferred(fb string) string {
	if fb == "" {
		return "suspend"
	}
	return fb
}
