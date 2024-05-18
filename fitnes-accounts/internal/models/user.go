package models

type User struct {
	ID          int64
	Email       string
	PassHash    []byte
	Name        string
	Surname     string
	Patronymic  string
	Role        string
	PhoneNumber string
}
