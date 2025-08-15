package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MrTeeett/sleepguardian/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write cfg: %v", err)
	}
	return path
}

func TestLoadSetsDefaults(t *testing.T) {
	cfgJSON := `{
        "log":{},
        "monitor":{"interval_ms":0,"weights":{"net":1,"disk":1},"active_threshold_mbps":1,"idle_threshold_pct":50},
        "exceptions":{},
        "sleep":{"mode":"none"},
        "tray":{"show":false}
    }`
	path := writeTempConfig(t, cfgJSON)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Monitor.IntervalMS != 1000 {
		t.Fatalf("expected default interval 1000, got %d", cfg.Monitor.IntervalMS)
	}
	if cfg.Monitor.ActiveThresholdBPS != 125000 {
		t.Fatalf("active threshold conversion incorrect: %d", cfg.Monitor.ActiveThresholdBPS)
	}
	if cfg.Monitor.IdleThresholdBPS != 62500 {
		t.Fatalf("idle threshold conversion incorrect: %d", cfg.Monitor.IdleThresholdBPS)
	}
}

func TestLoadThresholdError(t *testing.T) {
	cfgJSON := `{
        "log":{},
        "monitor":{"interval_ms":100,"weights":{"net":1,"disk":1},"active_threshold_mbps":0,"idle_threshold_pct":0},
        "exceptions":{},
        "sleep":{"mode":"none"},
        "tray":{"show":false}
    }`
	path := writeTempConfig(t, cfgJSON)
	if _, err := config.Load(path); err == nil {
		t.Fatalf("expected error when thresholds <=0")
	}
}

func TestMonitorInterval(t *testing.T) {
	cfg := config.Monitor{IntervalMS: 250}
	if d := cfg.Interval(); d.Milliseconds() != 250 {
		t.Fatalf("interval mismatch: %v", d)
	}
}
