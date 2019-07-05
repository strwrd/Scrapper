package model

import (
	"time"
)

// Archieve Object
type Archieve struct {
	ID        string
	Link      string
	Code      string
	Published string
	Journals  []*Journal
	CreatedAt time.Time
	UpdatedAt time.Time
}
