package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gohub/pkg/app"
	"gohub/pkg/limiter"
	"gohub/pkg/logger"
	"gohub/pkg/response"
	"net/http"
)

func LimitIP(limit string) gin.HandlerFunc {
	if app.IsTesting() {
		limit = "1000000-H"
	}

	return func(c *gin.Context) {
		// 针对 IP 限流
		key := limiter.GetKeyIP(c)
		if ok := limiterHandler(c, key, limit); !ok {
			return
		}
		c.Next()
	}
}

func LimitPerRoute(limit string) gin.HandlerFunc {
	if app.IsTesting() {
		limit = "1000000-H"
	}

	return func(c *gin.Context) {
		// 针对单个路由, 增加访问次数
		c.Set("limiter-once", false)

		// 针对 IP+路由限制
		key := limiter.GetKeyRouteWithIP(c)
		if ok := limiterHandler(c, key, limit); !ok {
			return
		}
		c.Next()
	}
}

func limiterHandler(c *gin.Context, key string, limit string) bool {
	// 获取超额的情况
	rate, err := limiter.CheckRate(c, key, limit)
	if err != nil {
		logger.LogIf(err)
		response.Abort500(c)
		return false
	}

	// 设置标头
	c.Header("X-RateLimit-Limit", cast.ToString(rate.Limit))         //最大方位次数
	c.Header("X-RateLimit-Remaining", cast.ToString(rate.Remaining)) // 剩余访问次数
	c.Header("X-RateLimit-Reset", cast.ToString(rate.Reset))         // 到该时间点会重置 limit

	// 超额
	if rate.Reached {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"message": "接口请求太频繁",
		})
		return false
	}

	return true
}
