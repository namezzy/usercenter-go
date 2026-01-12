package common

// BaseResponse 通用返回类
type BaseResponse struct {
	Code        int         `json:"code"`
	Data        interface{} `json:"data"`
	Message     string      `json:"message"`
	Description string      `json:"description"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) *BaseResponse {
	return &BaseResponse{
		Code:        0,
		Data:        data,
		Message:     "ok",
		Description: "",
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(errorCode ErrorCode) *BaseResponse {
	return &BaseResponse{
		Code:        errorCode.Code,
		Data:        nil,
		Message:     errorCode.Message,
		Description: errorCode.Description,
	}
}

// NewErrorResponseWithDesc 创建带描述的错误响应
func NewErrorResponseWithDesc(errorCode ErrorCode, description string) *BaseResponse {
	return &BaseResponse{
		Code:        errorCode.Code,
		Data:        nil,
		Message:     errorCode.Message,
		Description: description,
	}
}
