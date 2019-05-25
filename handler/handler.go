package handler

import (
	"IM/pkg/errno"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type Response struct {
	Code            int         `json:"code"`
	Message         string      `json:"message"`
	Data            interface{} `json:"data,omitempty"`
	ValidationForms interface{} `json:"validationForms,omitempty"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	var (
		code             int
		message          string
		validationForms  interface{}
		errValidateForms []*errno.ErrValidateForm
	)
	// 判断error的类型
	switch typed := err.(type) {
	case validator.ValidationErrors:
		code, message, errValidateForms = errno.DecodeValidationErrors(typed)
	case *errno.Err:
		code, message = errno.DecodeErr(typed)
	case *errno.Errno:
		code, message = errno.DecodeErrNo(typed)
	default:
	}
	// 判断
	if errValidateForms == nil {
		validationForms = nil
	} else {
		validationForms = errValidateForms
	}
	// always return http.StatusOK
	c.JSON(http.StatusOK, Response{
		Code:            code,
		Message:         message,
		Data:            data,
		ValidationForms: validationForms,
	})
}

// 构建新的表单验证器
func NewValidator() *validator.Validate {
	return validator.New()
}
