package common

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type NowTime struct {
	StrTime string
	UnixTime int
	UnixTimeNano int
}

type JccUserRole struct{
	Id 			int64 			`json:"id"`
	Sid			int64 			`json:"sid" description:"超级ID"`
	Pid 		int64 			`json:"pid" description:"父级Id"`
	Path    	string 			`json:"path" description:"前端路由"`
	Key     	string 			`json:"key" description:"后台方法接口"`
	Name     	string 			`json:"name" description:"名称"`
	IsStart    	int64 			`json:"is_start" description:"是否启用"`
	IsDel    	int64 			`json:"is_del" description:"是否删除"`
	IsShow    	int64 			`json:"is_show" description:"是否隐藏"`
	Level    	int64 			`json:"level" description:"描述"`
}

//通过 url jcc-path body 发送post请求 返回结果返回
func HttpPost(url,jccpath string, data []byte) ([]byte,float64,error){
	body := bytes.NewReader(data)

	request, err := http.NewRequest("POST", url, body)

	if err != nil {
		fmt.Println("请求错误500REQ01:",err.Error())
		return  nil, 0,err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("JCC-Path", jccpath)

	beforTime := GetNowTime()
	resp, err := (&http.Client{}).Do(request)
	afterTime := GetNowTime()
	timeOutSec := GetTimeOut(beforTime,afterTime)

	code := resp.StatusCode
	if code != 200 {
		return nil,timeOutSec,errors.New("请求错误500REQ04 请求状态码错误"+string(code))
	}
	if err != nil {
		fmt.Println("请求错误500REQ02", err)
		return nil,timeOutSec, err
	}
	defer resp.Body.Close()

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("请求错误500REQ03", err)
		return nil,timeOutSec, err
	}

	return respByte,timeOutSec, nil
}
//读取json文件 返回结果
func ReadJson(filePath string) (result string){
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	buf := bufio.NewReader(file)
	for {
		s, err := buf.ReadString('\n')  // 按行读取
		result += s
		if err != nil {
			if err == io.EOF{   // 判断是否为文件尾
				break
			}else{
				fmt.Println("ERROR:", err)
				return
			}
		}
	}
	return result
}
//报错增加
func IfError(err error){
	t:= GetNowTime()
	var content string
	content = "["+t.StrTime+"] "
	if err != nil {
		err = ErrLogWright("./err.log",content+err.Error()+"\n")
	}
	if err != nil {
		fmt.Println("写入到文件出错")
	}
}
// 文件写入
func ErrLogWright(fileName string, content string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("cacheFileList.yml file create failed. err: " + err.Error())
	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
	return err
}
// 取当前时间日期
func GetNowTime()NowTime{
	current := time.Now()
	var t NowTime
	t.StrTime = current.Format("2006-01-02 15:04:05")
	t.UnixTime = int(current.Unix())
	t.UnixTimeNano = int(current.UnixNano())
	return t
}
// 记录超时时间
func GetTimeOut(befor,after NowTime) float64 {
	return float64(after.UnixTimeNano)/1000000000.00 - float64(befor.UnixTimeNano)/1000000000.00
}
// 判断超时时间
func IsTimeOut(fact,need float64,content string){
	need = need * 1e9
	if fact >= need {
		_ = ErrLogWright("./timeout.log", fmt.Sprintf("[%s]接口超时 *%s* 当前使用 %.10f 秒 规定使用 %.10f 秒\n",GetNowTime().StrTime,content,fact,need))
	}
}
// 根据接口地址获取 path 试用 可能有BUG
func GetRolePath(apiRole string,db *gorm.DB) string {
	apiRole = "api/money/moneyCalculate"
	var userPath JccUserRole
	db.Where("`key` like ?", "%"+apiRole+"%").First(&userPath)
	//db.First(&userPath)
	return userPath.Path
}