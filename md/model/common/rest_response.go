package common

const (
	HttpSuccess     = 200 // 请求成功
	HttpAuthFailure = 401 // 认证失败
	HttpFailure     = 500 // 请求失败
)

// 通用返回json数据结构体
type RestResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 主动抛出异常结构体
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// 填写成功信息
func NewSuccess(message string) RestResponse {
	return RestResponse{
		Code:    HttpSuccess,
		Message: message,
	}
}

// 填写成功信息、数据
func NewSuccessData(message string, data interface{}) RestResponse {
	return RestResponse{
		Code:    HttpSuccess,
		Message: message,
		Data:    data,
	}
}

// 填写错误信息
func NewError(message string) ErrorResponse {
	return ErrorResponse{
		Code:    HttpFailure,
		Message: message,
	}
}

// 填写错误编号、错误信息
func NewErrorCode(code int, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// 填写错误信息、error对象
func NewErr(message string, err error) ErrorResponse {
	return ErrorResponse{
		Code:    HttpFailure,
		Message: message,
		Err:     err,
	}
}

// 填写错误编号、错误信息、error对象
func NewErrCode(code int, message string, err error) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
