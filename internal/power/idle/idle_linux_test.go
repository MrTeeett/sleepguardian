//go:build linux

package idle

import "testing"

func TestUserIdleSeconds(t *testing.T) {
	r := New()
	if _, err := r.UserIdleSeconds(); err != nil && err.Error() != "xprintidle not available" {
		t.Fatalf("unexpected error: %v", err)
	}
}
