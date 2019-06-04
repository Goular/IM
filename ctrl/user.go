package ctrl

import (
	"../model"
	"../service"
	"../util"
	"fmt"
	"math/rand"
	"net/http"
)

func UserLogin(writer http.ResponseWriter,
	request *http.Request) {
	//数据库操作
	//逻辑处理
	//restapi json/xml返回
	//1.获取前端传递的参数
	//mobile,passwd
	//解析参数
	//如何获得参数
	//解析参数
	err := request.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	mobile := request.PostForm.Get("mobile")
	passwd := request.PostForm.Get("passwd")

	// 模拟
	user, err := userService.Login(mobile, passwd)
	if err != nil {
		util.RespFail(writer, err.Error())
	} else {
		util.RespOk(writer, user, "")
	}

	//fmt.Println(mobile, passwd)
	//
	//loginok := false
	//if mobile == "18600000000" && passwd == "123456" {
	//	loginok = true
	//}
	//if loginok {
	//	//{"id":1,"token":"xx"}
	//	data := make(map[string]interface{})
	//	data["id"] = 1
	//	data["token"] = "test"
	//	util.RespOk(writer, data, "")
	//} else {
	//	util.RespFail(writer, "密码不正确")
	//
	//}
	//go get github.com/go-xorm/xorm
	//go get github.com/go-sql-driver/mysql
	//返回json ok

	//如何返回JSON

	//io.WriteString(writer,request.PostForm.Get("mobile"))
}

var userService service.UserService

func UserRegister(writer http.ResponseWriter,
	request *http.Request) {

	request.ParseForm()
	//
	mobile := request.PostForm.Get("mobile")
	//
	plainpwd := request.PostForm.Get("passwd")
	nickname := fmt.Sprintf("user%06d", rand.Int31())
	avatar := ""
	sex := model.SEX_UNKNOW

	user, err := userService.Register(mobile, plainpwd, nickname, avatar, sex)
	if err != nil {
		util.RespFail(writer, err.Error())
	} else {
		util.RespOk(writer, user, "")

	}

}
