package common

// ErrorCode 错误码
type ErrorCode struct {
	Code        int
	Message     string
	Description string
}

var (
	SUCCESS      = ErrorCode{Code: 0, Message: "ok", Description: ""}
	PARAMS_ERROR = ErrorCode{Code: 40000, Message: "请求参数错误", Description: ""}
	NULL_ERROR   = ErrorCode{Code: 40001, Message: "请求数据为空", Description: ""}
	NOT_LOGIN    = ErrorCode{Code: 40100, Message: "未登录", Description: ""}
	NO_AUTH      = ErrorCode{Code: 40101, Message: "无权限", Description: ""}
	SYSTEM_ERROR = ErrorCode{Code: 50000, Message: "系统内部异常", Description: ""}
)
