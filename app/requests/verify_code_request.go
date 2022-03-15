package requests

import (
	"gohub/app/requests/validators"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type VerifyCodePhoneRequest struct {
	CaptchaID     string `json:"captcha_id,omitempty" valid:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer,omitempty" valid:"captcha_answer"`

	Phone string `json:"phone,omitempty" valid:"phone"`
}

func VerifyCodePhone(data interface{}, c *gin.Context) map[string][]string {
	// 定制认证规则
	rules := govalidator.MapData{
		"phone":          []string{"required", "digits:11"},
		"captcha_id":     []string{"required"},
		"captcha_answer": []string{"required", "digits:6"},
	}

	// 定制错误消息
	// 2. 定制错误消息
	messages := govalidator.MapData{
		"phone": []string{
			"required:手机号为必填项，参数名称 phone",
			"digits:手机号长度必须为 11 位的数字",
		},
		"captcha_id": []string{
			"required:图片验证码的 ID 为必填",
		},
		"captcha_answer": []string{
			"required:图片验证码答案必填",
			"digits:图片验证码长度必须为 6 位的数字",
		},
	}

	errs := validate(data, rules, messages)

	// 图片验证码
	_data := data.(*VerifyCodePhoneRequest)
	errs = validators.ValidateCaptcha(_data.CaptchaID, _data.CaptchaAnswer, errs)

	return errs

}

type VerifyCodeEmailRequest struct {
	CaptchaID     string `json:"captcha_id,omitempty" valid:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer,omitempty" valid:"captcha_answer"`
	Email         string `json:"email,omitempty" valid:"email"`
}

func VerifyCodeEmail(data interface{}, c *gin.Context) map[string][]string {
	// 定制认证规则
	rules := govalidator.MapData{
		"captcha_id":     []string{"required"},
		"captcha_answer": []string{"required", "digits:6"},
		"email":          []string{"required", "email", "min:4", "max:30"},
	}

	// 定制错误消息
	messages := govalidator.MapData{
		"captcha_id": []string{
			"required:图片验证码的 ID 必填",
		},
		"email": []string{
			"required:邮箱必填",
			"min:Email 长度需大于 4",
			"max:email 长度需大于 30",
			"email:email 格式不正确",
		},
		"captcha_answer": []string{
			"required:图片验证码答案必填",
			"digits:图片验证码长度必须为 6 位的数字",
		},
	}

	errs := validate(data, rules, messages)

	// 图片验证码
	_data := data.(*VerifyCodeEmailRequest)
	errs = validators.ValidateCaptcha(_data.CaptchaID, _data.CaptchaAnswer, errs)

	return errs
}
