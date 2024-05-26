package handler

import (
	"fitnes-gateway/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	userRoleCtx = "userRole"
	roleUser    = "User"
)

func (h *Handler) getAllLessons(c *gin.Context) {
	const op = "lessons.getLessons"

	page, err := strconv.Atoi(c.Query("weekDay"))
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: err.Error()})
		return
	}
	if page < 0 || page > 6 {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: "page must be between 0 and 6"})
		return
	}
	lessonList, err := h.service.LessonService.GetLessons(c, page)
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrResponse{Message: err.Error()})
		return
	}
	response := models.LessonsWeekDayResponse[models.Lesson]{
		List:  lessonList,
		Count: len(lessonList),
	}
	c.JSON(http.StatusOK, response)
}
func (h *Handler) getLesson(c *gin.Context) {
	const op = "lessons.getLesson"

	var input models.RequestWithId

	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: err.Error()})
		return
	}
	lesson, err := h.service.LessonService.GetLessonById(c, input.LessonId)
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, lesson)
}
func (h *Handler) getLessonByTrainer(c *gin.Context) {
	const op = "lessons.getLessonsByTrainer"

	page, err := strconv.Atoi(c.Query("page"))
	limit, err := strconv.Atoi(c.Query("limit"))
	trainerId, err := strconv.Atoi(c.Query("trainerId"))
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: err.Error()})
		return
	}
	input := models.PaginateRequestWithId{Page: int64(page), Limit: int64(limit), GetId: int64(trainerId)}
	lessonList, err := h.service.LessonService.GetLessonsByTrainerId(c, input)
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrResponse{Message: err.Error()})
		return
	}
	response := models.PaginateResponse[models.Lesson]{
		CurPage:  input.Page,
		NextPage: input.Page + 1,
		PrePage:  input.Page - 1,
		Limit:    input.Limit,
		ElCount:  len(lessonList),
		List:     lessonList,
	}
	if len(lessonList) < limit {
		response.NextPage = input.Page
	}
	if response.PrePage < 0 {
		response.PrePage = 0
	}
	c.JSON(http.StatusOK, response)
}
func (h *Handler) getLessonByUser(c *gin.Context) {
	const op = "lessons.getLessonsByUser"

	page, err := strconv.Atoi(c.Query("page"))
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: err.Error()})
		return
	}
	input := models.PaginateRequest{Page: int64(page), Limit: int64(limit)}
	// Получаем из контекста, тк пользователь должен быть авторизован для этого запроса
	userId := c.GetInt64(userIdCtx)

	lessonList, err := h.service.LessonService.GetLessonsByUserId(c, input, userId)
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrResponse{Message: err.Error()})
		return
	}
	response := models.PaginateResponse[models.Lesson]{
		CurPage:  input.Page,
		NextPage: input.Page + 1,
		PrePage:  input.Page - 1,
		Limit:    input.Limit,
		ElCount:  len(lessonList),
		List:     lessonList,
	}
	if len(lessonList) < limit {
		response.NextPage = input.Page
	}
	if response.PrePage < 0 {
		response.PrePage = 0
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) createLesson(c *gin.Context) {
	const op = "lessons.createLesson"

	var input models.Lesson
	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: err.Error()})
		return
	}

	if input.Time == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: "time is empty"})
		return
	}
	createdLessonId, err := h.service.LessonService.CreateLesson(c, input)
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"lesson_id": createdLessonId,
	})
}
func (h *Handler) deleteLesson(c *gin.Context) {
	const op = "lessons.deleteLesson"

	var input models.RequestWithId
	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: err.Error()})
		return
	}

	deleteResult, err := h.service.LessonService.DeleteLesson(c, input.LessonId)
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"deleteResult": deleteResult,
	})
}
func (h *Handler) editLesson(c *gin.Context) {
	const op = "lessons.editLesson"

	var input models.Lesson
	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: err.Error()})
		return
	}

	createdLessonId, err := h.service.LessonService.EditLesson(c, input)
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"lesson_id": createdLessonId,
	})
}

/*
	func (h *Handler) closeLesson(c *gin.Context) {
		const op = "lessons.closeLesson"

		var input models.RequestWithId
		if err := c.BindJSON(&input); err != nil {
			h.logger.Error("%s: %w", op, err)
			c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: err.Error()})
			return
		}

		closeResult, err := h.service.LessonService.CloseLesson(c, input.LessonId)
		if err != nil {
			h.logger.Error("%s: %w", op, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"close_result": closeResult,
		})
	}
*/
func (h *Handler) signOnLesson(c *gin.Context) {
	const op = "lessons.closeLesson"

	var input models.RequestWithId
	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: err.Error()})
		return
	}

	userId := c.GetInt64(userIdCtx)
	userRole := c.GetString(userRoleCtx)
	if userRole != roleUser {
		c.AbortWithStatusJSON(http.StatusForbidden, models.ErrResponse{Message: "Эта функция подходит только обычным пользователям"})
		return
	}
	subMessage, err := h.service.LessonService.SignLesson(c, input.LessonId, userId)
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": subMessage,
	})
}
