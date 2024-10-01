package model

import "time"

type Song struct {
	Id          int64
	Group       string
	Name        string
	ReleaseDate time.Time
	Verses      []Verse
	Link        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Verse struct {
	VerseNumber int    `json:"verse_number"`
	Text        string `json:"text"`
}
