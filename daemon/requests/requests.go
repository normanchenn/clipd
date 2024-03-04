package requests

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/normanchenn/clipd/daemon/history"
)

type Request struct {
	Action string         `json:"action"`
	Params map[string]int `json:"params"`
}

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

	var request Request
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error decoding request: ", err)
		return
	}

	reqJSON, _ := json.Marshal(request)
	fmt.Println("Request JSON:", string(reqJSON))
	switch request.Action {
	case "get":
		handleGet(conn, request.Params, clipboard_history)
	default:
		fmt.Fprintf(conn, "Invalid Request: %s", request.Action)
	}
}

func handleGet(conn net.Conn, params map[string]int, clipboard_history *history.History) {
	at_value, at_ok := params["at"]
	last_value, last_ok := params["last"]
	from_value, from_ok := params["from"]
	to_value, to_ok := params["to"]
	if len(params) == 0 { // get most recent (no args)
		item := clipboard_history.GetItem(0)
		returnItems(conn, []*history.HistoryItem{item})
	} else if len(params) == 1 && at_ok { // get n
		item := clipboard_history.GetItem(at_value)
		returnItems(conn, []*history.HistoryItem{item})
	} else if len(params) == 1 && last_ok { // get last n
		items := clipboard_history.GetItemRange(0, last_value-1)
		returnItems(conn, items)
	} else if len(params) == 2 && from_ok && to_ok { // get from n to m
		items := clipboard_history.GetItemRange(from_value, to_value)
		returnItems(conn, items)
	} else {
		fmt.Fprintln(conn, "Invalid request")
	}
}

func returnItems(conn net.Conn, items []*history.HistoryItem) {
	var ret string
	for _, item := range items {
		// fmt.Fprintln(os.Stdout, item.GetContent())
		ret += item.GetContent()
	}
	fmt.Fprint(conn, ret)
}
