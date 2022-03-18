package auth

import (
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/models/user"
	"gohub/app/requests"
	"gohub/pkg/response"
)

type PasswordController struct {
	v1.BaseAPIController
}

func (pc *PasswordController) ResetByPhone(c *gin.Context) {
	// 验证表单
	request := requests.ResetByPhoneRequest{}

	if ok := requests.Validate(c, &request, requests.ResetByPhone); !ok {
		return
	}

	// 更新密码
	userModel := user.GetByPhone(request.Phone)
	if userModel.ID == 0 {
		response.Abort404(c)
	} else {
		userModel.Password = request.Password
		userModel.Save()

		response.Success(c)
	}
}

func (pc *PasswordController) ResetByEmail(c *gin.Context) {
	request := requests.ResetByEmailRequest{}

	if ok := requests.Validate(c, &request, requests.ResetByEmail); !ok{
		return
	}

	// 2. 更新密码
	userModel := user.GetByEmail(request.Email)
	if userModel.ID == 0 {
		response.Abort404(c)
	} else {
		userModel.Password = request.Password
		userModel.Save()
		response.Success(c)
	}}