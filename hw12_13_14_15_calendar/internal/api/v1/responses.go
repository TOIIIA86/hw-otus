package v1

import "time"

type EventResponse struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       time.Time  `json:"end_at"`
	NotifyAt    *time.Time `json:"notify_at"`
}

type EventListResponse struct {
	Events []EventResponse `json:"events"`
}
