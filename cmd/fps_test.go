package cmd

import (
	"strings"
	"testing"
)

func TestInvalidFPSRejectedBeforeInput(t *testing.T) {
	old := fps
	defer func() { fps = old }()
	for _, value := range []int{0, -1} {
		fps = value
		err := run(rootCmd, []string{"synthetic-nonexistent-input"})
		if err == nil || !strings.Contains(err.Error(), "fps must be positive") {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}
