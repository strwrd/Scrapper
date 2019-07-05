package model

import (
	"time"
)

// Journal Object
type Journal struct {
	ID         string
	ArchieveID string
	Title      string
	Authors    string
	Abstract   string
	Link       string
	PDFLink    string
	Published  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
