package requests

import (
	"fmt"
	"gohub/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// 验证函数类型
type ValidatorFunc func(interface{}, *gin.Context) map[string][]string

func Validate(c *gin.Context, obj interface{}, handler ValidatorFunc) bool {
	// 1. 解析请求
	if err := c.ShouldBind(obj); err != nil {
		response.BadRequest(c, err)

		fmt.Println(err.Error())
		return false
	}

	// 2. 表单验证
	errs := handler(obj, c)

	// 3. 判断验证是否通过
	if len(errs) > 0 {

		response.ValidateError(c, errs)
		return false
	}

	return true
}

func validate(data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	// 配置选项
	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		Messages:      messages,
		TagIdentifier: "valid",
	}

	return govalidator.New(opts).ValidateStruct()
}
