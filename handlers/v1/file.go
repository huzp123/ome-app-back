package v1

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"ome-app-back/services"
)

// FileAPI 处理文件相关接口
type FileAPI struct {
	fileService *services.FileService
}

// NewFileAPI 创建文件API处理实例
func NewFileAPI(fileService *services.FileService) *FileAPI {
	return &FileAPI{
		fileService: fileService,
	}
}

// GetFile 获取文件内容
func (a *FileAPI) GetFile(c *gin.Context) {
	// 获取文件路径参数
	filePath := c.Param("filepath")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文件路径不能为空",
			"data": nil,
		})
		return
	}

	// 安全性检查：防止任意文件访问
	if strings.Contains(filePath, "..") || !strings.HasPrefix(filePath, "uploads/") {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "禁止访问该路径",
			"data": nil,
		})
		return
	}

	// 读取文件
	data, mimeType, err := a.fileService.GetFile(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "文件不存在或无法读取: " + err.Error(),
			"data": nil,
		})
		return
	}

	// 设置内容类型并返回文件
	c.Header("Content-Type", mimeType)
	c.Header("Content-Disposition", "inline; filename="+filepath.Base(filePath))
	c.Data(http.StatusOK, mimeType, data)
}

// GetUserFile 获取用户文件（需要验证权限）
func (a *FileAPI) GetUserFile(c *gin.Context) {
	// 记录请求路径
	requestPath := c.Request.URL.Path
	fmt.Printf("[文件访问] 请求路径: %s\n", requestPath)

	// 获取用户ID和文件路径
	userID := getUserID(c)
	fmt.Printf("[文件访问] 获取到的用户ID: %d\n", userID)

	if userID == 0 {
		fmt.Printf("[文件访问] 未授权访问: %s\n", requestPath)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
			"data": nil,
		})
		return
	}

	filePath := c.Param("filepath")
	fmt.Printf("[文件访问] 解析到的文件路径参数: '%s'\n", filePath)

	if filePath == "" {
		fmt.Printf("[文件访问] 文件路径为空\n")
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文件路径不能为空",
			"data": nil,
		})
		return
	}

	// 处理路径中的前导斜杠
	if strings.HasPrefix(filePath, "/") {
		filePath = filePath[1:]
		fmt.Printf("[文件访问] 处理后的文件路径: '%s'\n", filePath)
	}

	// 权限检查：仅允许访问用户自己的文件
	expectedPrefix := fmt.Sprintf("uploads/user_%d", userID)
	fmt.Printf("[文件访问] 期望的路径前缀: '%s', 实际路径: '%s'\n", expectedPrefix, filePath)

	if !strings.HasPrefix(filePath, expectedPrefix) {
		fmt.Printf("[文件访问] 权限检查失败: 路径 '%s' 不符合前缀 '%s'\n", filePath, expectedPrefix)
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "无权访问该文件",
			"data": nil,
		})
		return
	}

	fmt.Printf("[文件访问] 权限检查通过，准备读取文件: %s\n", filePath)

	// 读取文件
	data, mimeType, err := a.fileService.GetFile(filePath)
	if err != nil {
		fmt.Printf("[文件访问] 读取文件失败: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "文件不存在或无法读取: " + err.Error(),
			"data": nil,
		})
		return
	}

	fmt.Printf("[文件访问] 文件读取成功, 大小: %d字节, 类型: %s\n", len(data), mimeType)

	// 设置内容类型并返回文件
	c.Header("Content-Type", mimeType)
	c.Header("Content-Disposition", "inline; filename="+filepath.Base(filePath))
	c.Data(http.StatusOK, mimeType, data)
}
