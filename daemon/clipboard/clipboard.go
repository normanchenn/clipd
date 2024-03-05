package clipboard

import (
	"bufio"
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

	_, err = file.WriteString(IDENTIFIER + " " + clipboard + "\n")
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

	file_info, err := file.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error seeking file")
		return err
	}

	file_size := file_info.Size()
	if file_size == 0 {
		return nil
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error seeking end")
		return err
	}

	cur := int64(0)
	offset := int64(0)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, IDENTIFIER) {
			cur = int64(file_size - offset)
			threshold--
			if threshold == 0 {
				break
			}
		}
		offset += int64(len(line) + 1)
		_, err := file.Seek(-offset, io.SeekEnd)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error seeking next line")
			return err
		}
	}

	// start from cur, and then add to the clipboard
	_, err = file.Seek(cur, io.SeekStart)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error seeking start")
		return err
	}

	now := time.Now()
	item := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, IDENTIFIER) && item == "" {
			item += strings.TrimPrefix(line, IDENTIFIER)
		} else if strings.Contains(line, IDENTIFIER) && item != "" {
			clipboard.AddItem(item, now)
			item = ""
		} else {
			item += line
		}
	}
	if item != "" {
		clipboard.AddItem(item, time.Now())
	}
	return nil
}
