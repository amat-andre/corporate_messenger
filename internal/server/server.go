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

func New(port int, logger *slog.Logger) *Server {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Messenger API"))
	})

	return &Server{
		logger: logger,
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  time.Minute,
		},
	}
}

func (serv *Server) Run() error {
	serv.logger.Info("server started", "addr", serv.httpServer.Addr)

	err := serv.httpServer.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (serv *Server) Shutdown(ctx context.Context) error {
	serv.logger.Info("server stopping")
	return serv.httpServer.Shutdown(ctx)
}