package data

import (
	"time"
)

//go:generate reform

// Work represents a row in work table.
//reform:work
type Work struct {
	WorkID    int64     `reform:"work_id,pk"`
	CreatedAt time.Time `reform:"created_at"`
	Name      string    `reform:"name"`
}
