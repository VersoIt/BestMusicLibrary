package service

import (
	"BestMusicLibrary/internal/client"
	"BestMusicLibrary/internal/model"
	"BestMusicLibrary/internal/repository"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type SongService struct {
	songRepos  repository.Song
	songClient client.ExternalSongApiClient
}

const (
	defaultPagePagingAmount  = 0
	defaultLimitPagingAmount = 5
)

func NewSongService(repos repository.Song, songClient client.ExternalSongApiClient) *SongService {
	return &SongService{songRepos: repos, songClient: songClient}
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
	enrichedSong, err := s.enrichSongWithAPI(song)
	if err != nil {
		return 0, err
	}

	return s.songRepos.AddSong(enrichedSong)
}

// EnrichSongWithAPI Обогащение данных с использованием стороннего сервиса
func (s *SongService) enrichSongWithAPI(song model.Song) (enrichedSong model.Song, err error) {
	logrus.Error("enriching song '%s'\n", song.Name)

	songDetails, err := s.songClient.FetchSongDetails(song.Group, song.Name)
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
	var cleanedVerses []model.Verse
	for index, verse := range verses {
		cleanedVerses = append(cleanedVerses, model.Verse{VerseNumber: index, Text: strings.TrimSpace(verse)})
	}
	return cleanedVerses
}
