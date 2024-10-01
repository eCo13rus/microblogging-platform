package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
	DBHost        string
	DBPort        int
	DBUser        string
	DBPassword    string
	DBName        string
	RedisHost     string
	RedisPort     int
	RedisPassword string
	KafkaBrokers  string
	RabbitMQURL   string
}

func Load() (*Config, error) {
	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))

	return &Config{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        dbPort,
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     redisPort,
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		KafkaBrokers:  os.Getenv("KAFKA_BROKERS"),
		RabbitMQURL:   os.Getenv("RABBITMQ_URL"),
	}, nil
}
