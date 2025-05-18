package ierr

import (
	"fmt"
	"wzkj-common/pkg/util"
)

type IError struct {
	Code     Code   `json:"code"`
	Msg      string `json:"msg"` //跟code一致的
	MoreInfo string `json:"more_info"`
	Data     util.M `json:"data"`
}

func NewIError(code Code, moreInfo string) *IError {
	// var finalCode = CodeNum*1000 + code
	var finalCode = code
	return &IError{Code: finalCode, MoreInfo: moreInfo}
}

func NewErrorData(code Code, moreInfo string, data util.M) *IError {
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
