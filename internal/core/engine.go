package core

import (
	"context"
	"log/slog"
	"math"
	"sync/atomic"
	"time"

	"github.com/MrTeeett/sleepguardian/internal/config"
	diskmon "github.com/MrTeeett/sleepguardian/internal/monitors/diskmon"
	netmon "github.com/MrTeeett/sleepguardian/internal/monitors/netmon"
	procmon "github.com/MrTeeett/sleepguardian/internal/monitors/procmon"
	"github.com/MrTeeett/sleepguardian/internal/power/idle"
	"github.com/MrTeeett/sleepguardian/internal/power/inhibit"
	"github.com/MrTeeett/sleepguardian/internal/power/trigger"
)

type Engine struct {
	cfg    config.Config
	ctx    context.Context
	cancel context.CancelFunc

	net  netmon.Reader
	disk diskmon.Reader
	proc procmon.Reader

	inhib inhibit.Inhibitor
	idle  idle.Reader
	trig  trigger.Trigger

	paused atomic.Bool
	active atomic.Bool
	since  atomic.Value // time.Time
}

func NewEngine(ctx context.Context, cfg config.Config) (*Engine, error) {
	cctx, cancel := context.WithCancel(ctx)
	e := &Engine{cfg: cfg, ctx: cctx, cancel: cancel}
	e.net = netmon.New()
	e.disk = diskmon.New()
	e.proc = procmon.New()
	e.inhib = inhibit.New()
	e.idle = idle.New()
	e.trig = trigger.New()
	e.since.Store(time.Now())
	return e, nil
}

func (e *Engine) Close() error { e.inhib.Release(); e.cancel(); return nil }
func (e *Engine) TogglePause() {
	p := !e.paused.Load()
	e.paused.Store(p)
	if p {
		_ = e.inhib.Release()
	}
}
func (e *Engine) TriggerSuspend()   { _ = e.trig.Suspend("guardian menu") }
func (e *Engine) TriggerHibernate() { _ = e.trig.Hibernate("guardian menu") }
func (e *Engine) Status() (bool, time.Time, bool) {
	return e.active.Load(), e.since.Load().(time.Time), e.paused.Load()
}

func (e *Engine) Run() error {
	tk := time.NewTicker(e.cfg.Monitor.Interval())
	defer tk.Stop()

	actWin := time.Duration(e.cfg.Monitor.MinActiveSec) * time.Second
	idlWin := time.Duration(e.cfg.Monitor.MinIdleSec) * time.Second

	var lastOn, lastOff time.Time

	for {
		select {
		case <-e.ctx.Done():
			return nil
		case <-tk.C:
			if e.paused.Load() {
				continue
			}

			ns, nerr := e.net.Snapshot(e.cfg.Monitor.InterfacesInclude, append(e.cfg.Monitor.InterfacesExclude, e.cfg.Exceptions.InterfacesExclude...))
			ds, derr := e.disk.Snapshot(e.cfg.Monitor.DisksInclude, e.cfg.Monitor.DisksExclude)
			if nerr != nil || derr != nil {
				slog.Debug("snapshot err", "net", nerr, "disk", derr)
			}

			netBps := sumBps(ns.BytesIn, ns.BytesOut, ns.PerSec)
			diskBps := sumBps(ds.ReadBytes, ds.WriteBytes, ds.PerSec)

			if len(e.cfg.Exceptions.ProcessNames)+len(e.cfg.Exceptions.ProcessPaths) > 0 {
				stats, _ := e.proc.ByNamesOrPaths(e.cfg.Exceptions.ProcessNames, e.cfg.Exceptions.ProcessPaths)
				var exclDisk uint64
				var exclHasNet bool
				for _, s := range stats {
					exclDisk += s.ReadBytes + s.WriteBytes
					if s.HasNetSock {
						exclHasNet = true
					}
				}
				if exclHasNet && likelyOnlyExcludedHaveNet(stats) {
					netBps = 0
				}
				diskBps = max0(diskBps - float64(exclDisk))
			}

			score := e.cfg.Monitor.Weights.Net*netBps + e.cfg.Monitor.Weights.Disk*diskBps
			now := time.Now()
			if !e.active.Load() {
				if score >= float64(e.cfg.Monitor.ActiveThresholdBPS) {
					if lastOn.IsZero() {
						lastOn = now
					}
					if now.Sub(lastOn) >= actWin {
						_ = e.inhib.Acquire("active transfer")
						e.active.Store(true)
						e.since.Store(now)
						slog.Info("inhibit ON", "score", int64(score), "netBps", int64(netBps), "diskBps", int64(diskBps))
						lastOff = time.Time{}
					}
				} else {
					lastOn = time.Time{}
				}
			} else {
				if score <= float64(e.cfg.Monitor.IdleThresholdBPS) {
					if lastOff.IsZero() {
						lastOff = now
					}
					if now.Sub(lastOff) >= idlWin {
						_ = e.inhib.Release()
						e.active.Store(false)
						e.since.Store(now)
						slog.Info("inhibit OFF", "score", int64(score))
						if e.cfg.Sleep.ImmediateOnEnd && e.shouldSleepNow() {
							go e.triggerSleep()
						}
						lastOn = time.Time{}
					}
				} else {
					lastOff = time.Time{}
				}
			}
		}
	}
}

func (e *Engine) shouldSleepNow() bool {
	sec, err := e.idle.UserIdleSeconds()
	if err != nil {
		return false
	}
	return sec >= uint64(e.cfg.Sleep.IdleGraceSec)
}
func (e *Engine) triggerSleep() {
	mode := e.cfg.Sleep.Mode
	if mode == "system" {
		mode = e.trig.SystemPreferred(e.cfg.Sleep.Fallback)
	}
	switch mode {
	case "suspend":
		_ = e.trig.Suspend("guardian finished")
	case "hibernate":
		_ = e.trig.Hibernate("guardian finished")
	}
}

func sumBps(a, b uint64, perSec bool) float64 {
	if !perSec {
		return float64(a + b)
	}
	return float64(a + b)
}
func likelyOnlyExcludedHaveNet(stats []procmon.ProcStat) bool {
	if len(stats) == 0 {
		return false
	}
	for _, s := range stats {
		if !s.HasNetSock {
			return false
		}
	}
	return true
}
func max0(f float64) float64 { return math.Max(0, f) }
