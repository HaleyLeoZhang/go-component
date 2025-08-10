package xgin

import (
	"context"
	"net/http"

	"github.com/HaleyLeoZhang/go-component/driver/xlog"
)

// ---------------------------------------------------------------------
// 		业务错误 Error 模型
// ---------------------------------------------------------------------
type BusinessError struct {
	Code    int
	Message string
	error   *error
}

func (b *BusinessError) Error() string {
	return b.Message
}

// ---------------------------------------------------------------------
// 		HTTP 响应模型
// ---------------------------------------------------------------------

type ResponseModel struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

// HTTP 响应模型
func (o *Gin) Response(ctx context.Context, err error, data interface{}) {
	code := HTTP_RESPONSE_CODE_SUCCESS
	message := ""
	if err != nil {
		switch err.(type) {
		case *BusinessError:
			businessError := err.(*BusinessError)
			code = businessError.Code
			message = businessError.Message
			data = nil
			xlog.Infof(ctx, "Response BusinessError(%+v)", err)
		default:
			code = HTTP_RESPONSE_CODE_UNKNOWN_FAIL
			message = "服务繁忙"
			data = nil
			xlog.Errorf(ctx, "Response Error(%+v)", err)
		}
	}
	o.GinContext.JSON(http.StatusOK, ResponseModel{
		Code: code,
		Msg:  message,
		Data: data,
	})
	return
}

func NewBusinessError(msg string, code int) (err *BusinessError) {
	err = &BusinessError{Message: msg, Code: code}
	return
}
