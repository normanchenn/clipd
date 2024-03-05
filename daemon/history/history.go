package history

import (
	"container/list"
	"fmt"
	"os"
	"sync"
	"time"
)

type History struct {
	history *list.List
	size    int
	sync.Mutex
}

type HistoryItem struct {
	content   string
	timestamp time.Time
}

func InitHistory() *History {
	return &History{history: list.New(), size: 0}
}

func InitHistoryItem(content string, timestamp time.Time) *HistoryItem {
	return &HistoryItem{content: content, timestamp: timestamp}
}

func (history *History) PrintHistory() {
	fmt.Println(os.Stdout, "-----------------------------")
	for i := history.history.Front(); i != nil; i = i.Next() {
		fmt.Fprintln(os.Stdout, i.Value.(*HistoryItem).GetContent())
	}
	fmt.Println(os.Stdout, "-----------------------------")
}
func (history *History) GetSize() int {
	return history.size
}

func (history_item *HistoryItem) GetContent() string {
	return history_item.content
}

func (history_item *HistoryItem) GetTimestamp() time.Time {
	return history_item.timestamp
}

func (history *History) AddItem(item string, time time.Time) {
	history.Lock()
	defer history.Unlock()

	historyItem := InitHistoryItem(item, time)
	history.history.PushFront(historyItem)
	history.size++
}

func (history *History) GetItem(index int) *HistoryItem {
	history.Lock()
	defer history.Unlock()

	if index < 0 || index > history.history.Len()-1 {
		return nil
	}
	cur := history.history.Front()
	for i := 0; i < index; i++ {
		cur = cur.Next()
	}
	return cur.Value.(*HistoryItem)
}

func (history *History) GetItemRange(start int, end int) []*HistoryItem {
	history.Lock()
	defer history.Unlock()

	if start < 0 || start > history.history.Len()-1 || end < 0 || end > history.history.Len()-1 || start > end {
		return nil
	}
	result := make([]*HistoryItem, 0)
	cur := history.history.Front()
	for i := 0; i < start; i++ {
		cur = cur.Next()
	}
	for i := start; i <= end; i++ {
		result = append(result, cur.Value.(*HistoryItem))
		cur = cur.Next()
	}
	return result
}
