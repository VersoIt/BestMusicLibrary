package service

import (
	"BestMusicLibrary/internal/model"
	"BestMusicLibrary/internal/repository"
	"context"
	"strings"
	"time"
)

type SongFetchData struct {
	ReleaseDate string
	Text        string
	Link        string
}

type SongDataFetcher interface {
	FetchSongDetails(group, song string) (SongFetchData, error)
}

type SongService struct {
	songRepos       repository.Song
	songDataFetcher SongDataFetcher
}

const (
	defaultPagePagingAmount  = 0
	defaultLimitPagingAmount = 5
)

func NewSongService(repos repository.Song, songFetcher SongDataFetcher) *SongService {
	return &SongService{songRepos: repos, songDataFetcher: songFetcher}
}

// GetSongs Получение данных библиотеки с фильтрацией по всем полям и пагинацией
func (s *SongService) GetSongs(group, song string, rawPage, rawLimit int) ([]model.Song, error) {
	page, limit := handlePagingData(rawPage, rawLimit)
	return s.songRepos.GetSongs(group, song, page, limit)
}

// GetSongVerses Получение текста песни с пагинацией по куплетам
func (s *SongService) GetSongVerses(id int64, rawPage, rawLimit int) ([]model.Verse, error) {
	page, limit := handlePagingData(rawPage, rawLimit)
	return s.songRepos.GetSongVerses(id, page, limit)
}

// DeleteSong Удаление песни
func (s *SongService) DeleteSong(id int64) error {
	return s.songRepos.DeleteSong(id)
}

// UpdateSong Изменение песни
func (s *SongService) UpdateSong(song model.Song, text string) error {
	song.Verses = textToVerses(text)
	return s.songRepos.UpdateSong(song)
}

// AddSong Добавление песни
func (s *SongService) AddSong(song model.Song) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	enrichedSong, err := s.enrichSongWithAPI(ctx, song)
	if err != nil {
		return 0, err
	}

	return s.songRepos.AddSong(enrichedSong)
}

// EnrichSongWithAPI Обогащение данных с использованием стороннего сервиса
func (s *SongService) enrichSongWithAPI(ctx context.Context, song model.Song) (enrichedSong model.Song, err error) {
	errChan := make(chan error)
	resChan := make(chan SongFetchData)

	go func() {
		songDetails, clientError := s.songDataFetcher.FetchSongDetails(song.Group, song.Name)
		if clientError != nil {
			errChan <- clientError
			return
		}
		resChan <- songDetails
	}()

	var songDetails SongFetchData
	select {
	case clientError := <-errChan:
		return model.Song{}, clientError
	case res := <-resChan:
		songDetails = res
	case <-ctx.Done():
		return model.Song{}, ctx.Err()
	}

	if err != nil {
		return song, err
	}

	layout := "02.01.2006"
	releaseDate, err := time.Parse(layout, songDetails.ReleaseDate)

	enrichedSong.ReleaseDate = releaseDate
	enrichedSong.Verses = textToVerses(songDetails.Text)
	enrichedSong.Link = songDetails.Link

	return enrichedSong, nil
}

func handlePagingData(rawPage, rawLimit int) (page, limit int) {
	page = rawPage
	limit = rawLimit
	if page <= 0 {
		page = defaultPagePagingAmount
	}
	if limit <= 0 {
		limit = defaultLimitPagingAmount
	}
	return
}

func textToVerses(text string) []model.Verse {
	verses := strings.Split(text, "\n\n")
	cleanedVerses := make([]model.Verse, 0)
	count := 0
	for _, verse := range verses {
		cleanedVerse := strings.TrimSpace(verse)
		if cleanedVerse != "" {
			cleanedVerses = append(cleanedVerses, model.Verse{VerseNumber: count, Text: cleanedVerse})
			count++
		}
	}
	return cleanedVerses
}
