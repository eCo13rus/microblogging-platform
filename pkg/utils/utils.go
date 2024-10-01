package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"microblogging-platform/internal/models"
	"regexp"
	"time"
)

func ValidateUser(user *models.User) error {
	if user.Username == "" {
		return errors.New("имя пользователя не может быть пустым")
	}
	if len(user.Username) < 3 || len(user.Username) > 50 {
		return errors.New("имя пользователя должно быть от 3 до 50 символов")
	}
	if user.Email == "" {
		return errors.New("email не может быть пустым")
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(user.Email) {
		return errors.New("некорректный формат email")
	}
	if user.Password == "" {
		return errors.New("пароль не может быть пустым")
	}
	if len(user.Password) < 8 {
		return errors.New("пароль должен содержать минимум 8 символов")
	}
	return nil
}

func GenerateJWTToken(userID uint) (string, error) {
	mySigningKey := []byte("AllYourBase")
	claims := models.MyCustomClaims{
		UserId: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(mySigningKey)
}
