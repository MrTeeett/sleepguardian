package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/MrTeeett/sleepguardian/internal/config"
	"github.com/MrTeeett/sleepguardian/internal/core"
	"github.com/MrTeeett/sleepguardian/internal/logx"
	"github.com/MrTeeett/sleepguardian/internal/osutil"
	"github.com/MrTeeett/sleepguardian/internal/tray"
)

func main() {
	var cfgPathFlag string
	flag.StringVar(&cfgPathFlag, "config", "", "Путь к config.json (если пусто — ищем рядом с программой)")
	flag.Parse()

	cfgPath := cfgPathFlag
	if cfgPath == "" {
		cfgPath = configPathNearExe("config.json")
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		slog.Error("config load failed", "err", err, "path", cfgPath)
		os.Exit(1)
	}

	logx.Setup(cfg.Log)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	engine, err := core.NewEngine(ctx, *cfg)
	if err != nil {
		slog.Error("engine init failed", "err", err)
		os.Exit(1)
	}
	defer engine.Close()

	if cfg.Tray.Show {
		go tray.Run(tray.Options{
			StatusFn:       engine.Status,
			OnPauseToggle:  engine.TogglePause,
			OnSleepNow:     engine.TriggerSuspend,
			OnHibernateNow: engine.TriggerHibernate,
			OnOpenLog:      func() { _ = osutil.OpenFile(cfg.Log.File) },
			OnOpenConfig:   func() { _ = osutil.OpenFile(cfgPath) },
			OnExit:         func() { cancel() },
		})
	}

	if err := engine.Run(); err != nil {
		slog.Error("engine stopped", "err", err)
	}
}

func configPathNearExe(name string) string {
	exe, err := os.Executable()
	if err != nil {
		return filepath.Join(".", name)
	}
	exe, _ = filepath.EvalSymlinks(exe)
	return filepath.Join(filepath.Dir(exe), name)
}

func init() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	_ = time.Now
}
