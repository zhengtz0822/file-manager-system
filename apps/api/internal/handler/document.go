package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"file-manager-service/internal/pkg/jwt"
	"file-manager-service/internal/pkg/storage"
	"file-manager-service/internal/repository"
	"file-manager-service/internal/service"

	"github.com/gin-gonic/gin"
)

// DocumentHandler 文档处理器
type DocumentHandler struct {
	docService *service.DocumentService
}

// NewDocumentHandler 创建文档处理器
func NewDocumentHandler() *DocumentHandler {
	return &DocumentHandler{
		docService: service.NewDocumentService(),
	}
}

// InitUpload 初始化分片上传
// @Summary 初始化分片上传
// @Description 初始化一个大文件分片上传
// @Tags 文档上传
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.InitUploadRequest true "上传信息"
// @Success 200 {object} Response{data=service.InitUploadResponse}
// @Router /documents/chunks/init [post]
func (h *DocumentHandler) InitUpload(c *gin.Context) {
	var req service.InitUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.docService.InitUpload(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "初始化上传失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "初始化成功",
		Data:    resp,
	})
}

// UploadChunk 上传分片
// @Summary 上传分片
// @Description 上传单个分片
// @Tags 文档上传
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param upload_id formData string true "上传ID"
// @Param chunk_number formData int true "分片序号"
// @Param file formData file true "分片文件"
// @Success 200 {object} Response
// @Router /documents/chunks/upload [post]
func (h *DocumentHandler) UploadChunk(c *gin.Context) {
	uploadID := c.PostForm("upload_id")
	var chunkNumber int
	if err := c.ShouldBind(&chunkNumber); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "文件上传失败",
			Error:   err.Error(),
		})
		return
	}

	// 保存分片
	if err := storage.SaveChunk(uploadID, chunkNumber, file); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "保存分片失败",
			Error:   err.Error(),
		})
		return
	}

	// 更新上传进度
	if err := h.docService.UploadChunk(uploadID, chunkNumber, file.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "更新上传进度失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "分片上传成功",
	})
}

// CompleteUpload 完成上传
// @Summary 完成上传
// @Description 合并所有分片并完成上传
// @Tags 文档上传
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.CompleteUploadRequest true "完成上传信息"
// @Success 200 {object} Response{data=service.CompleteUploadResponse}
// @Router /documents/chunks/complete [post]
func (h *DocumentHandler) CompleteUpload(c *gin.Context) {
	// 获取账号信息（应用账号或用户名）
	tokenType, exists := c.Get("token_type")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: "未授权",
		})
		return
	}

	var account string
	var userID *int
	var appID *int

	if tokenType == "app" {
		// 应用 Token - 使用应用标识
		appIdentifier, exists := c.Get("app_identifier")
		if !exists {
			c.JSON(http.StatusUnauthorized, Response{
				Code:    401,
				Message: "未授权",
			})
			return
		}
		account = appIdentifier.(string)

		// 获取 app_id
		if appIDValue, exists := c.Get("app_id"); exists {
			val := appIDValue.(int)
			appID = &val
		}
	} else {
		// 用户 Token
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, Response{
				Code:    401,
				Message: "未授权",
			})
			return
		}
		account = username.(string)

		// 获取 user_id
		if userIDValue, exists := c.Get("user_id"); exists {
			val := userIDValue.(int)
			userID = &val
		}
	}

	var req service.CompleteUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	// 设置上传人信息
	req.Account = account
	req.TokenType = tokenType.(string)
	req.UserID = userID
	req.AppID = appID

	resp, err := h.docService.CompleteUpload(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "完成上传失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "上传完成",
		Data:    resp,
	})
}

// CancelUpload 取消上传
// @Summary 取消上传
// @Description 取消分片上传并清理临时文件
// @Tags 文档上传
// @Accept json
// @Produce json
// @Security Bearer
// @Param upload_id path string true "上传ID"
// @Success 200 {object} Response
// @Router /documents/chunks/{upload_id} [delete]
func (h *DocumentHandler) CancelUpload(c *gin.Context) {
	uploadID := c.Param("upload_id")

	if err := h.docService.CancelUpload(uploadID); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "取消上传失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "上传已取消",
	})
}

// UploadSingleFile 单文件上传
// @Summary 单文件上传
// @Description 直接上传单个文件（适用于小文件）
// @Tags 文档上传
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param file formData file true "文件"
// @Success 200 {object} Response{data=service.UploadSingleFileResponse}
// @Failure 400 {object} Response "参数错误或文件类型不支持"
// @Failure 413 {object} Response "文件大小超过限制"
// @Router /documents/upload [post]
func (h *DocumentHandler) UploadSingleFile(c *gin.Context) {
	// 获取账号信息（应用账号或用户名）
	tokenType, exists := c.Get("token_type")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: "未授权",
		})
		return
	}

	var account string
	var userID *int
	var appID *int

	if tokenType == "app" {
		// 应用 Token - 使用应用标识
		appIdentifier, exists := c.Get("app_identifier")
		if !exists {
			c.JSON(http.StatusUnauthorized, Response{
				Code:    401,
				Message: "未授权",
			})
			return
		}
		account = appIdentifier.(string)

		// 获取 app_id
		if appIDValue, exists := c.Get("app_id"); exists {
			val := appIDValue.(int)
			appID = &val
		}
	} else {
		// 用户 Token
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, Response{
				Code:    401,
				Message: "未授权",
			})
			return
		}
		account = username.(string)

		// 获取 user_id
		if userIDValue, exists := c.Get("user_id"); exists {
			val := userIDValue.(int)
			userID = &val
		}
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "参数错误: 未找到文件",
			Error:   err.Error(),
		})
		return
	}

	// 打开文件
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "文件读取失败",
			Error:   err.Error(),
		})
		return
	}
	defer fileContent.Close()

	// 调用服务层上传文件
	resp, err := h.docService.UploadSingleFile(account, file.Filename, file.Size, fileContent, tokenType.(string), userID, appID)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "上传失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "上传成功",
		Data:    resp,
	})
}

// GetDocument 获取文档信息
// @Summary 获取文档信息
// @Description 根据ID获取文档详细信息
// @Tags 文档管理
// @Produce json
// @Security Bearer
// @Param id path string true "文档ID"
// @Success 200 {object} Response{data=model.Document}
// @Router /documents/{id} [get]
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	id := c.Param("id")

	doc, err := h.docService.GetDocument(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "文档不存在",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取成功",
		Data:    doc,
	})
}

// ListDocuments 获取文档列表
// @Summary 获取文档列表
// @Description 分页获取文档列表，支持关键词搜索
// @Tags 文档管理
// @Produce json
// @Security Bearer
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} Response{data=service.ListDocumentsResponse}
// @Router /documents [get]
func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	var req service.ListDocumentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.docService.List(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取列表失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取成功",
		Data:    resp,
	})
}

// DeleteDocument 删除文档
// @Summary 删除文档
// @Description 删除指定文档
// @Tags 文档管理
// @Produce json
// @Security Bearer
// @Param id path string true "文档ID"
// @Success 200 {object} Response
// @Router /documents/{id} [delete]
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	id := c.Param("id")

	if err := h.docService.DeleteDocument(id); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "删除文档失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "删除成功",
	})
}

// DownloadDocument 下载文档
// @Summary 下载文档
// @Description 下载指定文档
// @Tags 文档访问
// @Produce application/octet-stream
// @Security Bearer
// @Param id path string true "文档ID"
// @Success 200 {file} file
// @Router /documents/{id}/download [get]
func (h *DocumentHandler) DownloadDocument(c *gin.Context) {
	id := c.Param("id")

	doc, err := h.docService.GetDocument(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "文档不存在",
		})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(doc.StoragePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "文件不存在",
		})
		return
	}

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+doc.FileName)
	c.Header("Content-Type", "application/octet-stream")

	// 发送文件
	c.File(doc.StoragePath)
}

// PreviewDocument 预览文档
// @Summary 预览文档
// @Description 在线预览文档（支持图片和PDF）
// @Tags 文档访问
// @Produce image/jpeg,application/pdf
// @Security Bearer
// @Param id path string true "文档ID"
// @Success 200 {file} file
// @Router /documents/{id}/preview [get]
func (h *DocumentHandler) PreviewDocument(c *gin.Context) {
	id := c.Param("id")

	doc, err := h.docService.GetDocument(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "文档不存在",
		})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(doc.StoragePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "文件不存在",
		})
		return
	}

	// 检查文件类型是否支持预览
	ext := filepath.Ext(doc.StoragePath)
	supportedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".pdf":  true,
	}

	if !supportedExtensions[ext] {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "不支持的预览格式",
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", doc.FileType)
	c.File(doc.StoragePath)
}

// GetCurrentUser 获取当前用户
func (h *DocumentHandler) GetCurrentUser(c *gin.Context) {
	claimsInterface, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: "未授权",
		})
		return
	}

	claims := claimsInterface.(*jwt.Claims)
	userRepo := repository.NewUserRepository()
	user, err := userRepo.FindByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "用户不存在",
			Error:   err.Error(),
		})
		return
	}

	// 构造前端需要的用户数据格式
	userData := map[string]interface{}{
		"id":         user.ID,
		"userid":     fmt.Sprintf("%d", user.ID), // 前端需要的 userid
		"name":       user.Username,               // 前端需要的 name 字段
		"username":   user.Username,
		"created_at": user.CreatedAt,
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取成功",
		Data:    userData,
	})
}
