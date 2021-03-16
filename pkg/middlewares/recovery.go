package middlewares

import (
	"JoGo/boot"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
)

//Recovery 中间件捕获异常的panic统一返回
func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			// panic 捕获
			if err := recover(); err != nil {
				boot.ZapLogger.Error(fmt.Sprintf("goroutine panic: %v", err),
					zap.String("request_id", ctx.GetHeader("request_id")),
					zap.String("stack", string(debug.Stack())),
				)
				ctx.JSON(http.StatusOK, gin.H{
					"success": false,
					"code":    http.StatusInternalServerError,
					"msg":     "服务器错误，请稍后重试",
				})
				//终止Context 传递
				ctx.Abort()
				return
			}
		}()
		ctx.Next()
	}
}
