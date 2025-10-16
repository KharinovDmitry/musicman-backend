package payment

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) NewPayment(ctx *gin.Context) {

}

func (h *Handler) GetPayments(ctx *gin.Context) {}
