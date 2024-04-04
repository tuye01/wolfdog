package api

import (
	"github.com/gin-gonic/gin"
	"wolfdog/api/dto"
	"wolfdog/internal/handlers"
	"wolfdog/utils/verify"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", dto.Index)
	r.POST("/login", handlers.Login)
	r.POST("/login/mobile", handlers.LoginByMobileCode)
	groupJwt := r.Group("/", verify.AuthMiddleware())
	groupJwt.POST("/renewal", handlers.Renewal)
	groupJwt.POST("/logout", handlers.Logout)
	groupJwt.POST("/sendsms", handlers.SendSms)
	groupJwt.POST("/signup/mobile", handlers.SignupByMobile)
	groupJwt.POST("/signup/mobile/exist", handlers.MobileIsExists)
	groupJwt.GET("/my/info", handlers.Info) //用户信息
	groupJwt.GET("/pong", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
