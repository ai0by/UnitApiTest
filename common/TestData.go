package common

import "errors"

//统一返回数据
type RetData struct {
	Code float64 `json:"code"`
	Msg string `json:"msg"`
	Data LoginReturnData `json:"data"`
	Custom interface{} `json:"custom"`
}
//登录数据
type Login struct {
	Boss string `json:"boss"`
	Password string `json:"password"`
}
// 单一用户数据
type LoginReturnData struct {
	Name string `json:"name"`
	Openid string `json:"openid"`
	Status int `json:"status"`
	Token string `json:"token"`
	Uid int `json:"uid"`
}

func (res RetData) CheckCodeEqual() error {
	if res.Code == 0.00 {
		IfError(errors.New("LoginTest 错误的返回信息 : " + res.Msg))
		return errors.New("Not a successful request")
	}
	return nil
}