package BestMusicLibrary

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

const ReadTimeoutSeconds = 10
const WriteTimeoutSeconds = 10

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    ReadTimeoutSeconds * time.Second,
		WriteTimeout:   WriteTimeoutSeconds * time.Second}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
