package tmux

import (
	"os"
	"testing"
)

func TestInTmux(t *testing.T) {
	// Save and restore
	orig := os.Getenv("TMUX")
	defer os.Setenv("TMUX", orig)

	os.Setenv("TMUX", "")
	if InTmux() {
		t.Error("expected false when TMUX is empty")
	}

	os.Setenv("TMUX", "/tmp/tmux-501/default,12345,0")
	if !InTmux() {
		t.Error("expected true when TMUX is set")
	}
}

func TestCurrentPane(t *testing.T) {
	orig := os.Getenv("TMUX_PANE")
	defer os.Setenv("TMUX_PANE", orig)

	os.Setenv("TMUX_PANE", "%3")
	if got := CurrentPane(); got != "%3" {
		t.Errorf("CurrentPane() = %q, want %%3", got)
	}

	os.Setenv("TMUX_PANE", "")
	if got := CurrentPane(); got != "" {
		t.Errorf("CurrentPane() = %q, want empty", got)
	}
}

func TestTargetPane_NotInTmux(t *testing.T) {
	orig := os.Getenv("TMUX_PANE")
	defer os.Setenv("TMUX_PANE", orig)

	os.Setenv("TMUX_PANE", "")
	_, err := TargetPane()
	if err == nil {
		t.Error("expected error when not in tmux")
	}
}
