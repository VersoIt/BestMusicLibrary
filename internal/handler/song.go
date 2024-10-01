package handler

import (
	"BestMusicLibrary/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) GetSongs(w http.ResponseWriter, r *http.Request) {
	type songResponse struct {
		Id          int64     `json:"id"`
		Group       string    `json:"group"`
		Name        string    `json:"name"`
		ReleaseDate time.Time `json:"release_date"`
		Link        string    `json:"link"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	if err := handleRequestMethod(w, http.MethodGet, r.Method); err != nil {
		logrus.Error(err)
		return
	}

	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	logrus.WithFields(logrus.Fields{
		"group": group,
		"song":  song,
		"page":  page,
		"limit": limit,
	}).Info("received query parameters")

	pageNum, limitNum, err := parsePagingData(page, limit)

	if err != nil {
		handleError(w, err)
		logrus.WithFields(logrus.Fields{
			"page":  page,
			"limit": limit,
		}).Error("error parsing paging data")
		return
	}

	logrus.WithFields(logrus.Fields{
		"page":  pageNum,
		"limit": limitNum,
	}).Info("parsed paging data")

	songs, err := h.service.Song.GetSongs(group, song, pageNum, limitNum)
	if err != nil {
		handleError(w, err)
		logrus.WithFields(logrus.Fields{
			"group": group,
			"song":  song,
		}).Error("error fetching songs")
		return
	}

	logrus.WithFields(logrus.Fields{
		"count": len(songs),
	}).Info("fetched songs")

	songResponses := make([]songResponse, 0, len(songs))
	for _, s := range songs {
		songResponses = append(songResponses, songResponse{
			Id:          s.Id,
			Group:       s.Group,
			Name:        s.Name,
			ReleaseDate: s.ReleaseDate,
			Link:        s.Link,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		})
	}

	err = json.NewEncoder(w).Encode(songResponses)
	if err != nil {
		handleError(w, err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"response_count": len(songResponses),
	}).Info("response successfully sent")
}

type newSongRequest struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

func (h *Handler) AddSong(w http.ResponseWriter, r *http.Request) {
	if err := handleRequestMethod(w, http.MethodPost, r.Method); err != nil {
		logrus.Error(err)
		return
	}
	var songRequest newSongRequest
	err := json.NewDecoder(r.Body).Decode(&songRequest)
	if err != nil {
		handleError(w, err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"song":  songRequest.Song,
		"group": songRequest.Group,
	}).Info("decoded request body")

	songId, err := h.service.Song.AddSong(model.Song{Name: songRequest.Song, Group: songRequest.Group})
	if err != nil {
		handleError(w, err)
		return
	}

	logrus.WithField("songId", songId).Info("song successfully added")

	w.WriteHeader(http.StatusCreated)
	_, err = fmt.Fprint(w, fmt.Sprintf("%d", songId))
	if err != nil {
		handleError(w, err)
	}

	logrus.WithField("songId", songId).Info("response successfully sent")
}

func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	if err := handleRequestMethod(w, http.MethodDelete, r.Method); err != nil {
		logrus.Error(err)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	logrus.WithFields(logrus.Fields{
		"id": id,
	}).Info("received query parameters")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error(err)
		return
	}

	err = h.service.Song.DeleteSong(int64(id))
	if err != nil {
		handleError(w, err)
		return
	}
	logrus.WithField("id", id).Info("song successfully deleted")
	w.WriteHeader(http.StatusOK)
	logrus.Info("response successfully sent")
}

func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {

	type songUpdate struct {
		Id          int64     `json:"id"`
		Group       string    `json:"group"`
		Name        string    `json:"name"`
		ReleaseDate time.Time `json:"release_date"`
		Text        string    `json:"text"`
		Link        string    `json:"link"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	if err := handleRequestMethod(w, http.MethodPut, r.Method); err != nil {
		logrus.Error(err)
		return
	}
	var song songUpdate
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		handleError(w, err)
		return
	}
	logrus.WithFields(logrus.Fields{
		"id":           song.Id,
		"name":         song.Name,
		"group":        song.Group,
		"release_date": song.ReleaseDate,
		"link":         song.Link,
		"created_at":   song.CreatedAt,
		"updated_at":   song.UpdatedAt,
	}).Info("decoded request body for song update")

	err = h.service.Song.UpdateSong(model.Song{
		Id:          song.Id,
		Group:       song.Group,
		Name:        song.Name,
		ReleaseDate: song.ReleaseDate,
		Link:        song.Link,
		CreatedAt:   song.CreatedAt,
		UpdatedAt:   song.UpdatedAt,
	}, song.Text)

	if err != nil {
		handleError(w, err)
		return
	}

	logrus.WithField("id", song.Id).Info("song successfully updated")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetSongVerses(w http.ResponseWriter, r *http.Request) {
	if err := handleRequestMethod(w, http.MethodGet, r.Method); err != nil {
		logrus.Error(err)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	logrus.WithFields(logrus.Fields{
		"id":    id,
		"page":  page,
		"limit": limit,
	}).Info("parsed query parameters")

	pageNum, limitNum, err := parsePagingData(page, limit)

	if err != nil {
		handleError(w, err)
		return
	}

	text, err := h.service.Song.GetSongVerses(int64(id), pageNum, limitNum)
	if err != nil {
		handleError(w, err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"id":    id,
		"page":  pageNum,
		"limit": limitNum,
	}).Info("successfully retrieved song verses")

	err = json.NewEncoder(w).Encode(text)
	if err != nil {
		handleError(w, err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"id":     id,
		"verses": len(text), // Assuming text is a slice or array
	}).Info("response successfully sent")
}

func handleRequestMethod(w http.ResponseWriter, requiredMethod, currentMethod string) error {
	if currentMethod != requiredMethod {
		errText := fmt.Sprintf("method %s required!", requiredMethod)
		http.Error(w, errText, http.StatusMethodNotAllowed)
		return errors.New(errText)
	}
	return nil
}

func parsePagingData(page, limit string) (pageNum, limitNum int, err error) {
	pageNum = 0
	limitNum = 0
	if page != "" {
		pageNum, err = strconv.Atoi(page)
		if err != nil {
			return
		}
	}
	if limit != "" {
		limitNum, err = strconv.Atoi(limit)
		if err != nil {
			return
		}
	}
	return
}
