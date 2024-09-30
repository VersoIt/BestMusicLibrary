package handler

import (
	"BestMusicLibrary/internal/service"
	"BestMusicLibrary/pkg/model"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *service.Service
}

func (h *Handler) InitRoutes() {
	http.HandleFunc("/song/details", h.GetSongs)
	http.HandleFunc("/song/add", h.AddSong)
	http.HandleFunc("/song/delete", h.DeleteSong)
	http.HandleFunc("/song/edit", h.UpdateSong)
	http.HandleFunc("/song/text", h.GetSongText)
}

func (h *Handler) GetSongs(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	songs, err := h.service.Song.GetSongs(group, song, page, limit)
	if err != nil {
		handleError(w, err)
		return
	}
	err = json.NewEncoder(w).Encode(songs)
	if err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AddSong(w http.ResponseWriter, r *http.Request) {
	var song model.Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		handleError(w, err)
		return
	}
	err = h.service.Song.AddSong(song)
	if err != nil {
		handleError(w, err)
		return
	}

	err = h.service.Song.AddSong(song)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	err = h.service.Song.DeleteSong(int64(id))
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	var song model.Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetSongText(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		handleError(w, err)
		return
	}

	text, err := h.service.Song.GetSongText(int64(id), page, limit)
	if err != nil {
		handleError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(text)
	if err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	log.Println(err)
}
