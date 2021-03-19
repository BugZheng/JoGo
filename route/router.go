package route

import (
	"JoGo/app/api"
	_ "JoGo/docs"
	"JoGo/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"os"
)

//RegisterApp 路由注册
func RegisterApp(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "☺ welcome to golang app")
	})
	// 探侦地址，用于健康检查
	r.HEAD("/listen", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})
	//在正式环境无需要这个路由开启
	if env := os.Getenv("DEPLOY_ENV"); env == "dev" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	//使用use + 中间件 就可以使用中间件的功能
	v1 := r.Group("/api/v1").Use(middlewares.Recovery()).Use(middlewares.RateLimiter())
	{
		v1.GET("demo", api.Demo)
		v1.GET("demoApi", api.DemoAPI)
	}
}
