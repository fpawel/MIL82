package data

import (
	"time"
)

//go:generate reform

// Entry represents a row in entry table.
//reform:entry
type Entry struct {
	EntryID   int64     `reform:"entry_id,pk"`
	CreatedAt time.Time `reform:"created_at"`
	Message   string    `reform:"message"`
	Level     string    `reform:"level"`
	WorkID    int64     `reform:"work_id"`
}
