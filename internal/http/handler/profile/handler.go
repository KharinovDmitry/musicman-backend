package profile

import (
	"github.com/gin-gonic/gin"
	"github.com/musicman-backend/internal/domain/constant"
	"github.com/musicman-backend/internal/http/dto"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

// GetMyProfile
// @Summary Получить профиль пользователя
// @Description Возвращает информацию о текущем авторизованном пользователе
// @Tags profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserProfile "Успешное получение профиля"
// @Failure 401 {object} dto.ApiError "Пользователь не авторизован"
// @Router /profile/me [get]
func (h *Handler) GetMyProfile(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dto.UserProfile{
		UUID:            ctx.GetString(constant.CtxUserUUID),
		Login:           ctx.GetString(constant.CtxUserLogin),
		SubscribeStatus: ctx.GetBool(constant.CtxSubscribeStatus),
	})
}
