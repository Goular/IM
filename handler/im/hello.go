package im

import (
	"IM/handler"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

// 登录表单
type LoginForm struct {
	Username string `form:"username" validate:"required,min=4" json:"username"`
	Password string `form:"password" validate:"required,min=6" json:"password"`
}

// IM系统用户登录
func Login(c *gin.Context) {
	// 变量定义
	var (
		loginForm LoginForm
		err       error
		validate  *validator.Validate
	)
	// 绑定模型
	if err = c.ShouldBind(&loginForm); err != nil {
		fmt.Println(err)
		handler.SendResponse(c, err, nil)
		return
	}
	// 校验表单
	validate = handler.NewValidator()
	if err = validate.Struct(loginForm); err != nil {
		handler.SendResponse(c, err, nil)
		return
	}
	// 业务逻辑处理
	handler.SendResponse(c, nil, loginForm)
}
