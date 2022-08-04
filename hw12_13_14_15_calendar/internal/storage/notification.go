package storage

import "time"

type Notification struct {
	ID      int
	UserID  int
	Title   string
	EventAt time.Time
}
