package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/musicman-backend/internal/domain/constant"
	"github.com/musicman-backend/internal/domain/entity"
	"github.com/musicman-backend/internal/http/dto"
	token_extracter "github.com/musicman-backend/pkg/token-extracter"
)

type TokenVerifier interface {
	VerifyToken(ctx context.Context, tokenString string) (entity.JWTClaims, error)
}

func AuthMiddleware(verifier TokenVerifier) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := token_extracter.ExtractToken(ctx.Request)
		if err != nil {
			slog.Warn("invalid auth header",
				slog.String("err", err.Error()),
				slog.String("header", ctx.Request.Header.Get("Authorization")),
			)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.NewApiError("некорекктный заголовок авторизации"))
		}

		claims, err := verifier.VerifyToken(ctx, token)
		if err != nil {
			slog.Warn("invalid token",
				slog.String("token", token),
				slog.String("err", err.Error()),
			)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.NewApiError("некорекктный токен авторизации"))
			return
		}

		ctx.Set(constant.CtxUserUUID, claims.UserUUID.String())
		ctx.Set(constant.CtxUserLogin, claims.Login)
		ctx.Set(constant.CtxSubscribeStatus, claims.Subscribe)
		ctx.Next()
	}
}
