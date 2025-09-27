package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h Handler) Health(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
