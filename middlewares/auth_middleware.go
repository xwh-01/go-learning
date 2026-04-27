package middlewares

import (
	"exchangeapp/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			utils.HandleAppError(ctx, utils.ErrMissingAuthHeader)
			ctx.Abort()
			return
		}
		username, err := utils.ParseJWT(token)

		if err != nil {
			utils.HandleAppError(ctx, utils.ErrInvalidToken)
			ctx.Abort()
			return
		}

		ctx.Set("username", username)
		ctx.Next()
	}
}
