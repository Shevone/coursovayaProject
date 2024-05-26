package lessons

import (
	"context"
	"fitnes-lessons/internal/models"
	lessonsFitnesv1 "github.com/Shevone/proto-fitnes/gen/go/lessons"
	"regexp"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type serverApi struct {
	lessonsFitnesv1.UnimplementedLessonsServiceServer
	lessonsService Service
}

type Service interface {
	CreateLesson(ctx context.Context, lesson *models.Lesson) (int64, error)
	DeleteLesson(ctx context.Context, lessonId int64) (bool, error)

	GetLesson(ctx context.Context, lessonId int64) (*models.Lesson, error)
	GetLessonsByTrainerId(ctx context.Context, trainerId int64, page int32, limit int32) ([]*models.Lesson, error)
	GetLessonsByUserId(ctx context.Context, trainerId int64, page int32, limit int32) ([]*models.Lesson, error)
	GetLessonsByWeekDay(ctx context.Context, weekDay int32) ([]*models.Lesson, error)

	SignUpForLessonOrCancel(ctx context.Context, lessonId int64, userId int64) (string, error)
	// CloseLesson(ctx context.Context, lessonId int64, trainerId int64) (bool, error)
	EditLesson(ctx context.Context, lesson *models.Lesson) (int64, error)
}

func Register(gRPCServer *grpc.Server, lessonsService Service) {
	lessonsFitnesv1.RegisterLessonsServiceServer(
		gRPCServer,
		&serverApi{
			lessonsService: lessonsService,
		})
}

func (s serverApi) GetAllLessons(ctx context.Context, request *lessonsFitnesv1.GetLessonsRequest) (*lessonsFitnesv1.GetResponse, error) {
	lessons, err := s.lessonsService.GetLessonsByWeekDay(ctx, request.DayOfWeek)
	if err != nil {
		return nil,
			status.Error(codes.Internal,
				"failed to get lessons")
	}
	response := &lessonsFitnesv1.GetResponse{}
	for _, lesson := range lessons {
		lessonGrpcModel := lessonsFitnesv1.Lesson{
			LessonId:       lesson.LessonId,
			Title:          lesson.Title,
			Time:           lesson.Time,
			TrainerId:      lesson.TrainerId,
			AvailableSeats: lesson.AvailableSeats,
			Description:    lesson.Description,
			DayOfWeek:      lesson.DayOfWeek,
			SeatsCount:     lesson.CurUsersCount,
		}
		switch lesson.Difficult {
		case "HARD":
			lessonGrpcModel.Difficulty = lessonsFitnesv1.Difficulty_HARD
		case "MEDIUM":
			lessonGrpcModel.Difficulty = lessonsFitnesv1.Difficulty_MEDIUM
		case "EASY":
			lessonGrpcModel.Difficulty = lessonsFitnesv1.Difficulty_EASY
		}
		response.Lessons = append(response.Lessons, &lessonGrpcModel)
	}
	return response, nil
}

func (s serverApi) GetLesson(ctx context.Context, request *lessonsFitnesv1.GetLessonRequest) (*lessonsFitnesv1.Lesson, error) {
	if request.GetLessonId() <= 0 {
		return nil,
			status.Error(
				codes.InvalidArgument,
				"lesson id must be greater than 0")
	}
	lesson, err := s.lessonsService.GetLesson(ctx, request.GetLessonId())
	if err != nil {
		return nil,
			status.Error(codes.Internal,
				"failed to get lesson")
	}
	lessonResponse := &lessonsFitnesv1.Lesson{
		LessonId:       lesson.LessonId,
		Title:          lesson.Title,
		Time:           lesson.Time,
		DayOfWeek:      lesson.DayOfWeek,
		TrainerId:      lesson.TrainerId,
		AvailableSeats: lesson.AvailableSeats,
		Description:    lesson.Description,
		SeatsCount:     lesson.CurUsersCount,
	}
	switch lesson.Difficult {
	case "HARD":
		lessonResponse.Difficulty =
			lessonsFitnesv1.Difficulty_HARD
	case "MEDIUM":
		lessonResponse.Difficulty =
			lessonsFitnesv1.Difficulty_MEDIUM
	case "EASY":
		lessonResponse.Difficulty =
			lessonsFitnesv1.Difficulty_EASY
	}
	return lessonResponse, nil
}

func (s serverApi) CreateLesson(ctx context.Context, request *lessonsFitnesv1.CreateRequest) (*lessonsFitnesv1.CreateResponse, error) {
	// Валидируем данные поступившие нам
	if request.TrainerId <= 0 {
		return nil,
			status.Error(codes.InvalidArgument,
				"error with trainer id")
	}
	if request.AvailableSeats <= 0 {
		return nil,
			status.Error(codes.InvalidArgument,
				"the number of seats per session must be greater than 0")
	}
	if !isValidTimeFormat(request.Time) {
		return nil,
			status.Error(
				codes.InvalidArgument,
				"invalid time format")
	}
	if request.Title == "" {
		return nil,
			status.Errorf(codes.InvalidArgument,
				"title is required")
	}
	lesson := models.NewLesson(
		request.GetTitle(),
		request.GetTime(),
		request.GetDayOfWeek(),
		request.GetTrainerId(),
		request.GetAvailableSeats(),
		request.GetDescription(),
		request.GetDifficulty().String(),
	)
	lesson.CurUsersCount = request.AvailableSeats
	lessonId, err := s.lessonsService.CreateLesson(ctx, lesson)
	if err != nil {
		return nil,
			status.Error(codes.Internal,
				"failed to create lesson")
	}
	return &lessonsFitnesv1.CreateResponse{
		LessonId: lessonId,
	}, nil
}
func isValidTimeFormat(input string) bool {
	_, err := time.Parse("15:04", input)
	return err == nil
}
func validateDate(input string) bool {
	// Паттерн для проверки формата даты
	pattern := `^(202[0-9])-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`
	// Компиляция регулярного выражения
	regexp := regexp.MustCompile(pattern)
	// Проверка соответствия строки паттерну
	return regexp.MatchString(input)
}

func (s serverApi) DeleteLesson(ctx context.Context, request *lessonsFitnesv1.DeleteRequest) (*lessonsFitnesv1.DeleteResponse, error) {
	if request.LessonId < 0 {
		return nil,
			status.Error(codes.InvalidArgument,
				"lesson id must be greater than 0")
	}
	result, err := s.lessonsService.DeleteLesson(ctx, request.GetLessonId())
	if err != nil {
		return nil,
			status.Error(codes.Internal,
				"failed to delete lesson")
	}
	return &lessonsFitnesv1.DeleteResponse{
		Result: result,
	}, nil

}

func (s serverApi) EditLesson(ctx context.Context, request *lessonsFitnesv1.EditRequest) (*lessonsFitnesv1.EditResponse, error) {
	// Валидируем данные поступившие нам
	if request.Lesson.TrainerId <= 0 {
		return nil,
			status.Error(codes.InvalidArgument,
				"error with trainer id")
	}
	if request.Lesson.AvailableSeats <= 0 {
		return nil,
			status.Error(codes.InvalidArgument,
				"the number of seats per session must be greater than 0")
	}

	if !isValidTimeFormat(request.Lesson.Time) {
		return nil,
			status.Error(codes.InvalidArgument,
				"invalid time format")
	}
	if request.Lesson.Title == "" {
		return nil,
			status.Errorf(codes.InvalidArgument,
				"title is required")
	}
	lesson := models.NewLesson(request.Lesson.GetTitle(), request.Lesson.GetTime(), request.Lesson.GetDayOfWeek(), request.Lesson.GetTrainerId(), request.Lesson.GetAvailableSeats(), request.Lesson.GetDescription(), request.Lesson.GetDifficulty().String())
	lesson.LessonId = request.Lesson.LessonId

	lessnId, err := s.lessonsService.EditLesson(ctx, lesson)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create lesson")
	}
	return &lessonsFitnesv1.EditResponse{LessonId: lessnId}, nil
}

func (s serverApi) GetLessonsByTrainer(ctx context.Context, request *lessonsFitnesv1.GetRequest) (*lessonsFitnesv1.GetResponse, error) {
	if request.GetId < 0 {
		return nil,
			status.Error(codes.InvalidArgument,
				"trainer id must be greater than 0")
	}
	lessons, err := s.lessonsService.GetLessonsByTrainerId(ctx, request.GetId, request.Page, request.CountEl)
	if err != nil {
		return nil,
			status.Error(codes.Internal,
				"failed to get lesson by trainer id")
	}
	response := &lessonsFitnesv1.GetResponse{}
	for _, lesson := range lessons {
		lessonGrpcModel := lessonsFitnesv1.Lesson{
			LessonId:       lesson.LessonId,
			Title:          lesson.Title,
			Time:           lesson.Time,
			TrainerId:      lesson.TrainerId,
			AvailableSeats: lesson.AvailableSeats,
			Description:    lesson.Description,
			DayOfWeek:      lesson.DayOfWeek,
			SeatsCount:     lesson.CurUsersCount,
		}
		switch lesson.Difficult {
		case "HARD":
			lessonGrpcModel.Difficulty = lessonsFitnesv1.Difficulty_HARD
		case "MEDIUM":
			lessonGrpcModel.Difficulty = lessonsFitnesv1.Difficulty_MEDIUM
		case "EASY":
			lessonGrpcModel.Difficulty = lessonsFitnesv1.Difficulty_EASY
		}
		response.Lessons = append(response.Lessons, &lessonGrpcModel)
	}
	return response, nil

}

func (s serverApi) GetLessonsByUser(ctx context.Context, request *lessonsFitnesv1.GetRequest) (*lessonsFitnesv1.GetResponse, error) {
	if request.GetId < 0 {
		return nil,
			status.Error(codes.InvalidArgument,
				"trainer id must be greater than 0")
	}
	lessons, err := s.lessonsService.GetLessonsByUserId(ctx, request.GetId, request.Page, request.CountEl)
	if err != nil {
		return nil,
			status.Error(codes.Internal,
				"failed to get lesson by trainer id")
	}
	response := &lessonsFitnesv1.GetResponse{}
	for _, lesson := range lessons {
		lessonGrpcModel := lessonsFitnesv1.Lesson{
			LessonId:       lesson.LessonId,
			Title:          lesson.Title,
			Time:           lesson.Time,
			TrainerId:      lesson.TrainerId,
			AvailableSeats: lesson.AvailableSeats,
			Description:    lesson.Description,
			DayOfWeek:      lesson.DayOfWeek,
			SeatsCount:     lesson.CurUsersCount,
		}
		switch lesson.Difficult {
		case "HARD":
			lessonGrpcModel.Difficulty = lessonsFitnesv1.Difficulty_HARD
		case "MEDIUM":
			lessonGrpcModel.Difficulty = lessonsFitnesv1.Difficulty_MEDIUM
		case "EASY":
			lessonGrpcModel.Difficulty = lessonsFitnesv1.Difficulty_EASY
		}
		response.Lessons = append(response.Lessons, &lessonGrpcModel)
	}
	return response, nil
}

func (s serverApi) CloseLesson(ctx context.Context, request *lessonsFitnesv1.CloseRequest) (*lessonsFitnesv1.CloseResponse, error) {
	return nil, status.Error(codes.Unimplemented, "не используется")
	//if request.LessonId < 0 {
	//	return nil,
	//		status.Error(codes.InvalidArgument,
	//			"lesson id must be greater than 0")
	//}
	//if request.TrainerId < 0 {
	//	return nil,
	//		status.Error(codes.InvalidArgument,
	//			"trainer id must be greater than 0")
	//}
	//result, err := s.lessonsService.CloseLesson(ctx, request.LessonId, request.TrainerId)
	//if err != nil {
	//	return nil,
	//		status.Error(codes.Internal,
	//			"failed to close lesson")
	//}
	//var message string
	//if result {
	//	message = "Запись на занятие закрыта"
	//} else {
	//	message = "Запись уже была закрыта"
	//}
	//return &lessonsFitnesv1.CloseResponse{Message: message}, nil

}

func (s serverApi) SignUpForLessonOrCancel(ctx context.Context, request *lessonsFitnesv1.SignUpOrCancelRequest) (*lessonsFitnesv1.SignUpOrCancelResponse, error) {
	if request.LessonId < 0 {
		return nil,
			status.Error(codes.InvalidArgument,
				"lesson id must be greater than 0")
	}
	if request.UserId < 0 {
		return nil,
			status.Error(codes.InvalidArgument,
				"trainer id must be greater than 0")
	}
	message, err := s.lessonsService.SignUpForLessonOrCancel(ctx, request.LessonId, request.UserId)
	if err != nil {
		return nil,
			status.Error(codes.Internal,
				"failed to close lesson")
	}
	return &lessonsFitnesv1.SignUpOrCancelResponse{Message: message}, nil
}

func (s serverApi) mustEmbedUnimplementedLessonsServiceServer() {
	//TODO implement me
	panic("implement me")
}
