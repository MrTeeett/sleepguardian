//go:build windows

package idle

import (
	"syscall"
	"time"
	"unsafe"
)

type impl struct{}

func newImpl() Reader { return &impl{} }

type lastinputinfo struct {
	cbSize uint32
	dwTime uint32
}

func (i *impl) UserIdleSeconds() (uint64, error) {
	user32 := syscall.NewLazyDLL("user32.dll")
	proc := user32.NewProc("GetLastInputInfo")
	var lii lastinputinfo
	lii.cbSize = uint32(unsafe.Sizeof(lii))
	proc.Call(uintptr(unsafe.Pointer(&lii)))
	tick := time.Now().UnixMilli()
	idleMs := uint64(uint32(tick) - lii.dwTime)
	return idleMs / 1000, nil
}
