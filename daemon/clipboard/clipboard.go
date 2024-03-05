package clipboard

import (
	"fmt"
	"github.com/normanchenn/clipd/daemon/history"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	IDENTIFIER = "CLIPD"
)

func GetClipboard() (string, error) {
	cmd := exec.Command("pbpaste")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func WriteClipboardToFile(filepath string, permissions os.FileMode, clipboard string) error {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, permissions)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file")
		return err
	}
	defer file.Close()

	_, err = file.Seek(0, 2)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error seeking EOF")
		return err
	}

	_, err = file.WriteString(IDENTIFIER + clipboard + "\n")
	return err
}

func WriteFileToClipboard(filepath string, permissions os.FileMode, clipboard *history.History, threshold int) error {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, permissions)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		fmt.Fprintln(os.Stderr, "Error opening file")
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	cur := 0
	lines := strings.Split(string(content), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if strings.Contains(line, IDENTIFIER) {
			cur = i
			threshold--
			if threshold == 0 {
				break
			}
		}
	}

	now := time.Now()
	item := ""
	for i := cur; i < len(lines); i++ {
		line := lines[i]
		if strings.Contains(line, IDENTIFIER) && item == "" {
			item += strings.TrimPrefix(line, IDENTIFIER)
		} else if strings.Contains(line, IDENTIFIER) && item != "" {
			clipboard.AddItem(item, now)
			item = ""
			item += strings.TrimPrefix(line, IDENTIFIER)
		} else {
			item += line
		}
	}
	if item != "" {
		clipboard.AddItem(item, time.Now())
	}
	clipboard.PrintHistory()
	return nil
}
