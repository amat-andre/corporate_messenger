package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func New(port int, handler http.Handler, logger *slog.Logger) *Server {

	// handler.Handle("GET /", http.FileServer(http.Dir(webDir))) // когда будет готов фронт, раскомментировать и указать путь к нему

	return &Server{
		logger: logger,
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  time.Minute,
		},
	}
}

func (server *Server) Run() error {
	server.logger.Info("server started", "addr", server.httpServer.Addr)

	err := server.httpServer.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (server *Server) Shutdown(ctx context.Context) error {
	server.logger.Info("server stopping")
	return server.httpServer.Shutdown(ctx)
}