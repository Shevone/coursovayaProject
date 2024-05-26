package models

import (
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-playground/validator/v10"
)

// User структура пользователя
// Применяется при регистрации
type User struct {
	Id          int64  `json:"id"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Role        string `json:"role" validate:"oneof=Admin User Trainer"`
	PhoneNumber string `json:"phoneNumber"`
}

// LoginInput структура для авторизации пользователя
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TokenClaims структура, куда записываем данные из jwt токена
type TokenClaims struct {
	jwt.StandardClaims
	Id    int64 `json:"uid"`
	Email string
	Name  string
	Role  string
}
