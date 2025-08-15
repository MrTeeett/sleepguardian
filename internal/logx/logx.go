package logx

import (
	"log/slog"
	"os"
	"strings"

	"github.com/MrTeeett/sleepguardian/internal/config"
)

func Setup(c config.Log) {
	var lvl slog.Level
	switch strings.ToLower(c.Level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	out := os.Stdout
	if c.File != "" {
		if f, err := os.OpenFile(c.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
			out = f
		}
	}
	var h slog.Handler
	if strings.ToLower(c.Format) == "json" {
		h = slog.NewJSONHandler(out, &slog.HandlerOptions{Level: lvl})
	} else {
		h = slog.NewTextHandler(out, &slog.HandlerOptions{Level: lvl})
	}
	slog.SetDefault(slog.New(h))
}
