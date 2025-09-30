package middleware

import (
	"beep/internal/errors"
	"log/slog"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery panic错误处理中间件
func Recovery() gin.RecoveryFunc {
	return func(c *gin.Context, err any) {
		if e, ok := err.(error); ok {
			stack := debug.Stack()
			serviceErr, ok := errors.AsServiceError(e)
			if ok {
				// 处理service error
				slog.Error("service error", "error", serviceErr, "detail", serviceErr.Details, "stack", string(stack))
				c.JSON(200, gin.H{
					"code":    serviceErr.Code,
					"message": serviceErr.Message,
				})
				c.Abort()
				return
			}
			slog.Error("recover error", "error", e, "stack", string(stack))
			c.JSON(200, gin.H{
				"code":    500,
				"message": "Internal Server Error",
			})
			c.Abort()
		}
	}
}
