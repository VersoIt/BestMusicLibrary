package service

import (
	"BestMusicLibrary/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTextToVerses(t *testing.T) {
	text := "This is verse 1.\n\nThis is verse 2.\n\nThis is verse 3."
	expected := []model.Verse{
		{VerseNumber: 0, Text: "This is verse 1."},
		{VerseNumber: 1, Text: "This is verse 2."},
		{VerseNumber: 2, Text: "This is verse 3."},
	}
	result := textToVerses(text)
	assert.Equal(t, expected, result, "the verses should match the expected output")
}

func TestTextToVersesEmptyLines(t *testing.T) {
	text := "\n\nThis is verse 1.\n\n\n\nThis is verse 2.\n\n"
	expected := []model.Verse{
		{VerseNumber: 0, Text: "This is verse 1."},
		{VerseNumber: 1, Text: "This is verse 2."},
	}
	result := textToVerses(text)
	assert.Equal(t, expected, result, "the verses should handle empty lines correctly")
}

func TestTextToVersesEmptyText(t *testing.T) {
	text := ""
	expected := make([]model.Verse, 0)
	result := textToVerses(text)
	assert.Equal(t, expected, result, "the verses should handle empty lines correctly")
}
