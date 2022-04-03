package routes

import (
	controllers "gohub/app/http/controllers/api/v1"
	"gohub/app/http/controllers/api/v1/auth"
	"gohub/app/http/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")

	v1.Use(middlewares.LimitIP("200-H"))
	{
		authGroup := v1.Group("/auth")
		authGroup.Use(middlewares.LimitIP("1000-H"))
		{
			suc := new(auth.SignupController)
			// 判断手机号是否注册
			authGroup.POST("/signup/phone/exist", middlewares.GuestJWT(), middlewares.LimitPerRoute("60-H"), suc.IsPhoneExist)
			authGroup.POST("/signup/email/exist", middlewares.GuestJWT(), middlewares.LimitPerRoute("60-H"), suc.IsEmailExist)
			authGroup.POST("/signup/using-phone", middlewares.GuestJWT(), suc.SignupUsingPhone)
			authGroup.POST("/signup/using-email", middlewares.GuestJWT(), suc.SignupUsingEmail)

			// 发送验证码
			vcc := new(auth.VerifyCodeController)
			authGroup.POST("/verify-codes/captcha", middlewares.LimitPerRoute("50-H"), vcc.ShowCaptcha)
			authGroup.POST("/verify-codes/phone", middlewares.LimitPerRoute("20-H"), vcc.SendUsingPhone)
			authGroup.POST("/verify-codes/email", middlewares.LimitPerRoute("20-H"), vcc.SendUsingEmail)

			// 登录
			lgc := new(auth.LoginController)
			authGroup.POST("/login/using-phone", middlewares.GuestJWT(), lgc.LoginByPhone)
			authGroup.POST("/login/using-password", middlewares.GuestJWT(), lgc.LoginByPassword)
			authGroup.POST("/login/refresh-token", middlewares.AuthJWT(), lgc.RefreshToken)

			// 重置密码
			pwc := new(auth.PasswordController)
			authGroup.POST("/password-reset/using-phone", middlewares.GuestJWT(), pwc.ResetByPhone)
			authGroup.POST("/password-reset/using-email", middlewares.GuestJWT(), pwc.ResetByEmail)
		}

		uc := new(controllers.UsersController)
		v1.GET("/user", middlewares.AuthJWT(), uc.CurrentUser)
		usersGroup := v1.Group("/users")
		{
			usersGroup.GET("", uc.Index)
			usersGroup.PUT("", middlewares.AuthJWT(), uc.UpdateProfile)
			usersGroup.PUT("/email", middlewares.AuthJWT(), uc.UpdateEmail)

		}

		cgc := new(controllers.CategoriesController)
		cgcGroup := v1.Group("/categories")
		{
			cgcGroup.GET("", cgc.Index)
			cgcGroup.POST("", middlewares.AuthJWT(), cgc.Store)
			cgcGroup.PUT("/:id", middlewares.AuthJWT(), cgc.Update)
			cgcGroup.DELETE("/:id", middlewares.AuthJWT(), cgc.Delete)
			usersGroup.PUT("/email", middlewares.AuthJWT(), uc.UpdateEmail)
			usersGroup.PUT("/phone", middlewares.AuthJWT(), uc.UpdatePhone)
			usersGroup.PUT("/password", middlewares.AuthJWT(), uc.UpdatePassword)

		}

		tpc := new(controllers.TopicsController)
		tpcGroup := v1.Group("/topics")
		{
			tpcGroup.GET("", tpc.Index)
			tpcGroup.POST("", middlewares.AuthJWT(), tpc.Store)
			tpcGroup.PUT("/:id", middlewares.AuthJWT(), tpc.Update)
			tpcGroup.GET("/:id", tpc.Show)
			tpcGroup.DELETE("/:id", middlewares.AuthJWT(), tpc.Delete)
		}

		lsc := new(controllers.LinksController)
		linksGroup := v1.Group("/links")
		{
			linksGroup.GET("", lsc.Index)
		}
	}
}
