package lib

import (
	"fitnes-account/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user models.User, salt string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Добавляем в токен всю необходимую информацию
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["name"] = user.Name
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(duration).Unix()

	// Подписываем токен, используя секретный ключ приложения

	tokenString, err := token.SignedString([]byte(salt))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
