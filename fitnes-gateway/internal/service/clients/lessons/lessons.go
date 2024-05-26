package lessons

import (
	"context"
	"fitnes-gateway/internal/models"
	"fmt"
	pb "github.com/Shevone/proto-fitnes/gen/go/lessons"
	"github.com/gin-gonic/gin"
	grpcrlog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

var (
	LevelEasy   = "Easy"
	LevelMedium = "Medium"
	LevelHard   = "Hard"
)

type LessonsService struct {
	api pb.LessonsServiceClient
	log *slog.Logger
}

func (l LessonsService) GetLessons(ctx *gin.Context, weekDay int) ([]models.Lesson, error) {
	const op = "grpc.getLessons"

	getLessonsResponse, err := l.api.GetAllLessons(ctx, &pb.GetLessonsRequest{DayOfWeek: int32(weekDay)})
	if err != nil {
		l.log.Error(err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	lessonList := make([]models.Lesson, 0, cap(getLessonsResponse.Lessons))
	for _, lessonFromReq := range getLessonsResponse.Lessons {
		lesson := models.Lesson{
			LessonId:       lessonFromReq.LessonId,
			Title:          lessonFromReq.Title,
			Time:           lessonFromReq.Time,
			DayOfWeek:      lessonFromReq.DayOfWeek,
			TrainerId:      lessonFromReq.TrainerId,
			AvailableSeats: lessonFromReq.AvailableSeats,
			Difficult:      lessonFromReq.Difficulty.String(),
			FreeSeats:      lessonFromReq.SeatsCount,
		}
		lessonList = append(lessonList, lesson)
	}
	return lessonList, nil
}

func (l LessonsService) GetLessonById(ctx *gin.Context, lessonId int64) (*models.Lesson, error) {
	const op = "grpc.GetLesson"

	getLessonResponse, err := l.api.GetLesson(ctx, &pb.GetLessonRequest{LessonId: lessonId})
	if err != nil {
		l.log.Error(err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	lesson := models.Lesson{
		LessonId:       getLessonResponse.LessonId,
		Title:          getLessonResponse.Title,
		Time:           getLessonResponse.Time,
		DayOfWeek:      getLessonResponse.DayOfWeek,
		TrainerId:      getLessonResponse.TrainerId,
		AvailableSeats: getLessonResponse.AvailableSeats,
		Difficult:      getLessonResponse.Difficulty.String(),
		FreeSeats:      getLessonResponse.SeatsCount,
	}
	return &lesson, nil
}

func (l LessonsService) GetLessonsByUserId(ctx *gin.Context, request models.PaginateRequest, userId int64) ([]models.Lesson, error) {
	const op = "grpc.getUserLessons"

	getLessonsResponse, err := l.api.GetLessonsByUser(ctx,
		&pb.GetRequest{Page: int32(request.Page), CountEl: int32(request.Limit), GetId: userId})
	if err != nil {
		l.log.Error(err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	lessonList := make([]models.Lesson, 0, cap(getLessonsResponse.Lessons))
	for _, lessonFromReq := range getLessonsResponse.Lessons {
		lesson := models.Lesson{
			LessonId:       lessonFromReq.LessonId,
			Title:          lessonFromReq.Title,
			Time:           lessonFromReq.Time,
			DayOfWeek:      lessonFromReq.DayOfWeek,
			TrainerId:      lessonFromReq.TrainerId,
			AvailableSeats: lessonFromReq.AvailableSeats,
			Difficult:      lessonFromReq.Difficulty.String(),
			FreeSeats:      lessonFromReq.SeatsCount,
		}
		lessonList = append(lessonList, lesson)
	}
	return lessonList, nil
}

func (l LessonsService) GetLessonsByTrainerId(ctx *gin.Context, request models.PaginateRequestWithId) ([]models.Lesson, error) {
	const op = "grpc.getTrainerLessons"

	getLessonsResponse, err := l.api.GetLessonsByTrainer(ctx,
		&pb.GetRequest{Page: int32(request.Page), CountEl: int32(request.Limit), GetId: request.GetId})
	if err != nil {
		l.log.Error(err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	lessonList := make([]models.Lesson, 0, cap(getLessonsResponse.Lessons))
	for _, lessonFromReq := range getLessonsResponse.Lessons {
		lesson := models.Lesson{
			LessonId:       lessonFromReq.LessonId,
			Title:          lessonFromReq.Title,
			Time:           lessonFromReq.Time,
			DayOfWeek:      lessonFromReq.DayOfWeek,
			TrainerId:      lessonFromReq.TrainerId,
			AvailableSeats: lessonFromReq.AvailableSeats,
			Difficult:      lessonFromReq.Difficulty.String(),
			FreeSeats:      lessonFromReq.SeatsCount,
		}
		lessonList = append(lessonList, lesson)
	}
	return lessonList, nil
}

func (l LessonsService) CreateLesson(ctx *gin.Context, newLesson models.Lesson) (int64, error) {
	const op = "grpc.CreateLesson"

	requestLessonModel := &pb.CreateRequest{
		Title:          newLesson.Title,
		Time:           newLesson.Time,
		DayOfWeek:      newLesson.DayOfWeek,
		TrainerId:      newLesson.TrainerId,
		AvailableSeats: newLesson.AvailableSeats,
	}
	switch newLesson.Difficult {
	case LevelEasy:
		requestLessonModel.Difficulty = pb.Difficulty_EASY
	case LevelMedium:
		requestLessonModel.Difficulty = pb.Difficulty_MEDIUM
	case LevelHard:
		requestLessonModel.Difficulty = pb.Difficulty_HARD
	}
	createResponse, err := l.api.CreateLesson(ctx, requestLessonModel)
	if err != nil {
		l.log.Error(err.Error())
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return createResponse.LessonId, nil
}

func (l LessonsService) EditLesson(ctx *gin.Context, editLesson models.Lesson) (int64, error) {
	const op = "grpc.EditLesson"
	requestLessonModel := &pb.Lesson{
		LessonId:       editLesson.LessonId,
		Title:          editLesson.Title,
		Time:           editLesson.Time,
		DayOfWeek:      editLesson.DayOfWeek,
		TrainerId:      editLesson.TrainerId,
		AvailableSeats: editLesson.AvailableSeats,
	}
	switch editLesson.Difficult {
	case LevelEasy:
		requestLessonModel.Difficulty = pb.Difficulty_EASY
	case LevelMedium:
		requestLessonModel.Difficulty = pb.Difficulty_MEDIUM
	case LevelHard:
		requestLessonModel.Difficulty = pb.Difficulty_HARD
	}
	editResponse, err := l.api.EditLesson(ctx, &pb.EditRequest{Lesson: requestLessonModel})
	if err != nil {
		l.log.Error(err.Error())
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return editResponse.LessonId, nil

}

func (l LessonsService) DeleteLesson(ctx *gin.Context, deleteLessonId int64) (bool, error) {
	const op = "grpc.DeleteLesson"

	deleteResponse, err := l.api.DeleteLesson(ctx, &pb.DeleteRequest{LessonId: deleteLessonId})
	if err != nil {
		l.log.Error(err.Error())
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return deleteResponse.Result, err
}

func (l LessonsService) CloseLesson(ctx *gin.Context, lessonId int64) (string, error) {
	const op = "grpc.CloseLesson"

	closeLessonResp, err := l.api.CloseLesson(ctx, &pb.CloseRequest{LessonId: lessonId})
	if err != nil {
		l.log.Error(err.Error())
		return "false", fmt.Errorf("%s: %w", op, err)
	}
	return closeLessonResp.Message, err
}

func (l LessonsService) SignLesson(ctx *gin.Context, lessonId int64, userId int64) (string, error) {
	const op = "grpc.SignLesson"

	signRes, err := l.api.SignUpForLessonOrCancel(ctx, &pb.SignUpOrCancelRequest{LessonId: lessonId, UserId: userId})
	if err != nil {
		l.log.Error(err.Error())
		return "false", fmt.Errorf("%s: %w", op, err)
	}
	return signRes.Message, nil
}

// NewLessonClient конструктор клиента grpc
func NewLessonClient(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*LessonsService, error) {
	const op = "grpc.New"

	// Конфигурируем то, в каких случаях делаем retry
	retryOpts := []grpcretry.CallOption{
		// Коды при которых делаем retry
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		// Максимальное количество попыток новых
		grpcretry.WithMax(uint(retriesCount)),
		// таймаут ретраев
		grpcretry.WithPerRetryTimeout(timeout),
	}
	// Логирование запросов
	logOpt := []grpcrlog.Option{
		// На какие события мы реагируем
		// По сути тело запроса и ответа
		grpcrlog.WithLogOnEvents(grpcrlog.PayloadReceived, grpcrlog.PayloadSent),
	}

	// Подключение к сервису
	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// Оборачиваем 2 interceptors в 1
		// Похоже на middleware
		grpc.WithChainUnaryInterceptor(
			grpcrlog.UnaryClientInterceptor(InterceptorLogger(log), logOpt...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	// Засовываем наше подключение в сервис
	return &LessonsService{
		api: pb.NewLessonsServiceClient(cc),
		log: log,
	}, nil
}

// InterceptorLogger адаптирует slog logger под logger interceptor - a
// Обертка над логгером, чтоб им мог пользоватлься интерсептор

func InterceptorLogger(l *slog.Logger) grpcrlog.Logger {

	// Возвращает функцию, которая будет вызваться внутри интерсептора
	// А в свою очередь в эту функцию мы засунули наш логгер и метод Log

	return grpcrlog.LoggerFunc(
		func(ctx context.Context, level grpcrlog.Level, msg string, fields ...any) {
			l.Log(ctx, slog.Level(level), msg, fields...)
		},
	)
}
