package middleware

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"microblogging-platform/internal/models"
	"microblogging-platform/pkg/logger"
	"net/http"
	"strings"
)

func Logging(l logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Info("Получен запрос: " + r.Method + " " + r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

func Recovery(l logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					l.Error("Паника: ", err.(error))
					http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("CORS middleware: %s %s", r.Method, r.URL.Path)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Auth middleware: начало обработки запроса")

			authHeader := r.Header.Get("Authorization")
			log.Printf("Auth middleware: полученный заголовок Authorization: %s", authHeader)

			if authHeader == "" {
				log.Println("Auth middleware: отсутствует токен авторизации")
				http.Error(w, "Отсутствует токен авторизации", http.StatusUnauthorized)
				return
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				log.Println("Auth middleware: неверный формат токена")
				http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
				return
			}

			tokenString := bearerToken[1]
			log.Printf("Auth middleware: парсинг токена: %s", tokenString)

			token, err := jwt.ParseWithClaims(tokenString, &models.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte("AllYourBase"), nil
			})

			if err != nil {
				log.Printf("Auth middleware: ошибка при парсинге токена: %v", err)
				http.Error(w, "Неверный токен", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(*models.MyCustomClaims); ok && token.Valid {
				fmt.Printf("%v %v", claims.UserId, claims.StandardClaims.ExpiresAt)
				ctx := context.WithValue(r.Context(), "user_id", claims.UserId)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				log.Println("Auth middleware: неверный токен")
				http.Error(w, "Неверный токен", http.StatusUnauthorized)
			}

		})
	}
}
