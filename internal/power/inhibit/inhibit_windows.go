//go:build windows

package inhibit

import "syscall"

var (
	modkernel32                         = syscall.NewLazyDLL("kernel32.dll")
	procSetThreadExecutionState         = modkernel32.NewProc("SetThreadExecutionState")
	ES_CONTINUOUS               uintptr = 0x80000000
	ES_SYSTEM_REQUIRED          uintptr = 0x00000001
	ES_AWAYMODE_REQUIRED        uintptr = 0x00000040
)

type impl struct{ active bool }

func newImpl() Inhibitor { return &impl{} }

func (i *impl) Acquire(_ string) error {
	if i.active {
		return nil
	}
	procSetThreadExecutionState.Call(ES_CONTINUOUS | ES_SYSTEM_REQUIRED | ES_AWAYMODE_REQUIRED)
	i.active = true
	return nil
}
func (i *impl) Release() error {
	if !i.active {
		return nil
	}
	procSetThreadExecutionState.Call(ES_CONTINUOUS)
	i.active = false
	return nil
}
