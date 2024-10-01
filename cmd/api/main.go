package main

import (
	"microblogging-platform/internal/api"
	"microblogging-platform/internal/config"
	"microblogging-platform/pkg/database"
	"microblogging-platform/pkg/logger"
)

func main() {
	// Инициализация логгера
	logger := logger.NewLogger()

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Ошибка загрузки конфигурации", err)
	}

	// Подключение к базе данных
	db, err := database.NewGORMConnection()
	if err != nil {
		logger.Fatal("Ошибка подключения к базе данных", err)
	}

	// Инициализация и запуск сервера
	server := api.NewServer(cfg, logger, db)
	if err := server.Run(); err != nil {
		logger.Fatal("Ошибка при запуске сервера", err)
	}
}
