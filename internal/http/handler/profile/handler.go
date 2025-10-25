package profile

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain/constant"
	"github.com/musicman-backend/internal/domain/entity"
	"github.com/musicman-backend/internal/http/dto"
	"log/slog"
	"net/http"
)

type UserRepo interface {
	GetUserByUUID(ctx context.Context, userUUID uuid.UUID) (entity.User, error)
}

type Handler struct {
	userRepo UserRepo
}

func NewHandler(userRepo UserRepo) *Handler {
	return &Handler{
		userRepo: userRepo,
	}
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
	userUUIDstr := ctx.GetString(constant.CtxUserUUID)
	userUUID, err := uuid.Parse(userUUIDstr)
	if err != nil {
		slog.Error("failed to parse user uuid", slog.String("err", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	profile, err := h.userRepo.GetUserByUUID(ctx, userUUID)
	if err != nil {
		slog.Error("failed to get user", slog.String("err", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, dto.UserProfile{
		UUID:   profile.UUID,
		Login:  profile.Login,
		Tokens: profile.Tokens,
	})
}
