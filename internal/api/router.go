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
	// Применяем CORS middleware ко всем маршрутам
	s.router.Use(middleware.CORS)
	s.router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Создаем подмаршрутизатор для API
	api := s.router.PathPrefix("/api").Subrouter()

	// Auth routes
	authHandler := handlers.NewAuthHandler(s.db, s.logger)
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST", "OPTIONS")

	// User routes
	userHandler := handlers.NewUserHandler(s.db, s.logger)
	api.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	// api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")

	// Post routes
	postHandler := handlers.NewPostHandler(s.db, s.logger)
	api.HandleFunc("/posts", postHandler.GetPosts).Methods("GET")
	api.HandleFunc("/posts/{id}", postHandler.GetPost).Methods("GET")

	// Защищенные маршруты
	protected := api.PathPrefix("/").Subrouter()
	protected.Use(middleware.Auth(s.config.JWTSecret))
	protected.HandleFunc("/users/me", userHandler.GetCurrentUser).Methods("GET")

	// Защищенные User routes
	protected.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	protected.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Защищенные Post routes
	protected.HandleFunc("/posts", postHandler.CreatePost).Methods("POST")
	protected.HandleFunc("/posts/{id}", postHandler.UpdatePost).Methods("PUT")
	protected.HandleFunc("/posts/{id}", postHandler.DeletePost).Methods("DELETE")

	// Comment routes
	commentHandler := handlers.NewCommentHandler(s.db, s.logger)
	protected.HandleFunc("/posts/{postId}/comments", commentHandler.GetComments).Methods("GET")
	protected.HandleFunc("/posts/{postId}/comments", commentHandler.CreateComment).Methods("POST")
	protected.HandleFunc("/comments/{id}", commentHandler.UpdateComment).Methods("PUT")
	protected.HandleFunc("/comments/{id}", commentHandler.DeleteComment).Methods("DELETE")

	// Like routes
	likeHandler := handlers.NewLikeHandler(s.db, s.logger)
	protected.HandleFunc("/posts/{postId}/likes", likeHandler.LikePost).Methods("POST")
	protected.HandleFunc("/posts/{postId}/likes", likeHandler.UnlikePost).Methods("DELETE")
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
