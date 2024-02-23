package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/normanchenn/clipd/daemon/clipboard"
	"github.com/normanchenn/clipd/daemon/history"
)

var clipboard_history = history.NewHistory()

func main() {
	filepath, interval, threshold, permissions, socketpath, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading .env variables")
		return
	}

	var mu sync.Mutex
	clipboard_channel := make(chan string)
	listener, err := net.Listen("unix", socketpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error listening")
		return
	}
	defer listener.Close()

	go poll(clipboard_channel, filepath, interval, threshold, permissions, &mu)
	go handleRequests(clipboard_channel, listener, &mu)

	select {}
}

func loadConfig() (string, time.Duration, int, os.FileMode, string, error) {
	// err := godotenv.Load(".env")
	// if err != nil {
	// fmt.Fprintln(os.Stderr, "Error loading .env file")
	// return "", 0, 0, 0, "", err
	// }
	// filepath := os.Getenv("FILEPATH")
	// interval_str := os.Getenv("INTERVAL")
	// threshold_str := os.Getenv("THRESHOLD")
	// permissions_str := os.Getenv("PERMISSIONS")
	// socketpath := os.Getenv("SOCKETPATH")
	var err error = nil
	filepath := "/Users/normanchen/Desktop/clipd.log"
	interval_str := "10"
	threshold_str := "500"
	permissions_str := "0777"
	socketpath := "/tmp/clipd.sock"
	if filepath == "" || interval_str == "" || threshold_str == "" || permissions_str == "" || socketpath == "" {
		fmt.Fprintln(os.Stderr, "Error loading .env variables")
		return "", 0, 0, 0, "", err
	}
	interval, _ := strconv.Atoi(interval_str)
	interval_time := time.Duration(interval) * time.Millisecond
	threshold, _ := strconv.Atoi(threshold_str)
	permissions, _ := strconv.ParseInt(permissions_str, 8, 32)
	file_permissions := os.FileMode(permissions)
	return filepath, interval_time, threshold, file_permissions, socketpath, nil
}

func poll(clipboard_channel <-chan string, filepath string, interval time.Duration, threshold int, permissions os.FileMode, mu *sync.Mutex) {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, permissions)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file")
		return
	}
	defer file.Close()

	prev := ""
	for {
		cur, err := clipboard.GetClipboard()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting clipboard")
			continue
		}

		if cur != prev {
			fmt.Println("New clipboard: ", cur)
			mu.Lock()
			clipboard_history.AddItem(cur, time.Now())
			mu.Unlock()
			err = clipboard.WriteClipboard(file, cur)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error writing to file")
				continue
			}
			prev = cur
		}
		time.Sleep(interval)
	}
}

func handleRequests(clipboard_channel <-chan string, listener net.Listener, mu *sync.Mutex) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error accepting connection")
			continue
		}
		go handleRequest(conn, mu)
	}
}

func handleRequest(conn net.Conn, mu *sync.Mutex) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	request, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading request: ", err)
		return
	}
	request = strings.TrimSpace(request)
	parts := strings.Fields(request)
	if len(parts) < 1 {
		fmt.Fprintln(os.Stderr, "Invalid request 1")
		return
	}

	switch parts[0] {
	case "get":
		handleGet(conn, mu, parts)
	default:
		fmt.Fprintf(conn, "Invalid request 2: %s", parts[0])
	}
}

func handleGet(conn net.Conn, mu *sync.Mutex, parts []string) {
	// format is this: ["get", "last=10", "from=10", "to=15", "at=10"]
	if len(parts) == 1 { // get most recent (no other args)
		mu.Lock()
		defer mu.Unlock()
		item := clipboard_history.GetItem(0)
		fmt.Fprintln(conn, item.GetContent())
		fmt.Fprintln(conn, item.GetTimestamp())
		fmt.Fprintln(conn, item)
	} else if len(parts) == 2 && strings.HasPrefix(parts[1], "at=") { // get n (at exists, last, from , to don't exist)
		at, err := strconv.Atoi(parts[1][3:])
		if err != nil {
			fmt.Fprintln(conn, "Invalid request 3")
			return
		}
		mu.Lock()
		defer mu.Unlock()
		item := clipboard_history.GetItem(at)
		fmt.Fprintln(conn, item.GetContent())
		fmt.Fprintln(conn, item.GetTimestamp())
		fmt.Fprintln(conn, item)
	} else if len(parts) == 2 && strings.HasPrefix(parts[1], "last=") { // get last n (last exists, at, from, to don't exist)
		last, err := strconv.Atoi(parts[1][5:])
		if err != nil {
			fmt.Fprintln(conn, "Invalid request 4")
			return
		}
		mu.Lock()
		defer mu.Unlock()
		items := clipboard_history.GetItemRange(0, last-1)
		for _, item := range items {
			fmt.Fprintln(conn, "----")
			fmt.Fprintln(conn, item.GetContent())
			fmt.Fprintln(conn, item.GetTimestamp())
			fmt.Fprintln(conn, item)
		}
	} else if len(parts) == 3 && strings.HasPrefix(parts[1], "from=") && strings.HasPrefix(parts[2], "to=") { // get from n to m (from, to exist, last, at don't exist)
		from, err := strconv.Atoi(parts[1][5:])
		if err != nil {
			fmt.Fprintln(conn, "Invalid request 5")
			return
		}
		to, err := strconv.Atoi(parts[2][3:])
		if err != nil {
			fmt.Fprintln(conn, "Invalid request 6")
			return
		}
		mu.Lock()
		defer mu.Unlock()
		items := clipboard_history.GetItemRange(from, to)
		for _, item := range items {
			fmt.Fprintln(conn, "----")
			fmt.Fprintln(conn, item.GetContent())
			fmt.Fprintln(conn, item.GetTimestamp())
			fmt.Fprintln(conn, item)
		}
	} else {
		fmt.Fprintln(conn, "Invalid request 7")
	}

}
