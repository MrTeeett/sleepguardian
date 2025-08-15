package core

import (
	"errors"
	"testing"
	"time"

	"github.com/MrTeeett/sleepguardian/internal/config"
	procmon "github.com/MrTeeett/sleepguardian/internal/monitors/procmon"
)

type idleStub struct {
	secs uint64
	err  error
}

func (i idleStub) UserIdleSeconds() (uint64, error) { return i.secs, i.err }

type trigStub struct {
	susp, hiber int
	sys         string
}

func (t *trigStub) Suspend(string) error          { t.susp++; return nil }
func (t *trigStub) Hibernate(string) error        { t.hiber++; return nil }
func (t *trigStub) SystemPreferred(string) string { return t.sys }

type inhibStub struct{ released bool }

func (i *inhibStub) Acquire(string) error { return nil }
func (i *inhibStub) Release() error       { i.released = true; return nil }

func TestSumBps(t *testing.T) {
	if sum := sumBps(1, 2, false); sum != 3 {
		t.Fatalf("sum %v", sum)
	}
	if sum := sumBps(1, 2, true); sum != 3 {
		t.Fatalf("sum %v", sum)
	}
}

func TestLikelyOnlyExcludedHaveNet(t *testing.T) {
	stats := []procmon.ProcStat{{HasNetSock: true}, {HasNetSock: true}}
	if !likelyOnlyExcludedHaveNet(stats) {
		t.Fatalf("expected true")
	}
	stats = append(stats, procmon.ProcStat{HasNetSock: false})
	if likelyOnlyExcludedHaveNet(stats) {
		t.Fatalf("expected false")
	}
}

func TestMax0(t *testing.T) {
	if max0(-1) != 0 {
		t.Fatalf("neg")
	}
	if max0(5) != 5 {
		t.Fatalf("pos")
	}
}

func TestShouldSleepNow(t *testing.T) {
	e := &Engine{cfg: config.Config{Sleep: config.Sleep{IdleGraceSec: 10}}, idle: idleStub{secs: 15}}
	if !e.shouldSleepNow() {
		t.Fatalf("expected true")
	}
	e.idle = idleStub{secs: 5}
	if e.shouldSleepNow() {
		t.Fatalf("expected false")
	}
	e.idle = idleStub{err: errors.New("x")}
	if e.shouldSleepNow() {
		t.Fatalf("expected false on err")
	}
}

func TestTriggerSleep(t *testing.T) {
	tr := &trigStub{}
	e := &Engine{cfg: config.Config{Sleep: config.Sleep{Mode: "suspend"}}, trig: tr}
	e.triggerSleep()
	if tr.susp != 1 || tr.hiber != 0 {
		t.Fatalf("unexpected calls %d %d", tr.susp, tr.hiber)
	}

	tr = &trigStub{sys: "hibernate"}
	e = &Engine{cfg: config.Config{Sleep: config.Sleep{Mode: "system", Fallback: "hibernate"}}, trig: tr}
	e.triggerSleep()
	if tr.hiber != 1 {
		t.Fatalf("hibernate not called")
	}
}

func TestTogglePause(t *testing.T) {
	inh := &inhibStub{}
	e := &Engine{inhib: inh}
	e.TogglePause()
	if !e.paused.Load() {
		t.Fatalf("expected paused")
	}
	if !inh.released {
		t.Fatalf("release not called")
	}
}

func TestStatus(t *testing.T) {
	e := &Engine{}
	e.since.Store(time.Unix(0, 0))
	a, since, p := e.Status()
	if a || p || since.Unix() != 0 {
		t.Fatalf("unexpected status")
	}
}
