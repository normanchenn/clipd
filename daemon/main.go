package main

import (
	"os"
	"os/exec"
	"time"
)

var clipboard []string

const (
	FILEPATH    = "/Users/normanchen/Desktop/clipd.log"
	INTERVAL    = 10 * time.Millisecond
	THRESHOLD   = 500
	PERMISSIONS = 0777
)

func main() {
	poll()
}

func poll() {
	file, err := os.OpenFile(FILEPATH, os.O_RDWR|os.O_CREATE, PERMISSIONS)
	if err != nil {
		// log error
		return
	}
	defer file.Close()

	prev := ""
	for {
		cur, err := getClipboard()
		if err != nil {
			// log error
			continue
		}

		if cur != prev {
			clipboard = append(clipboard, cur)
			writeClipboard(file, cur)
			// if err != nil {
			// log error
			// }
			prev = cur
		}
		time.Sleep(INTERVAL)
	}
}

func getClipboard() (string, error) {
	cmd := exec.Command("pbpaste")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func writeClipboard(file *os.File, clipboard string) error {
	_, err := file.Seek(0, 2)
	if err != nil {
		return err
	}

	_, err = file.WriteString(clipboard + "\n")
	return err
}
