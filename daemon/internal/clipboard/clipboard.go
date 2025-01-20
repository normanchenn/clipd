package clipboard

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Get() (string, error) {
	cmd := exec.Command("pbpaste")

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get clipboard content: %w", err)
	}
	return out.String(), nil
}
