package middlewares

import (
	"go-server/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")

	bearerString := "Bearer "
	if !strings.HasPrefix(token, bearerString) {
		utils.UnauthorizedAccess(ctx)
		ctx.Abort()
		return
	}

	token = strings.TrimPrefix(token, bearerString)

	if token != "Gib me access" {
		utils.UnauthorizedAccess(ctx, "Please attach <<Gib me access>>")
		ctx.Abort()
		return
	}

	ctx.Next()
}
