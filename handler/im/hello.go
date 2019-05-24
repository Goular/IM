package im

import (
	"IM/handler"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 登录表单
type LoginForm struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

// IM系统用户登录
func Login(c *gin.Context) {
	var (
		loginForm LoginForm
		err       error
	)
	if err = c.ShouldBind(&loginForm); err != nil {
		fmt.Println(err)
		return
	}
	handler.SendResponse(c, nil, loginForm)
}
