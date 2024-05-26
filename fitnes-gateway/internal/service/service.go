package service

import (
	"fitnes-gateway/internal/models"
	"github.com/gin-gonic/gin"
)

type Accounts interface {
	// Public Methods
	CreateUser(ctx *gin.Context, user models.User) (int64, error)
	LoginUser(ctx *gin.Context, loginInput models.LoginInput) (string, error)
	GetUserProfile(ctx *gin.Context, userId int64) (*models.User, error)
	EditUserProfile(ctx *gin.Context, request models.EditProfileRequest) (string, error)

	// For admins
	GetUsers(ctx *gin.Context, request models.PaginateRequest) ([]models.User, error)
	EditUserRole(ctx *gin.Context, request models.EditUserRoleRequest) (string, error)
	EditUserPassword(ctx *gin.Context, request models.UpdatePasswordRequest) (string, error)

	ParseToken(accessToken string) (*models.TokenClaims, error)
}

type Lessons interface {
	GetLessons(ctx *gin.Context, weekDay int) ([]models.Lesson, error)
	GetLessonById(ctx *gin.Context, lessonId int64) (*models.Lesson, error)
	GetLessonsByUserId(ctx *gin.Context, request models.PaginateRequest, userId int64) ([]models.Lesson, error)
	GetLessonsByTrainerId(ctx *gin.Context, request models.PaginateRequestWithId) ([]models.Lesson, error)

	CreateLesson(ctx *gin.Context, newLesson models.Lesson) (int64, error)
	EditLesson(ctx *gin.Context, editLesson models.Lesson) (int64, error)
	DeleteLesson(ctx *gin.Context, deleteLessonId int64) (bool, error)

	CloseLesson(ctx *gin.Context, lessonId int64) (string, error)

	SignLesson(ctx *gin.Context, lessonId int64, userId int64) (string, error)
}

type Service struct {
	AccountService Accounts
	LessonService  Lessons
}

func NewService(accountSrv Accounts, lessonSrv Lessons) *Service {
	return &Service{
		AccountService: accountSrv,
		LessonService:  lessonSrv,
	}
}
