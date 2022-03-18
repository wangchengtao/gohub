package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"gohub/pkg/config"
	"gohub/pkg/logger"
	"gohub/pkg/redis"
	"strings"
)

func GetKeyIP(c *gin.Context) string {
	return c.ClientIP()
}

// 路由+IP, 针对单个路由做限流
func GetKeyRouteWithIP(c *gin.Context) string {
	return routeToKeyString(c.FullPath()) + GetKeyIP(c)
}

// 检查请求是否超额
func CheckRate(c *gin.Context, key string, formatted string) (limiter.Context, error) {
	// 实例化依赖的 limiter 包的 limiter Rate 对象
	var context limiter.Context
	rate, err := limiter.NewRateFromFormatted(formatted)
	if err != nil {
		logger.LogIf(err)
		return context, nil
	}

	// 初始化存储, 使用redis 对象
	store, err := sredis.NewStoreWithOptions(redis.Redis.Client, limiter.StoreOptions{
		Prefix: config.GetString("app.name") + ":limiter",
	})

	if err != nil {
		logger.LogIf(err)
		return context, err
	}

	limiterObj := limiter.New(store, rate)

	if c.GetBool("limiter-once") {
		return limiterObj.Peek(c, key)
	} else {
		c.Set("limiter-once", true)

		return limiterObj.Get(c, key)
	}

}

func routeToKeyString(routeName string) string{
	routeName = strings.ReplaceAll(routeName, "/", "-")
	routeName = strings.ReplaceAll(routeName, ":", "_")

	return routeName
}