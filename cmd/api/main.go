package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"microblogging-platform/internal/api"
	"microblogging-platform/internal/config"
	"microblogging-platform/pkg/database"
	"microblogging-platform/pkg/logger"
)

func main() {
	// Инициализация логгера
	newLogger := logger.NewLogger()

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		newLogger.Fatal("Ошибка загрузки конфигурации", err)
	}

	// Подключение к базе данных
	db, err := database.NewGORMConnection()
	if err != nil {
		newLogger.Fatal("Ошибка подключения к базе данных", err)
	}

	r := gin.Default()

	// Добавь CORS middleware
	defaultConfig := cors.DefaultConfig()
	defaultConfig.AllowOrigins = []string{"http://localhost:8080"}
	defaultConfig.AllowHeaders = append(defaultConfig.AllowHeaders, "Authorization")
	r.Use(cors.New(defaultConfig))

	// Инициализация и запуск сервера
	server := api.NewServer(cfg, newLogger, db)
	if err := server.Run(); err != nil {
		newLogger.Fatal("Ошибка при запуске сервера", err)
	}
}
