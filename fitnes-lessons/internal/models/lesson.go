package models

type Lesson struct {
	LessonId       int64
	Title          string
	Time           string
	DayOfWeek      int32
	TrainerId      int64
	AvailableSeats int32
	Description    string
	Difficult      string
	CurUsersCount  int32
	Users          []int64
}

func NewLesson(title string, time string, dayOfWeek int32, trainerId int64, availableSeats int32, desc string, diff string) *Lesson {
	return &Lesson{Title: title, Time: time, DayOfWeek: dayOfWeek, TrainerId: trainerId, AvailableSeats: availableSeats, Description: desc, Difficult: diff, Users: make([]int64, 0)}
}
