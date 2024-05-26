package models

// Lesson структура урока
type Lesson struct {
	LessonId       int64
	Title          string
	Time           string
	DayOfWeek      int32
	TrainerId      int64
	AvailableSeats int32
	Description    string
	Difficult      string
	FreeSeats      int32
	Users          []int64
}
