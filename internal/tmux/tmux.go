package tmux

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// InTmux returns true if the process is running inside tmux.
func InTmux() bool {
	return os.Getenv("TMUX") != ""
}

// CurrentPane returns the current tmux pane ID (e.g., "%2").
func CurrentPane() string {
	return os.Getenv("TMUX_PANE")
}

// TargetPane finds the "other" pane in the current window.
// Returns the pane ID of the first pane that isn't ours.
// Falls back to pane 0 of the current window.
func TargetPane() (string, error) {
	current := CurrentPane()
	if current == "" {
		return "", fmt.Errorf("not in tmux")
	}

	// List panes in the current window
	cmd := exec.Command("tmux", "list-panes", "-F", "#{pane_id}")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("list panes: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	for _, pane := range lines {
		pane = strings.TrimSpace(pane)
		if pane != "" && pane != current {
			return pane, nil
		}
	}

	return "", fmt.Errorf("no other pane found in current window")
}

// SendKeys sends keystrokes to a tmux pane.
// The target can be a pane ID ("%0"), or a session:window.pane spec.
func SendKeys(target string, text string) error {
	cmd := exec.Command("tmux", "send-keys", "-t", target, text, "Enter")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return fmt.Errorf("tmux send-keys: %s", errMsg)
	}
	return nil
}

// SendToOtherPane sends text + Enter to the other pane in the current window.
func SendToOtherPane(text string) error {
	target, err := TargetPane()
	if err != nil {
		return err
	}
	return SendKeys(target, text)
}
