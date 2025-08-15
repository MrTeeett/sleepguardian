package osutil

import "testing"

func TestOpenFileEmpty(t *testing.T) {
	if err := OpenFile(""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
