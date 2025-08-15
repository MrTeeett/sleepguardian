//go:build windows

package diskmon

import (
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	m := set([]string{"a", "b"})
	if !m["a"] || !m["b"] || len(m) != 2 {
		t.Fatalf("set %v", m)
	}
}

func TestParseU64(t *testing.T) {
	if v := parseU64("123"); v != 123 {
		t.Fatalf("%d", v)
	}
	if v := parseU64("12x3"); v != 12 {
		t.Fatalf("%d", v)
	}
	if v := parseU64("abc"); v != 0 {
		t.Fatalf("%d", v)
	}
}

func TestPerSec(t *testing.T) {
	if v := perSec(100, 50, 5); v != 10 {
		t.Fatalf("%d", v)
	}
	if v := perSec(10, 20, 5); v != 0 {
		t.Fatalf("%d", v)
	}
	if v := perSec(100, 50, 0); v != 0 {
		t.Fatalf("%d", v)
	}
}

func TestSnapshot(t *testing.T) {
	r := New().(*windowsReader)
	if _, err := r.Snapshot(nil, nil); err != nil {
		t.Fatalf("snapshot error: %v", err)
	}
	time.Sleep(10 * time.Millisecond)
	s, err := r.Snapshot(nil, nil)
	if err != nil {
		t.Fatalf("snapshot error: %v", err)
	}
	if !s.PerSec {
		t.Fatalf("expected persec")
	}
}
