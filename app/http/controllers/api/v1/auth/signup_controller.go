package auth

import (
	"fmt"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/models/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignupController struct {
	v1.BaseAPIController
}

func (sc *SignupController) IsPhoneExist(c *gin.Context) {
	// 请求对象
	type IsPhoneExistReqeust struct {
		Phone string `json:"phone"`
	}

	request := IsPhoneExistReqeust{}

	// 解析 json 请求
	if err := c.ShouldBindJSON(&request); err != nil {
		// 解析失败, 返回 422 状态码和错误信息
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})

		// 打印错误信息
		fmt.Println(err.Error())
		// 出错了, 中断请求
		return
	}

	// 检查数据库并返回响应
	c.JSON(http.StatusOK, gin.H{
		"exist": user.IsPhoneExist(request.Phone),
	})
}
