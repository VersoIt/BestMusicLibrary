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
	http.HandleFunc("/song/get", h.GetSongs)
	http.HandleFunc("/song/add", h.AddSong)
	http.HandleFunc("/song/delete", h.DeleteSong)
	http.HandleFunc("/song/update", h.UpdateSong)
	http.HandleFunc("/song/verses", h.GetSongVerses)
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	logrus.Error(err)
}
