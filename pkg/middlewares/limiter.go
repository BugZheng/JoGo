/**
 * @Author: BugZheng
 * @Description: 限流中间件
 * @File:  limiter
 * @Version: 1.0.0
 * @Date: 2021/03/19 11:23 上午
 */
package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"time"
)

func RateLimiter() gin.HandlerFunc {
	// 创建速率配置
	rate := limiter.Rate{
		Period: time.Second,
		Limit:  10,
	}
	// 将数据存入内存
	store := memory.NewStore()
	// 创建速率实例, 必须是真实的请求
	instance := limiter.New(store, rate, limiter.WithTrustForwardHeader(true))
	// 生成gin中间件
	return mgin.NewMiddleware(instance)
}
