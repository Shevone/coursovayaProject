package handler

import (
	"fitnes-gateway/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	roleAdmin = "Admin"
	roleNew   = "New"
)

func (h *Handler) register(c *gin.Context) {
	const op = "accountHandler.register"
	var inputUser models.User

	// Парсим тело json в inputUser
	if err := c.BindJSON(&inputUser); err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: err.Error()})
		return
	}
	id, err := h.service.AccountService.CreateUser(c, inputUser)
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrResponse{Message: err.Error()})
		return
	}
	h.logger.Info("user inited")
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
func (h *Handler) login(c *gin.Context) {
	const op = "accountHandler.login"
	var inputLogin models.LoginInput

	// Парсим тело json в input
	if err := c.BindJSON(&inputLogin); err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: err.Error()})
		return
	}
	token, err := h.service.AccountService.LoginUser(c, inputLogin)
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (h *Handler) profile(c *gin.Context) {
	const op = "accountHandler.profile"

	userId := c.GetInt64(userIdCtx)

	user, err := h.service.AccountService.GetUserProfile(c, userId)
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id":           user.Id,
		"email":        user.Email,
		"name":         user.Name,
		"surname":      user.Surname,
		"patronymic":   user.Patronymic,
		"phone_number": user.PhoneNumber,
		"role":         user.Role,
	})

}
func (h *Handler) editProfile(c *gin.Context) {
	const op = "accountHandler.editProfile"
	var inputUser models.EditProfileRequest

	// Парсим тело json в inputUser
	if err := c.BindJSON(&inputUser); err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: "invalid data"})
		return
	}
	_, err := inputUser.GetId()
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: "wrong user id"})
		return
	}
	resultMessage, err := h.service.AccountService.EditUserProfile(c, inputUser)
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": resultMessage,
	})
}
func (h *Handler) tokenValidation(c *gin.Context) {
	const op = "accountHandler.tokenValidation"
	header := c.GetHeader(authorizationHeader)

	// Валидируем заголовок
	if header == emptyString {
		h.logger.Error("%s: %w", op, "Authorization header is empty")
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrResponse{Message: "empty auth header"})
		c.Redirect(301, "/account/login")
		return
	}

	headerParts := strings.Split(header, space)
	if len(headerParts) != 2 {
		h.logger.Error("%s: %w", op, "Authorization header has invalid formal")
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrResponse{Message: "Authorization header has invalid formal"})
		return
	}

	// Парсим токен

	_, err := h.service.AccountService.ParseToken(headerParts[1])
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"valid": "ok",
	})

}
func (h *Handler) getUsers(c *gin.Context) {
	const op = "accountHandler.getUsers"

	page, err := strconv.Atoi(c.Query("page"))
	limit, err := strconv.Atoi(c.Query("limit"))

	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: err.Error()})
		return
	}
	input := models.PaginateRequest{Page: int64(page), Limit: int64(limit)}

	result, err := h.service.AccountService.GetUsers(c, input)
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrResponse{Message: err.Error()})
		return
	}

	response := models.PaginateResponse[models.User]{
		CurPage: input.Page,
		Limit:   input.Limit,
		ElCount: len(result),
		List:    result,
	}
	if len(result) < int(input.Limit) {
		// Если мы ответили числом меньшим чем запрошеное количество
		response.NextPage = input.Page
	} else {
		response.NextPage = input.Page + 1
	}
	if response.CurPage == 0 {
		response.PrePage = 0
	} else {
		response.PrePage = response.CurPage - 1
	}

	c.JSON(http.StatusOK, response)

}
func (h *Handler) updateUserRole(c *gin.Context) {
	const op = "accountHandler.updateUserRole"

	var input models.EditUserRoleRequest

	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: err.Error()})
		return
	}
	result, err := h.service.AccountService.EditUserRole(c, input)
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
func (h *Handler) updateUserPassword(c *gin.Context) {
	const op = "accountHandler.updateUserProfile"
	var input models.UpdatePasswordRequest

	if err := c.BindJSON(&input); err != nil {

		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrResponse{Message: err.Error()})
		return
	}
	result, err := h.service.AccountService.EditUserPassword(c, input)
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
