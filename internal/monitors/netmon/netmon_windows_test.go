//go:build windows

package netmon

import (
	"testing"
	"time"
)

func TestSumIF(t *testing.T) {
	m := map[string]bool{"a": true}
	if !sumIF(m, "a") || sumIF(m, "b") {
		t.Fatalf("sumIF failed")
	}
	if sumIF(nil, "c") == false {
		t.Fatalf("nil map true")
	}
}

func TestParseU64(t *testing.T) {
	if v := parseU64("123"); v != 123 {
		t.Fatalf("%d", v)
	}
	if v := parseU64("12a"); v != 0 {
		t.Fatalf("expected 0, got %d", v)
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
