package handler

import (
	"fitnes-gateway/internal/service"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Handler struct {
	logger  *slog.Logger
	service *service.Service
}

func NewHandler(service *service.Service, log *slog.Logger) *Handler {
	return &Handler{service: service, logger: log}
}

// InitRoutes Инициализируем наши endpoint
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:63342")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-ijt")

		// Проверяем метод запроса на OPTIONS (предварительный запрос)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	})

	auth := router.Group("/account")
	{
		auth.POST("/token-valid", h.tokenValidation)
		auth.POST("/register", h.register) // Регистрация
		auth.POST("/login", h.login)       // Авторизация
		auth.POST("/token-validate", h.userIdentity)

		withAuthorization := auth.Group("/", h.userIdentity)
		{
			withAuthorization.PUT("/edit-password", h.updateUserPassword)
			withAuthorization.POST("/profile", h.profile)         // Профиль
			withAuthorization.PUT("/edit-profile", h.editProfile) // Изменение данных профиля
			adminAuth := withAuthorization.Group("/for-admin", h.validateIsAdmin)
			{
				adminAuth.GET("/users", h.getUsers)
				adminAuth.PUT("/update-role", h.updateUserRole)
			}
		}

	}
	router.GET("/a", h.getAllLessons) // Получить все
	lessons := router.Group("/lessons", h.userIdentity)
	{
		lessons.PUT("/sign", h.signOnLesson)

		forTrainersAdmins := lessons.Group("/", h.validateIsAdminOrTrainer)
		{
			forTrainersAdmins.POST("/create", h.createLesson)
			forTrainersAdmins.DELETE("/delete", h.deleteLesson)
			forTrainersAdmins.PUT("/edit", h.editLesson)
			//forTrainersAdmins.PUT("/close", h.closeLesson) // Закрыть запись
		}

		lessonsGet := lessons.Group("/get")
		{
			lessonsGet.GET("/lesson", h.getLesson)              // Одно занятие
			lessonsGet.GET("/by-trainer", h.getLessonByTrainer) // По id тренера
			lessonsGet.GET("/by-user", h.getLessonByUser)       // По id пользователя
		}

	}

	return router
}
