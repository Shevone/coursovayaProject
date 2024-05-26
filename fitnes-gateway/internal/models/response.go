package models

// ErrResponse структура ошибки
type ErrResponse struct {
	Message string `json:"message"`
}

// PaginateResponse то, в каком виде отдадим пользователю все
type PaginateResponse[T any] struct {
	CurPage  int64 `json:"cur_page"`
	NextPage int64 `json:"next_page"`
	PrePage  int64 `json:"pre_page"`
	Limit    int64 `json:"limit"`
	ElCount  int   `json:"el_count"`
	List     []T   `json:"list"`
}
type LessonsWeekDayResponse[T any] struct {
	Count int `json:"count"`
	List  []T `json:"list"`
}
