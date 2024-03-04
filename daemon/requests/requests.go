package requests

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/normanchenn/clipd/daemon/history"
)

func HandleRequests(listener net.Listener, clipboard_history *history.History) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error accepting connection")
			continue
		}
		fmt.Fprintln(os.Stdout, "New connection")
		go handleRequest(conn, clipboard_history)
	}
}

func handleRequest(conn net.Conn, clipboard_history *history.History) {
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
		handleGet(conn, parts, clipboard_history)
	default:
		fmt.Fprintf(conn, "Invalid request 2: %s", parts[0])
	}
}

func handleGet(conn net.Conn, parts []string, clipboard_history *history.History) { // format is this: ["get", "last=10", "from=10", "to=15", "at=10"]
	if len(parts) == 1 { // get most recent (no other args)
		clipboard_history.Lock()
		defer clipboard_history.Unlock()
		item := clipboard_history.GetItem(0)
		printItems(conn, []*history.HistoryItem{item})
	} else if len(parts) == 2 && strings.HasPrefix(parts[1], "at=") { // get n (at exists, last, from , to don't exist)
		at, err := strconv.Atoi(parts[1][3:])
		if err != nil {
			fmt.Fprintln(conn, "Invalid request 3")
			return
		}
		fmt.Fprintln(os.Stdout, "parts: ", parts)
		fmt.Fprintln(os.Stdout, "at: ", at)

		clipboard_history.Lock()
		defer clipboard_history.Unlock()
		item := clipboard_history.GetItem(at)
		printItems(conn, []*history.HistoryItem{item})
	} else if len(parts) == 2 && strings.HasPrefix(parts[1], "last=") { // get last n (last exists, at, from, to don't exist)
		last, err := strconv.Atoi(parts[1][5:])
		if err != nil {
			fmt.Fprintln(conn, "Invalid request 4")
			return
		}

		clipboard_history.Lock()
		defer clipboard_history.Unlock()
		items := clipboard_history.GetItemRange(0, last-1)
		printItems(conn, items)
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

		clipboard_history.Lock()
		defer clipboard_history.Unlock()
		items := clipboard_history.GetItemRange(from, to)
		printItems(conn, items)
	} else {
		fmt.Fprintln(conn, "Invalid request 7")
	}
}

func printItems(conn net.Conn, items []*history.HistoryItem) {
	var ret string
	for _, item := range items {
		// fmt.Fprintln(conn, "--------------------------------")
		fmt.Fprintln(os.Stdout, "returning item: ", item.GetContent())
		// fmt.Fprintln(conn, item.GetContent())
		ret += item.GetContent() + "\n"
		// fmt.Fprintln(conn, item.GetTimestamp())
		// fmt.Fprintln(conn, item)
	}
	// if the last character is a newline, remove it
	ret = strings.TrimRight(ret, "\n")
	fmt.Fprintln(conn, ret)
}
