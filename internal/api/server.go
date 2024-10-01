package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"microblogging-platform/internal/config"
	"microblogging-platform/internal/handlers"
	"microblogging-platform/pkg/logger"
	"microblogging-platform/pkg/middleware"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Server struct {
	router *mux.Router
	logger logger.Logger
	config *config.Config
	db     *gorm.DB
}

func NewServer(cfg *config.Config, logger logger.Logger, db *gorm.DB) *Server {
	s := &Server{
		router: mux.NewRouter(),
		logger: logger,
		config: cfg,
		db:     db,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.Use(middleware.Logging(s.logger))
	s.router.Use(middleware.Recovery(s.logger))
	s.router.Use(middleware.CORS)

	postHandler := handlers.NewPostHandler(s.db, s.logger)
	s.router.HandleFunc("/api/posts", postHandler.GetPosts).Methods("GET")
	s.router.HandleFunc("/api/posts", postHandler.CreatePost).Methods("POST")
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.config.ServerAddress,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Ошибка при запуске сервера", err)
		}
	}()

	s.logger.Info(fmt.Sprintf("Сервер запущен на %s", s.config.ServerAddress))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
