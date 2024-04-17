package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	Log        *zap.Logger
}

func New() *Server {
	return &Server{}
}

func (s *Server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    "localhost:8080",
		Handler: handler,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) WaitForShutDown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.Log.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("shut down failed: %w", err)
	}

	s.Log.Info("Server exiting.")
	return nil
}
