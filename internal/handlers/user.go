package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
	"wolfdog/api/dto"
	"wolfdog/internal/consts"
	models2 "wolfdog/internal/models"
	"wolfdog/utils/common"
	"wolfdog/utils/handle"
	"wolfdog/utils/request"
	"wolfdog/utils/response"
	"wolfdog/utils/sms"
	"wolfdog/utils/verify"
)

// 手机密码
func Login(c *gin.Context) {
	var userMobile dto.UserMobilePasswd
	if err := c.BindJSON(&userMobile); err != nil {
		msg := handle.TransTagName(&consts.UserMobileTrans, err)
		response.ShowValidatorError(c, msg)
		return
	}
	model := models2.Users{Mobile: userMobile.Mobile}
	if has := model.GetRow(); !has {
		response.ShowError(c, "mobile_not_exists")
		return
	}
	if common.Sha1En(userMobile.Passwd+model.Salt) != model.Passwd {
		response.ShowError(c, "login_error")
		return
	}
	err := verify.GenerateToken(c, model.Id, model.Name, model.Mobile)
	if err != nil {
		response.ShowError(c, "fail")
		return
	}

	response.ShowSuccess(c, "success")
	return
}

// 注销登录
func Logout(c *gin.Context) {
	token, err := verify.ParseToken(c)
	if err != nil {
		return
	}
	_, err = consts.RedisDB.Del(c, strconv.FormatInt(token.UserID, 10)).Result()
	if err != nil {
		return
	}
	return
}

// 手机验证码登录
func LoginByMobileCode(c *gin.Context) {
	var userMobile dto.UserMobileCode
	if err := c.BindJSON(&userMobile); err != nil {
		msg := handle.TransTagName(&consts.UserMobileTrans, err)
		fmt.Println(msg)
		response.ShowValidatorError(c, msg)
		return
	}
	//验证code
	//if sms.SmsCheck("code"+userMobile.Mobile, userMobile.Code) {
	//	response.ShowError(c, "code_error")
	//	return
	//}
	model := models2.Users{Mobile: userMobile.Mobile}
	if has := model.GetRow(); !has {
		response.ShowError(c, "mobile_not_exists")
		return
	}
	err := verify.GenerateToken(c, model.Id, model.Name, model.Mobile)
	if err != nil {
		response.ShowError(c, "fail")
		return
	}
	response.ShowSuccess(c, "success")
	return
}
func MobileIsExists(c *gin.Context) {
	var userMobile dto.Mobile
	if err := c.BindJSON(&userMobile); err != nil {
		msg := handle.TransTagName(&consts.UserMobileTrans, err)
		fmt.Println(msg)
		response.ShowValidatorError(c, msg)
		return
	}
	if !common.CheckMobile(userMobile.Mobile) {
		response.ShowError(c, "mobile_error")
		return
	}
	model := models2.Users{Mobile: userMobile.Mobile}
	var data = map[string]bool{"is_exist": true}
	if has := model.GetRow(); !has {
		data["is_exist"] = false
	}
	response.ShowData(c, data)
	return
}

// 发送短信验证码
func SendSms(c *gin.Context) {
	var p dto.Mobile
	if err := c.BindJSON(&p); err != nil {
		msg := handle.TransTagName(&consts.MobileTrans, err)
		response.ShowValidatorError(c, msg)
		return
	}
	if !common.CheckMobile(p.Mobile) {
		response.ShowError(c, "mobile_error")
		return
	}
	//生成随机数
	code := common.GetRandomNum(6)
	msg := strings.Replace(sms.SMSTPL, "[code]", code, 1)
	err := sms.SendSms(p.Mobile, msg)
	if err != nil {
		response.ShowError(c, "fail")
		return
	}
	response.ShowError(c, "success")
	return

}

// 手机号注册
func SignupByMobile(c *gin.Context) {
	var userMobile dto.UserMobile
	if err := c.BindJSON(&userMobile); err != nil {
		msg := handle.TransTagName(&consts.UserMobileTrans, err)
		fmt.Println(msg)
		response.ShowValidatorError(c, msg)
		return
	}
	model := models2.Users{Mobile: userMobile.Mobile}
	if has := model.GetRow(); has {
		response.ShowError(c, "mobile_exists")
		return
	}
	//验证code
	if sms.SmsCheck("code"+userMobile.Mobile, userMobile.Code) {
		response.ShowError(c, "code_error")
		return
	}

	model.Salt = common.GetRandomBoth(4)
	model.Passwd = common.Sha1En(userMobile.Passwd + model.Salt)
	model.Ctime = int(time.Now().Unix())
	model.Status = models2.UsersStatusOk
	model.Mtime = time.Now()

	traceModel := models2.Trace{Ctime: model.Ctime}
	traceModel.Ip = common.IpStringToInt(request.GetClientIp(c))
	traceModel.Type = models2.TraceTypeReg

	deviceModel := models2.Device{Ctime: model.Ctime, Ip: traceModel.Ip, Client: c.GetHeader("User-Agent")}
	_, err := model.Add(&traceModel, &deviceModel)
	if err != nil {
		fmt.Println(err)
		response.ShowError(c, "fail")
		return
	}
	response.ShowSuccess(c, "success")
	return
}

// access token 续期
func Renewal(c *gin.Context) {
	model, err := verify.ParseToken(c)
	if err != nil {
		return
	}
	consts.RedisDB.Del(c, strconv.FormatInt(model.UserID, 10))
	err = verify.GenerateToken(c, model.UserID, model.UserName, model.Mobile)
	if err != nil {
		response.ShowError(c, "fail")
		return
	}
}
func Info(c *gin.Context) {
	uid := c.Query("uid")
	uIdInt, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		response.ShowErrorParams(c, err.Error())
		return
	}
	model := models2.Users{}
	model.Id = uIdInt
	row, err := model.GetRowById()
	if err != nil {
		fmt.Println(err)
		response.ShowValidatorError(c, err)
		return
	}
	fmt.Println(row)
	fmt.Println(row.Name)
	//隐藏手机号中间数字
	s := row.Mobile
	row.Mobile = string([]byte(s)[0:3]) + "****" + string([]byte(s)[6:])
	response.ShowData(c, row)
	return
}
