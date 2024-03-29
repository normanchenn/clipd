package poll

import (
	"fmt"
	"os"
	"time"

	"github.com/normanchenn/clipd/daemon/clipboard"
	"github.com/normanchenn/clipd/daemon/history"
)

func Poll(clipboard_history *history.History, filepath string, interval time.Duration, threshold int, permissions os.FileMode) {
	item := clipboard_history.GetItem(0)
	prev := ""
	if item != nil {
		prev = item.GetContent()
	}

	for {
		cur, err := clipboard.GetClipboard()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting clipboard")
			continue
		}

		if cur != prev {
			fmt.Fprintln(os.Stdout, "New clipboard: ", cur)
			clipboard_history.AddItem(cur, time.Now())
			err = clipboard.WriteClipboardToFile(filepath, permissions, cur)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error writing to file")
				continue
			}
			prev = cur
		}
		time.Sleep(interval)
	}
}
