package main

import (
	"fmt"
	"gohub/bootstrap"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.New()

	bootstrap.SetupRoute(router)

	err := router.Run(":3000")

	if err != nil {
		fmt.Println(err)
	}
}
