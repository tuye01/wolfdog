package dto

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.String(http.StatusOK, "hello world!!!")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "hello world!!!",
	})
}

type UserMobile struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required"`
	Passwd string `form:"passwd" json:"passwd" binding:"required,max=20,min=6"`
	Code   string `form:"code" json:"code" binding:"required,len=6"`
}
type UserMobileCode struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required"`
	Code   string `form:"code" json:"code" binding:"required,len=6"`
}

type UserMobilePasswd struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required"`
	Passwd string `form:"passwd" json:"passwd" binding:"required,max=20,min=6"`
}

type Mobile struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required"`
}
