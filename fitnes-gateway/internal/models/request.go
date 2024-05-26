package models

import "strconv"

// PaginateRequest вид запроса с пагинацией
type PaginateRequest struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

// PaginateRequestWithId структура запроса с пагинацией + некое id для нужд
type PaginateRequestWithId struct {
	GetId int64 `json:"get_id"`
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

// EditProfileRequest структура - изменения данных пользователя
type EditProfileRequest struct {
	Id         string `json:"userId"`
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	Patronymic string `json:"patronymic"`
}

func (e *EditProfileRequest) GetId() (int64, error) {
	id, err := strconv.Atoi(e.Id)
	return int64(id), err
}

// EditUserRoleRequest структура запроса для изменения роли пользоватея
type EditUserRoleRequest struct {
	UserId  int64 `json:"user_id"`
	NewRole int64 `json:"new_role"`
}

// RequestWithId структура запроса на получение занятия
type RequestWithId struct {
	LessonId int64 `json:"lesson_id"`
}
type UpdatePasswordRequest struct {
	UserId   int64  `json:"user_id"`
	Password string `json:"password"`
}
