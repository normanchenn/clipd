package clipboard

import (
	"os"
	"os/exec"
)

func GetClipboard() (string, error) {
	cmd := exec.Command("pbpaste")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func WriteClipboard(file *os.File, clipboard string) error {
	_, err := file.Seek(0, 2)
	if err != nil {
		return err
	}

	_, err = file.WriteString(clipboard + "\n")
	return err
}