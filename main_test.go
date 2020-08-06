package gotest

import (
	"GoTest/common"
	"GoTest/module"
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "gopkg.in/check.v1"
	"io"
	"testing"
)

var db *gorm.DB                        // 全局数据库
var loginData []common.LoginReturnData // 全局已登录用户数据

type ErrorTest interface {
	CheckCodeEqual(code int) error // 检查 登录状态码 返回错误信息，如无则返回空，否则抛出一个异常
	CheckTimeOut(time int) // 查看接口是否超时
}
type ApiTest interface {
	GetPath(string) // 获取当前的Path
	ErrorTest
}
func Test(t *testing.T) {
	TestingT(t)
	fmt.Println(loginData)
	// 协程处理 testing 框架
	//获取loginData 以便于开启后续的测试任务
	//for key,value := range loginData  {
	//
	//}
}

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestHelloWorld(c *C) {
	c.Assert(42, Equals, 42)
	c.Assert(io.ErrClosedPipe, ErrorMatches, "io: .*on closed pipe")
	c.Check(42, Equals, 42)
}
// 调用 login test 初始化数据
func (s *MySuite) SetUpSuite(c *C) {
	//增加配置文件
	config, err := goconfig.LoadConfigFile("./config/test.ini")
	common.IfError(err)
	// 配置数据库连接
	mysql, err := config.GetSection("mysql")
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",mysql["username"],mysql["password"],mysql["address"],mysql["port"],mysql["database"],mysql["charset"])
	db, err = gorm.Open("mysql", connStr)
	if err != nil{
		panic("连接数据库失败")
	}
	db.SingularTable(true)
	defer db.Close()

	requestUrl, err := config.GetValue("jcc","requestUrl")
	loginData = module.LoginTest(requestUrl,db)
}

func (s *MySuite) TearDownSuite(c *C) {
	fmt.Println("Tear Down Start ...")
}