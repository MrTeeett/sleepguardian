package logx

import (
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/MrTeeett/sleepguardian/internal/config"
)

func TestSetupWritesLog(t *testing.T) {
	tmp := t.TempDir()
	f := tmp + "/log.txt"
	c := Setup(config.Log{Level: "info", Format: "json", File: f})
	slog.Info("hello", "k", 1)
	if c != nil {
		_ = c.Close()
	}
	b, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(b), "\"msg\":\"hello\"") {
		t.Fatalf("log not written: %s", string(b))
	}
}
