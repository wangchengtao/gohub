package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type SignupPhoneExistRequest struct {
	Phone string `json:"phone,omitempty" valid:"phone"`
}

func ValidateSignupPhoneExist(data interface{}, c *gin.Context) map[string][]string {
	// 自定义验证规则
	rules := govalidator.MapData{
		"phone": []string{"required", "digits:11"},
	}

	messages := govalidator.MapData{
		"phone": []string{
			"required:手机号必填, 参数名称 phone",
			"digits:手机号长度必须 11 位数字",
		},
	}

	return validate(data, rules, messages)
}

type SignupEmailExistRequest struct {
	Email string `json:"email,omitempty" valid:"email"`
}

func ValidateSignupEmailExist(data interface{}, c *gin.Context) map[string][]string {
	// 自定义检测规则
	rules := govalidator.MapData{
		"email": []string{
			"required", "min:4", "max:30", "email",
		},
	}

	messages := govalidator.MapData{
		"email": []string{
			"required:邮箱必填",
			"min:长度需大于 4",
			"max:长度需小于 30",
			"email:格式不正确",
		},
	}

	return validate(data, rules, messages)
}
