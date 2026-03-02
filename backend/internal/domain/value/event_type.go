package value

import "fmt"

type EventType string

// 有効なイベントタイプの定数定義
const (
	EventTypeLive  EventType = "live"
	EventTypeCamp  EventType = "camp"
	EventTypeOther EventType = "other"
)

var validEventTypes = map[EventType]bool{
	EventTypeLive:  true,
	EventTypeCamp:  true,
	EventTypeOther: true,
}

func NewEventType(s string) (EventType, error) {
	et := EventType(s)
	if validEventTypes[et] {
		return et, nil
	}
	return "", fmt.Errorf("無効なイベントタイプです: %s", s)
}