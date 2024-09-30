package handler

import (
	"BestMusicLibrary/internal/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type SongHandler struct {
	service *service.Service
}

func (h *SongHandler) GetSongs(rw http.ResponseWriter, r *http.Request) error {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	songs, err := h.service.Song.GetSongs(group, song, page, pageSize)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get songs", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(songs)
}

func (h *SongHandler) GetSongText(rw http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *SongHandler) DeleteSong(rw http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *SongHandler) EditSong(rw http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *SongHandler) AddSong(rw http.ResponseWriter, r *http.Request) error {
	return nil
}
