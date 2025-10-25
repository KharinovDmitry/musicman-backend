package payment

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain/constant"
	"github.com/musicman-backend/internal/domain/entity"
	"github.com/musicman-backend/internal/http/dto"
	"log/slog"
	"net/http"
)

type Service interface {
	CreatePayment(ctx context.Context, returnURI string, userUUID uuid.UUID, amount int) (string, error)
}

type History interface {
	GetPaymentsByUser(ctx context.Context, userUUID uuid.UUID) ([]entity.Payment, error)
}

type Handler struct {
	service  Service
	history  History
	validate *validator.Validate
}

func NewHandler(service Service, history History) *Handler {
	return &Handler{
		service:  service,
		history:  history,
		validate: validator.New(),
	}
}

// NewPayment godoc
// @Summary Создание нового платежа
// @Description Создаёт платёж через YooKassa и перенаправляет пользователя на страницу оплаты.
// @Tags payments
// @Param Authorization header string true "Bearer токен"
// @Param request body dto.CreatePaymentRequest true "Данные для создания платежа. return_uri - ссылка на которую вернуть пользователя после оплаты. amount в копейках"
// @Success 301 {string} string "Redirect — ссылка на YooKassa"
// @Failure 400 {object} dto.ApiError "Невалидное тело запроса"
// @Failure 500 "Ошибка сервера или не удалось создать платёж"
// @Router /payments/new [post]
func (h *Handler) NewPayment(ctx *gin.Context) {
	userUUIDstr := ctx.GetString(constant.CtxUserUUID)
	userUUID, err := uuid.Parse(userUUIDstr)
	if err != nil {
		slog.Error("failed to parse user uuid", slog.String("err", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var req dto.CreatePaymentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		slog.Warn("failed to bind create payment request", slog.String("err", err.Error()))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError("невалидное тело запроса"))
		return
	}

	err = h.validate.Struct(req)
	if err != nil {
		slog.Warn("failed to validate create payment request", slog.String("err", err.Error()))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(err.Error()))
		return
	}

	redirect, err := h.service.CreatePayment(ctx, req.ReturnURI, userUUID, req.Amount)
	if err != nil {
		slog.Error("failed to create payment", slog.String("err", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, redirect)
}

// GetPayments godoc
// @Summary Получить список платежей пользователя
// @Description Возвращает историю платежей текущего авторизованного пользователя
// @Tags payments
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Success 200 {array} dto.UserPayment "Список платежей пользователя"
// @Failure 500 "Ошибка сервера или не удалось получить платежи"
// @Router /payments/history [get]
func (h *Handler) GetPayments(ctx *gin.Context) {
	userUUIDstr := ctx.GetString(constant.CtxUserUUID)
	userUUID, err := uuid.Parse(userUUIDstr)
	if err != nil {
		slog.Error("failed to parse user uuid", slog.String("err", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	payments, err := h.history.GetPaymentsByUser(ctx, userUUID)
	if err != nil {
		slog.Error("failed to get user payments", slog.String("err", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, dto.UserPaymentsFromEntities(payments))
}
