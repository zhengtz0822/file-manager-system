package handler

import (
	"file-manager-service/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ApplicationHandler struct {
	appService *service.ApplicationService
}

func NewApplicationHandler() *ApplicationHandler {
	return &ApplicationHandler{
		appService: service.NewApplicationService(),
	}
}

// CreateApplication 创建应用
// @Summary 创建应用
// @Description 创建新应用，系统自动生成应用账号和密钥
// @Tags Application
// @Accept json
// @Produce json
// @Param request body service.CreateApplicationRequest true "应用信息"
// @Success 200 {object} Response{data=service.CreateApplicationResponse}
// @Router /applications [post]
func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	var req service.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	app, err := h.appService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "创建应用失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "创建成功",
		Data:    app,
	})
}

// ListApplications 获取应用列表
// @Summary 获取应用列表
// @Description 分页获取应用列表
// @Tags Application
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} Response{data=service.ListApplicationsResponse}
// @Router /applications [get]
func (h *ApplicationHandler) ListApplications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	result, err := h.appService.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取应用列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取成功",
		Data:    result,
	})
}

// GetApplication 获取应用详情
// @Summary 获取应用详情
// @Description 根据ID获取应用详情
// @Tags Application
// @Accept json
// @Produce json
// @Param id path int true "应用ID"
// @Success 200 {object} Response{data=service.CreateApplicationResponse}
// @Router /applications/{id} [get]
func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的应用ID",
		})
		return
	}

	app, err := h.appService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "应用不存在",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取成功",
		Data:    app,
	})
}

// UpdateApplicationStatus 更新应用状态
// @Summary 更新应用状态
// @Description 启用或禁用应用
// @Tags Application
// @Accept json
// @Produce json
// @Param id path int true "应用ID"
// @Param request body map[string]int true "状态"
// @Success 200 {object} Response
// @Router /applications/{id}/status [put]
func (h *ApplicationHandler) UpdateApplicationStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的应用ID",
		})
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	if err := h.appService.UpdateStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "更新状态失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "更新成功",
	})
}

// DeleteApplication 删除应用
// @Summary 删除应用
// @Description 根据ID删除应用
// @Tags Application
// @Accept json
// @Produce json
// @Param id path int true "应用ID"
// @Success 200 {object} Response
// @Router /applications/{id} [delete]
func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的应用ID",
		})
		return
	}

	if err := h.appService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "删除失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "删除成功",
	})
}

// GetApplicationOptions 获取应用选项（不含敏感信息）
// @Summary 获取应用选项
// @Description 获取应用选项列表，用于下拉选择等场景，不包含密钥等敏感信息
// @Tags Application
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=[]service.ApplicationOption}
// @Router /applications/options [get]
func (h *ApplicationHandler) GetApplicationOptions(c *gin.Context) {
	options, err := h.appService.GetOptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取应用选项失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取成功",
		Data:    options,
	})
}
