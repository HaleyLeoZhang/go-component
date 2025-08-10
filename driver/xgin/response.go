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
	Code  int         `json:"code"`
	Msg   string      `json:"message"`
	Data  interface{} `json:"data"`
	ReqID string      `json:"req_id"` // 请求ID
}

// HTTP 响应模型
func (o *Gin) Response(ctx context.Context, err error, data interface{}) {
	code := HTTP_RESPONSE_CODE_SUCCESS
	message := ""
	reqID := xlog.GetLogID(ctx)
	// 没报错时
	if err == nil {
		o.GinContext.JSON(http.StatusOK, ResponseModel{
			Code:  code,
			Msg:   message,
			Data:  data,
			ReqID: reqID,
		})
		return
	}
	// 有报错时
	switch err.(type) {
	case *BusinessError: // 可预期错误
		businessError := err.(*BusinessError)
		code = businessError.Code
		message = businessError.Message
		data = nil
		xlog.Warnf(ctx, "Response BusinessError(%v)", err)
		//o.GinContext.JSON(http.StatusOK, ResponseModel{
		//	Code:  code,
		//	Msg:   message,
		//	Data:  data,
		//	ReqID: reqID,
		//})
	default: // 不可预期错误
		code = HTTP_RESPONSE_CODE_UNKNOWN_FAIL
		message = "服务繁忙"
		data = nil
		xlog.Errorf(ctx, "Response Error(%+v)", err)
		o.GinContext.JSON(http.StatusInternalServerError, ResponseModel{
			Code:  code,
			Msg:   message,
			Data:  data,
			ReqID: reqID,
		})
	}
	return
}

func NewBusinessError(msg string, code ...int) (err *BusinessError) {
	var realCode = HTTP_RESPONSE_CODE_UNKNOWN_FAIL
	if len(code) > 0 {
		realCode = code[0]
	}
	err = &BusinessError{Message: msg, Code: realCode}
	return
}
