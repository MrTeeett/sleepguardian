//go:build linux

package netmon

import "testing"

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
