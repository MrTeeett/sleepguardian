//go:build linux

package inhibit

import "testing"

func TestReleaseWithoutAcquire(t *testing.T) {
	i := New()
	if err := i.Release(); err != nil {
		t.Fatalf("release: %v", err)
	}
}
