package main

import (
	"flag"
	"fmt"
	"gohub/bootstrap"
	btsConfig "gohub/config"
	"gohub/pkg/config"

	"github.com/gin-gonic/gin"
)

func init() {
	btsConfig.Initialize()
}

func main() {
	// 配置初始化 依赖命令行 --env 参数
	var env string
	flag.StringVar(&env, "env", "", "加载 .env 文件")
	flag.Parse()
	config.InitConfig(env)

	// 初始化 DB
	bootstrap.SetupDB()

	// new 一个 Gin Engine 实例
	router := gin.New()

	// 初始化路由绑定
	bootstrap.SetupRoute(router)
	fmt.Println(config.Get("app.port"))

	// 运行服务
	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		fmt.Println(err)
	}
}
