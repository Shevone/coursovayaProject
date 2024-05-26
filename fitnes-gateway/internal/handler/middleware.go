package handler

import (
	"fitnes-gateway/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	emptyString         = ""
	space               = " "
	userIdCtx           = "userId"
	roleTrainer         = "Trainer"
)

func (h *Handler) userIdentity(c *gin.Context) {
	const op = "middleware.userIdentity"

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

	claims, err := h.service.AccountService.ParseToken(headerParts[1])
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrResponse{Message: err.Error()})
		return
	}

	c.Set("userId", claims.Id)
	c.Set("userName", claims.Name)
	c.Set("userRole", claims.Role)
	c.Set("userEmail", claims.Email)
}
func (h *Handler) validateIsAdmin(c *gin.Context) {
	const op = "middleware.isAdmin"
	header := c.GetHeader(authorizationHeader)

	// Валидируем заголовок
	if header == emptyString {
		h.logger.Error("%s: %w", op, "Authorization header is empty")
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrResponse{Message: "empty auth header"})
		return
	}

	headerParts := strings.Split(header, space)
	if len(headerParts) != 2 {
		h.logger.Error("%s: %w", op, "Authorization header has invalid formal")
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrResponse{Message: "Authorization header has invalid formal"})
		return
	}

	// Парсим токен

	claims, err := h.service.AccountService.ParseToken(headerParts[1])
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrResponse{Message: err.Error()})
		return
	}

	if claims.Role != roleAdmin {
		errMsg := "User has no admin rules"
		h.logger.Error("%s: %w", op, errMsg)
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrResponse{Message: errMsg})
		return
	}
}
func (h *Handler) validateIsAdminOrTrainer(c *gin.Context) {
	const op = "middleware.isAdminOrTrainer"
	header := c.GetHeader(authorizationHeader)

	// Валидируем заголовок
	if header == emptyString {
		h.logger.Error("%s: %w", op, "Authorization header is empty")
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrResponse{Message: "empty auth header"})
		return
	}

	headerParts := strings.Split(header, space)
	if len(headerParts) != 2 {
		h.logger.Error("%s: %w", op, "Authorization header has invalid formal")
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrResponse{Message: "Authorization header has invalid formal"})
		return
	}

	// Парсим токен

	claims, err := h.service.AccountService.ParseToken(headerParts[1])
	if err != nil {
		h.logger.Error("%s: %w", op, err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrResponse{Message: err.Error()})
		return
	}

	if claims.Role != roleAdmin && claims.Role != roleTrainer {
		errMsg := "User has no admin rules"
		h.logger.Error("%s: %w", op, errMsg)
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrResponse{Message: errMsg})
		return
	}
}
func (h *Handler) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Разрешить доступ с вашего домена
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:63342")
		// Другие необходимые заголовки CORS, если нужно
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")

		// Продолжаем выполнение следующего middleware или обработчика
		c.Next()
	}
}
