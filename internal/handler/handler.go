package handler

import (
	"BestMusicLibrary/internal/service"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func (h *Handler) InitRoutes() {
	http.HandleFunc("/songs/get", h.GetSongs)
	http.HandleFunc("/songs/add", h.AddSong)
	http.HandleFunc("/songs/delete", h.DeleteSong)
	http.HandleFunc("/songs/update", h.UpdateSong)
	http.HandleFunc("/songs/verses", h.GetSongVerses)
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	logrus.Error(err)
}
