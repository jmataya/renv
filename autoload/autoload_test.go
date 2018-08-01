package autoload

import (
	"os"
	"testing"
)

func TestAutoload(t *testing.T) {
	want := "bar"
	got := os.Getenv("FOO")

	if want != got {
		t.Errorf("FOO=%s, want %s", got, want)
	}
}
