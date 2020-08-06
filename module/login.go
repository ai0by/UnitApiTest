package module

import (
	"GoTest/common"
	"encoding/json"
	_ "errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)
func LoginTest(requestUrl string,db *gorm.DB) []common.LoginReturnData {
	var login = []common.Login{}
	apiAddress := "/api/Token/login"               // 设置 api接口的地址
	result := common.ReadJson("./json/login.json") // 设置当前接口的json配置路径
	needTimeOut := 0.00                            // 设置需要的超时时间
	err := json.Unmarshal([]byte(result), &login)
	common.IfError(err)
	var data []common.LoginReturnData
	failNum := 0
	successNum :=0
	for key := range login {
		b, err := json.Marshal(login[key])
		common.IfError(err)
		var res []byte
		common.IfError(err)
		res,timeOut,err := common.HttpPost(requestUrl+apiAddress, common.GetRolePath(apiAddress,db), b)
		//res,timeOut,err := HttpPost(requestUrl+apiAddress,"login", b)
		common.IsTimeOut(timeOut,needTimeOut,apiAddress)
		common.IfError(err)
		//resData := make(map[string]interface{})
		var resData common.RetData
		json.Unmarshal([]byte(strings.TrimSpace(string(res))), &resData)
		if resData.CheckCodeEqual() != nil {
			fmt.Println("Login Failed Plz Check err.log")
			failNum++
			continue
		} else{
			fmt.Println("Login Success Get Data ： ",resData.Data)
			successNum++
			data = append(data,resData.Data)
		}
	}
	fmt.Printf("LoginTest Test OK %d Success %d Failed \n", successNum, failNum)
	return data
}

