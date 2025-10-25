package payment

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain/constant"
	"github.com/musicman-backend/internal/http/dto"
	"log/slog"
	"net/http"
)

type Service interface {
	CreatePayment(ctx context.Context, returnURI string, userUUID uuid.UUID, amount int) (string, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

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

	redirect, err := h.service.CreatePayment(ctx, req.ReturnURI, userUUID, req.Amount)
	if err != nil {
		slog.Error("failed to create payment", slog.String("err", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, redirect)
}

func (h *Handler) GetPayments(ctx *gin.Context) {

}
