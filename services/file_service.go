package services

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"ome-app-back/config"
)

// FileService 处理文件上传相关服务
type FileService struct {
	uploadDir string
	maxSize   int64
}

// NewFileService 创建文件服务实例
func NewFileService(cfg *config.UploadConfig) *FileService {
	// 使用配置中的上传目录
	uploadDir := cfg.Dir

	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		fmt.Printf("警告：无法创建上传目录 %s: %v\n", uploadDir, err)
	}

	return &FileService{
		uploadDir: uploadDir,
		maxSize:   cfg.MaxSize,
	}
}

// UploadImage 上传图片文件
func (s *FileService) UploadImage(file *multipart.FileHeader, userID int64) (string, error) {
	// 检查文件大小
	if file.Size > s.maxSize {
		return "", fmt.Errorf("文件过大：%d 字节，最大允许 %d 字节", file.Size, s.maxSize)
	}

	// 创建用户目录
	userDir := filepath.Join(s.uploadDir, fmt.Sprintf("user_%d", userID))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return "", fmt.Errorf("创建用户目录失败: %v", err)
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), randomString(8), ext)
	fullPath := filepath.Join(userDir, filename)

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开上传文件失败: %v", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("复制文件内容失败: %v", err)
	}

	// 返回相对路径
	relativePath := filepath.Join("uploads", fmt.Sprintf("user_%d", userID), filename)
	return relativePath, nil
}

// GetImageBase64 读取图片并转换为Base64编码
func (s *FileService) GetImageBase64(filePath string) (string, error) {
	// 获取完整路径
	fullPath := filepath.Join(".", filePath)

	// 读取文件
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("读取图片文件失败: %v", err)
	}

	// 转为base64编码
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded, nil
}

// randomString 生成指定长度的随机字符串
func randomString(length int) string {
	// 使用crypto/rand生成安全的随机字符串
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLen := len(charset)

	// 创建一个字节切片用于存储随机字节
	randomBytes := make([]byte, length)

	// 生成随机字节
	if _, err := io.ReadFull(cryptorand.Reader, randomBytes); err != nil {
		// 如果生成失败，记录错误并回退到不太安全但足够可靠的方法
		fmt.Printf("警告：无法生成加密安全的随机字符串: %v\n", err)
		return fallbackRandomString(length)
	}

	// 将随机字节映射到字符集
	result := make([]byte, length)
	for i, b := range randomBytes {
		result[i] = charset[int(b)%charsetLen]
	}

	return string(result)
}

// fallbackRandomString 备用的随机字符串生成方法
func fallbackRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLen := len(charset)

	// 使用time.Now().UnixNano()作为种子
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(charsetLen)]
	}

	return string(result)
}

// GetFile 获取文件内容和MIME类型
func (s *FileService) GetFile(filePath string) ([]byte, string, error) {
	// 获取完整路径
	fullPath := filepath.Join(".", filePath)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("文件不存在: %s", filePath)
	}

	// 读取文件
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, "", fmt.Errorf("读取文件失败: %v", err)
	}

	// 确定MIME类型
	ext := filepath.Ext(fullPath)
	mimeType := getMimeType(ext)

	return data, mimeType, nil
}

// getMimeType 根据文件扩展名确定MIME类型
func getMimeType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}
