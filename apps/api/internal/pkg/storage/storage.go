package storage

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"file-manager-service/internal/config"
)

var cfg *config.StorageConfig

// Init 初始化存储
func Init(c *config.StorageConfig) {
	cfg = c
	// 创建必要的目录
	os.MkdirAll(cfg.ChunkPath, 0755)
	os.MkdirAll(cfg.DocumentPath, 0755)
}

// GetMimeType 获取文件 MIME 类型
func GetMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	ext = strings.TrimPrefix(ext, ".")

	mimeMap := map[string]string{
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"gif":  "image/gif",
		"pdf":  "application/pdf",
		"doc":  "application/msword",
		"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"xls":  "application/vnd.ms-excel",
		"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"ppt":  "application/vnd.ms-powerpoint",
		"pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"txt":  "text/plain",
		"md":   "text/markdown",
	}

	if mime, ok := mimeMap[ext]; ok {
		return mime
	}
	return "application/octet-stream"
}

// IsAllowedExtension 检查文件扩展名是否允许
func IsAllowedExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	ext = strings.TrimPrefix(ext, ".")

	for _, allowed := range cfg.AllowedExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}

// SaveChunk 保存分片文件
func SaveChunk(uploadID string, chunkIndex int, file *multipart.FileHeader) error {
	filename := fmt.Sprintf("%s_%d", uploadID, chunkIndex)
	dstPath := filepath.Join(cfg.ChunkPath, filename)

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// MergeChunks 合并分片文件
func MergeChunks(uploadID string, totalChunks int, targetPath string) error {
	// 创建目标文件
	dst, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// 合并所有分片
	for i := 1; i <= totalChunks; i++ {
		chunkPath := filepath.Join(cfg.ChunkPath, fmt.Sprintf("%s_%d", uploadID, i))
		src, err := os.Open(chunkPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(dst, src)
		src.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// CleanChunks 清理分片文件
func CleanChunks(uploadID string, totalChunks int) error {
	for i := 1; i <= totalChunks; i++ {
		chunkPath := filepath.Join(cfg.ChunkPath, fmt.Sprintf("%s_%d", uploadID, i))
		if err := os.Remove(chunkPath); err != nil {
			return err
		}
	}
	return nil
}

// GetFileMD5 计算文件 MD5
func GetFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GenerateStoragePath 生成存储路径
// account: 应用账号或用户名
// fileID: 文档ID
// filename: 文件名
func GenerateStoragePath(account string, fileID string, filename string) string {
	// 生成年月路径: yyyy-MM
	yearMonth := time.Now().Format("2006-01")
	ext := filepath.Ext(filename)
	// 路径格式: uploads/应用账号/yyyy-MM/文档id
	return filepath.Join(cfg.DocumentPath, account, yearMonth, fileID+ext)
}

// SaveSingleFile 保存单个文件（适用于小文件直接上传）
func SaveSingleFile(account string, fileID string, filename string, file *multipart.FileHeader) error {
	// 生成存储路径
	storagePath := GenerateStoragePath(account, fileID, filename)

	// 创建目录
	yearMonth := time.Now().Format("2006-01")
	dirPath := filepath.Join(cfg.DocumentPath, account, yearMonth)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(storagePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// 复制文件内容
	_, err = io.Copy(dst, src)
	return err
}
