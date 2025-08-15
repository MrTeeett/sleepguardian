package osutil

import (
	"os/exec"
	"runtime"
)

func OpenFile(path string) error {
	if path == "" {
		return nil
	}
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd", "/c", "start", "", path).Start()
	case "darwin":
		return exec.Command("open", path).Start()
	default:
		return exec.Command("xdg-open", path).Start()
	}
}
