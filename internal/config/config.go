package config

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

type Log struct{ Level, Format, File string }
type Weights struct{ Net, Disk float64 }

type Monitor struct {
	IntervalMS         int      `json:"interval_ms"`
	Weights            Weights  `json:"weights"`
	ActiveThresholdBPS int64    `json:"active_threshold_bps"`
	IdleThresholdBPS   int64    `json:"idle_threshold_bps"`
	MinActiveSec       int      `json:"min_active_sec"`
	MinIdleSec         int      `json:"min_idle_sec"`
	InterfacesInclude  []string `json:"interfaces_include"`
	InterfacesExclude  []string `json:"interfaces_exclude"`
	DisksInclude       []string `json:"disks_include"`
	DisksExclude       []string `json:"disks_exclude"`
}

type Exceptions struct {
	ProcessNames      []string `json:"process_names"`
	ProcessPaths      []string `json:"process_paths"`
	PortsExclude      []int    `json:"ports_exclude"`
	IPsExclude        []string `json:"ips_exclude"`
	InterfacesExclude []string `json:"interfaces_exclude"`
}

type Sleep struct {
	Mode           string `json:"mode"`     // system|suspend|hibernate|none
	Fallback       string `json:"fallback"` // suspend|hibernate|none
	IdleGraceSec   int    `json:"idle_grace_sec"`
	ImmediateOnEnd bool   `json:"immediate_on_finish"`
}

type Tray struct {
	Show bool `json:"show"`
}

type Config struct {
	Log        Log        `json:"log"`
	Monitor    Monitor    `json:"monitor"`
	Exceptions Exceptions `json:"exceptions"`
	Sleep      Sleep      `json:"sleep"`
	Tray       Tray       `json:"tray"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	if cfg.Monitor.IntervalMS <= 0 {
		cfg.Monitor.IntervalMS = 1000
	}
	if cfg.Monitor.ActiveThresholdBPS <= 0 || cfg.Monitor.IdleThresholdBPS <= 0 {
		return nil, errors.New("thresholds must be > 0")
	}
	return &cfg, nil
}
func (m Monitor) Interval() time.Duration { return time.Duration(m.IntervalMS) * time.Millisecond }
