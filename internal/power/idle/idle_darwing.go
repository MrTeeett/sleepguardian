//go:build darwin

package idle

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

type impl struct{}

func newImpl() Reader { return &impl{} }

func (i *impl) UserIdleSeconds() (uint64, error) {
	out, err := exec.Command("ioreg", "-c", "IOHIDSystem").Output()
	if err != nil {
		return 0, err
	}
	// ищем "HIDIdleTime = NNNNNN" (наносекунды)
	lines := bytes.Split(out, []byte("\n"))
	for _, l := range lines {
		s := string(l)
		if strings.Contains(s, "HIDIdleTime") {
			fs := strings.Fields(s)
			n, _ := strconv.ParseUint(fs[len(fs)-1], 10, 64)
			return n / 1_000_000_000, nil
		}
	}
	return 0, nil
}
