package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"github.com/normanchenn/clipd/daemon/clipboard"
	"github.com/normanchenn/clipd/daemon/history"
	"github.com/normanchenn/clipd/daemon/poll"
	"github.com/normanchenn/clipd/daemon/requests"
)

const (
	FILEPATH    = "/clipd/logs/clipd.log"
	INTERVAL    = "10"
	THRESHOLD   = "30"
	PERMISSIONS = "0777"
	SOCKETPATH  = "/tmp/clipd.sock"
)

var clipboard_history = history.InitHistory()

func main() {
	filepath, interval, threshold, permissions, socketpath, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading variables")
		return
	}

	if _, err := os.Stat(socketpath); err == nil {
		if err := os.Remove(socketpath); err != nil {
			fmt.Fprintln(os.Stderr, "Error removing socket file")
			return
		}
	}

	err = clipboard.WriteFileToClipboard(filepath, permissions, clipboard_history, threshold)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading existing file to clipboard")
	}

	signal_channel := make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt, syscall.SIGTERM)

	listener, err := net.Listen("unix", socketpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error listening", err)
		return
	}
	defer listener.Close()
	go poll.Poll(clipboard_history, filepath, interval, threshold, permissions)
	go requests.HandleRequests(listener, clipboard_history)

	<-signal_channel
	fmt.Fprintln(os.Stdout, "Exiting")
}

func loadConfig() (string, time.Duration, int, os.FileMode, string, error) {
	var err error = nil
	user, err := user.Current()
	if err != nil {
		fmt.Println("Error getting user: ", err)
		return "", 0, 0, 0, "", err
	}
	baseDir := user.HomeDir
	filepath := baseDir + FILEPATH
	interval_str := INTERVAL
	threshold_str := THRESHOLD
	permissions_str := PERMISSIONS
	socketpath := SOCKETPATH
	if filepath == "" || interval_str == "" || threshold_str == "" || permissions_str == "" || socketpath == "" {
		fmt.Fprintln(os.Stderr, "Error loading variables")
		return "", 0, 0, 0, "", err
	}
	interval, _ := strconv.Atoi(interval_str)
	interval_time := time.Duration(interval) * time.Millisecond
	threshold, _ := strconv.Atoi(threshold_str)
	permissions, _ := strconv.ParseInt(permissions_str, 8, 32)
	file_permissions := os.FileMode(permissions)
	return filepath, interval_time, threshold, file_permissions, socketpath, nil
}
