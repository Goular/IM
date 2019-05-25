package errno

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
)

type Errno struct {
	Code    int
	Message string
}

func (err Errno) Error() string {
	return err.Message
}

// Err represents an error
type Err struct {
	Code    int
	Message string
	Err     error
}

// 表单校验
type ErrValidateForm struct {
	Key   string `json:"key"`   // 表单Field的名称
	Tag   string `json:"tag"`   // 检验类型的名称
	Param string `json:"param"` // 检验类型的设置参考值
}

func New(errno *Errno, err error) *Err {
	return &Err{Code: errno.Code, Message: errno.Message, Err: err}
}

func (err *Err) Add(message string) error {
	err.Message += " " + message
	return err
}

func (err *Err) Addf(format string, args ...interface{}) error {
	err.Message += " " + fmt.Sprintf(format, args...)
	return err
}

func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Err)
}

//func IsErrUserNotFound(err error) bool {
//	code, _ := DecodeErr(err)
//	return code == ErrUserNotFound.Code
//}

// 对Err的结构体数据进行解码
func DecodeErr(err *Err) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}
	return err.Code, err.Message
}

// 对ErrNo的结构体数据进行解码
func DecodeErrNo(err *Errno) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}
	return err.Code, err.Message
}

// 对ErrValidateForm的结构体数据进行解码
func DecodeValidationErrors(errs validator.ValidationErrors) (code int, message string, errValidateForms []*ErrValidateForm) {
	var (
		fieldError      validator.FieldError
		errValidateForm *ErrValidateForm
	)
	if errs == nil {
		return OK.Code, OK.Message, nil
	}
	for _, fieldError = range errs {
		errValidateForm = &ErrValidateForm{
			Key:   fieldError.Field(),
			Tag:   fieldError.Tag(),
			Param: fieldError.Param(),
		}
		errValidateForms = append(errValidateForms, errValidateForm)
	}
	return ErrUnregularParams.Code, ErrUnregularParams.Message, errValidateForms
}
