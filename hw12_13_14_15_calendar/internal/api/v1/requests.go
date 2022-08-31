package v1

import "time"

type EventCreateRequest struct {
	UserID      int     `json:"user_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	StartAt     int64   `json:"start_at"`
	EndAt       int64   `json:"end_at"`
	NotifyAt    *int64  `json:"notify_at"`
}

type EventUpdateRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	StartAt     int64   `json:"start_at"`
	EndAt       int64   `json:"end_at"`
	NotifyAt    *int64  `json:"notify_at"`
}

func GetNotifyAt(notifyAt *int64) *time.Time {
	if notifyAt == nil {
		return nil
	}
	t := time.Unix(*notifyAt, 0)
	return &t
}
