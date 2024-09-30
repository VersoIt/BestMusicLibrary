package model

import "time"

type Song struct {
	id          int64
	Group       string
	Name        string
	ReleaseDate string
	Text        string
	Link        string
	CreatedAt   time.Time
}
