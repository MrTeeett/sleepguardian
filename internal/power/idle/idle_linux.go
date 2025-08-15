//go:build linux

package idle

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

type impl struct{}

func newImpl() Reader { return &impl{} }

func (i *impl) UserIdleSeconds() (uint64, error) {
	out, err := exec.Command("xprintidle").Output()
	if err != nil {
		return 0, errors.New("xprintidle not available")
	}
	msStr := strings.TrimSpace(string(out))
	ms, _ := strconv.ParseUint(msStr, 10, 64)
	return ms / 1000, nil
}
