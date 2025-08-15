//go:build windows

package procmon

import (
	"os"
	"path/filepath"
	"testing"
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
	if v := parseU64("12x"); v != 12 {
		t.Fatalf("%d", v)
	}
	if v := parseU64("abc"); v != 0 {
		t.Fatalf("%d", v)
	}
}

func TestFileNonEmpty(t *testing.T) {
	tmp := t.TempDir()
	empty := filepath.Join(tmp, "e.txt")
	os.WriteFile(empty, nil, 0644)
	filled := filepath.Join(tmp, "f.txt")
	os.WriteFile(filled, []byte("x"), 0644)
	if fileNonEmpty(empty) {
		t.Fatalf("empty detected non-empty")
	}
	if !fileNonEmpty(filled) {
		t.Fatalf("filled detected empty")
	}
}

func TestMustReadlink(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "t.txt")
	os.WriteFile(target, []byte("a"), 0644)
	link := filepath.Join(tmp, "l.txt")
	os.Symlink(target, link)
	got := mustReadlink(link)
	if filepath.Clean(target) != got {
		t.Fatalf("readlink %s != %s", got, target)
	}
}
