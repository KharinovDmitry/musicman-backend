package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/musicman-backend/internal/domain"
	"github.com/musicman-backend/internal/http/dto"
)

type Auth interface {
	Login(ctx context.Context, login, password string) (string, error)
	Register(ctx context.Context, login, password string) (string, error)
}

type Handler struct {
	auth Auth
}

func NewHandler(auth Auth) *Handler {
	return &Handler{
		auth: auth,
	}
}

// Login
// @Summary Аутентификация пользователя
// @Description Выполняет вход пользователя в систему и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Данные для входа"
// @Success 200 {object} dto.LoginResponse "Успешная аутентификация"
// @Failure 400 {object} dto.ApiError "Неверный формат запроса"
// @Failure 401 {object} dto.ApiError "Неверные учетные данные"
// @Failure 500 {object} dto.ApiError "Внутренняя ошибка сервера"
// @Router /auth/sign-in [post]
func (h *Handler) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		slog.Warn("invalid request", slog.String("err", err.Error()))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError("некорекктное тело запроса"))
		return
	}

	token, err := h.auth.Login(ctx, req.Username, req.Password)
	if err != nil && !errors.Is(err, domain.ErrInvalidCredentials) {
		slog.Error("failed ti login", slog.String("err", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if errors.Is(err, domain.ErrInvalidCredentials) {
		slog.Warn("invalid credentials", slog.String("err", err.Error()))
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.NewApiError("Невереный логин или пароль"))
		return
	}

	ctx.JSON(http.StatusOK, dto.LoginResponse{Token: token})
}

// Register обрабатывает запрос на регистрацию нового пользователя
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Данные для регистрации"
// @Success 200 {object} dto.RegisterResponse "Успешная регистрация"
// @Failure 400 {object} dto.ApiError "Неверный формат запроса"
// @Failure 409 {object} dto.ApiError "Пользователь уже существует"
// @Failure 500 {object} dto.ApiError "Внутренняя ошибка сервера"
// @Router /auth/sign-up [post]
func (h *Handler) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBind(&req); err != nil {
		slog.Warn("invalid request", slog.String("err", err.Error()))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError("некорекктное тело запроса"))
		return
	}

	token, err := h.auth.Register(ctx, req.Username, req.Password)
	if err != nil && !errors.Is(err, domain.ErrUserAlreadyExists) {
		slog.Error("failed to register", slog.String("err", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if errors.Is(err, domain.ErrUserAlreadyExists) {
		slog.Warn("user already exist", slog.String("err", err.Error()))
		ctx.AbortWithStatusJSON(http.StatusConflict, dto.NewApiError("Пользователь с таким логином уже существует"))
		return
	}

	ctx.JSON(http.StatusOK, dto.RegisterResponse{Token: token})
}
