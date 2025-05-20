package ierr

import (
	"fmt"

	"github.com/hzxiao/goutil"
)

type IError struct {
	Code     Code       `json:"code"`
	Msg      string     `json:"msg"` //跟code一致的
	MoreInfo string     `json:"more_info"`
	Data     goutil.Map `json:"data"`
}

func NewIError(code Code, moreInfo string) *IError {
	// var finalCode = CodeNum*1000 + code
	var finalCode = code
	return &IError{Code: finalCode, MoreInfo: moreInfo}
}

func NewErrorData(code Code, moreInfo string, data goutil.Map) *IError {
	// var finalCode = CodeNum*1000 + code
	var finalCode = code
	return &IError{Code: finalCode, MoreInfo: moreInfo, Data: data}
}

func NewIErrorf(code Code, format string, args ...interface{}) *IError {
	// var finalCode = CodeNum*1000 + code
	var finalCode = code
	return &IError{Code: finalCode, MoreInfo: fmt.Sprintf(format, args...)}
}

func (e *IError) Error() string {
	return fmt.Sprintf("code(%d),msg(%s),moreInfo(%s)", e.Code, e.Msg, e.MoreInfo)
}
