//go:build linux

package trigger

import "testing"

func TestSystemPreferred(t *testing.T) {
	tr := newImpl()
	if tr.SystemPreferred("") != "suspend" {
		t.Fatalf("default not suspend")
	}
	if tr.SystemPreferred("hibernate") != "hibernate" {
		t.Fatalf("fallback not respected")
	}
}
