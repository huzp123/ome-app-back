package errcode

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error 定义错误码结构
type Error struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Details []string `json:"details"`
}

// 常量错误码映射
var (
	Success                  = NewError(0, "成功")
	ServerError              = NewError(10000, "服务内部错误")
	InvalidParams            = NewError(10001, "无效参数")
	NotFound                 = NewError(10002, "找不到资源")
	UnauthorizedAuthNotExist = NewError(10003, "未授权认证失败")
	UnauthorizedTokenError   = NewError(10004, "未授权Token错误")
	UnauthorizedTokenTimeout = NewError(10005, "未授权Token超时")
	TooManyRequests          = NewError(10006, "请求过多")

	UserNotExist      = NewError(20001, "用户不存在")
	UserAlreadyExist  = NewError(20002, "用户已存在")
	UserPasswordError = NewError(20003, "用户密码错误")
	UserCreateFail    = NewError(20004, "创建用户失败")
	UserUpdateFail    = NewError(20005, "更新用户失败")
	UserDeleteFail    = NewError(20006, "删除用户失败")
	UserInvalidCode   = NewError(20007, "用户唯一编码无效")
)

// NewError 创建新的错误码
func NewError(code int, msg string) *Error {
	return &Error{
		Code:    code,
		Msg:     msg,
		Details: nil,
	}
}

// Error 实现error接口
func (e *Error) Error() string {
	return fmt.Sprintf("错误码: %d, 错误信息: %s", e.Code, e.Msg)
}

// WithDetails 添加详细错误信息
func (e *Error) WithDetails(details ...string) *Error {
	newError := *e
	for _, d := range details {
		newError.Details = append(newError.Details, d)
	}
	return &newError
}

// StatusCode 获取HTTP状态码
func (e *Error) StatusCode() int {
	switch e.Code {
	case Success.Code:
		return http.StatusOK
	case InvalidParams.Code:
		return http.StatusBadRequest
	case UnauthorizedAuthNotExist.Code,
		UnauthorizedTokenError.Code,
		UnauthorizedTokenTimeout.Code:
		return http.StatusUnauthorized
	case TooManyRequests.Code:
		return http.StatusTooManyRequests
	case NotFound.Code:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// Response 在gin上下文中输出错误响应
func (e *Error) Response(c *gin.Context) {
	response := gin.H{
		"code": e.Code,
		"msg":  e.Msg,
		"data": nil,
	}
	if len(e.Details) > 0 {
		response["details"] = e.Details
	}
	c.JSON(e.StatusCode(), response)
}
