//go:build windows

package trigger

import (
	"os/exec"
)

type impl struct{}

func newImpl() Trigger { return &impl{} }

func (i *impl) Suspend(_ string) error {
	return exec.Command("rundll32.exe", "powrprof.dll,SetSuspendState", "0,1,0").Run()
}
func (i *impl) Hibernate(_ string) error { return exec.Command("shutdown", "/h").Run() }
func (i *impl) SystemPreferred(fb string) string {
	if fb == "" {
		return "suspend"
	}
	return fb
}
