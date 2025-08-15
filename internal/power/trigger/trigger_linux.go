//go:build linux

package trigger

import "os/exec"

type impl struct{}

func newImpl() Trigger { return &impl{} }

func (i *impl) Suspend(_ string) error   { return exec.Command("systemctl", "suspend").Run() }
func (i *impl) Hibernate(_ string) error { return exec.Command("systemctl", "hibernate").Run() }
func (i *impl) SystemPreferred(fb string) string {
	if fb == "" {
		return "suspend"
	}
	return fb
}
