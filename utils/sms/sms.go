package sms

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"unicode/utf8"
	"wolfdog/internal/consts"
	"wolfdog/utils/common"
)

const (
	SMSTPL = "【xxxx】您正在申请手机注册，验证码为：[code]，若非本人操作请忽略！"
	//账号
	ACCOUNT = "***************"
	//密码
	PSWD = "***************"
	// 发送url，
	URL = "xxxxxxxxxxxxxxxxxxx"
)

func SmsCheck(key, code string) bool {
	key = consts.RedisSuf + key
	val := consts.RedisDB.Do(context.Background(), "GET", key).String()
	if val != code {
		return false
	}
	return true
}
func SmsSet(key, val string) (err error) {
	key = consts.RedisSuf + key
	cmd := consts.RedisDB.Do(context.Background(), "Set", key, val, "EX", 600)
	if err := cmd.Err(); err != nil {
		// 处理错误，如记录日志、返回错误给调用方等
		log.Printf("Redis command failed: %v", err)
		return err
	}
	return
}

func HttpPostForm(url string, data url.Values) (string, error) {

	resp, err := http.PostForm(url, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// 发送短信
func SendSms(mobile, msg string) error {
	if mobile == "" {
		return errors.New("mobile is not null")
	}
	if !common.CheckMobile(mobile) {
		return errors.New("mobile is irregular")
	}
	if utf8.RuneCountInString(msg) < 10 {
		return errors.New("Character length is not enough.")
	}
	//不同信道参数可能不同，具体查看其开发文档
	data_send := url.Values{"account": {ACCOUNT}, "pswd": {PSWD}, "mobile": {mobile}, "msg": {msg}}
	_, err := HttpPostForm(URL, data_send)
	return err
}
