package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response API响应结构
type Response struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	Details []string    `json:"details,omitempty"`
}

// ResponseSuccess 返回成功响应
func ResponseSuccess(c *gin.Context, data interface{}) {
	resp := Response{
		Code: 0,
		Msg:  "成功",
		Data: data,
	}
	c.JSON(http.StatusOK, resp)
}

// ResponseError 返回错误响应
func ResponseError(c *gin.Context, httpCode int, msg string, details ...string) {
	resp := Response{
		Code:    httpCode,
		Msg:     msg,
		Data:    nil,
		Details: details,
	}
	c.JSON(httpCode, resp)
}
