package models

type Lesson struct {
	LessonId       int64
	Title          string
	Time           string
	Date           string
	TrainerId      int64
	AvailableSeats int32
	Description    string
	Difficult      string
}

func NewLesson(title string, time string, date string, trainerId int64, availableSeats int32, desc string, diff string) *Lesson {
	return &Lesson{Title: title, Time: time, Date: date, TrainerId: trainerId, AvailableSeats: availableSeats, Description: desc, Difficult: diff}
}
