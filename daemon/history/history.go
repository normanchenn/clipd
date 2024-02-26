package history

import (
	"container/list"
	"time"
)

type History struct {
	history *list.List
}

type HistoryItem struct {
	content   string
	timestamp time.Time
}

func NewHistory() *History {
	return &History{history: list.New()}
}

func NewHistoryItem(content string, timestamp time.Time) *HistoryItem {
	return &HistoryItem{content: content, timestamp: timestamp}
}

func (history *History) AddItem(item string, time time.Time) {
	historyItem := NewHistoryItem(item, time)
	history.history.PushFront(historyItem)
}

func (history *History) GetItem(index int) *HistoryItem {
	if index < 0 || index >= history.history.Len()-1 {
		return nil
	}
	cur := history.history.Front()
	for i := 0; i < index; i++ {
		cur = cur.Next()
	}
	return cur.Value.(*HistoryItem)
}

func (history *History) GetItemRange(start int, end int) []*HistoryItem {
	if start < 0 || start >= history.history.Len() || end < 0 || end >= history.history.Len()-1 || start > end {
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

func (history_item *HistoryItem) GetContent() string {
	return history_item.content
}

func (history_item *HistoryItem) GetTimestamp() time.Time {
	return history_item.timestamp
}
