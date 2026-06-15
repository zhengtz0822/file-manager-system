package repository

import (
	"file-manager-service/internal/model"
)

// DocumentRepository 文档数据访问层
type DocumentRepository struct{}

// NewDocumentRepository 创建文档仓库
func NewDocumentRepository() *DocumentRepository {
	return &DocumentRepository{}
}

// Create 创建文档
func (r *DocumentRepository) Create(doc *model.Document) error {
	return model.DB.Create(doc).Error
}

// FindByID 根据 ID 查找文档
func (r *DocumentRepository) FindByID(id string) (*model.Document, error) {
	var doc model.Document
	err := model.DB.Where("id = ? AND status = ?", id, model.DocumentStatusNormal).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// FindByUploadID 根据 UploadID 查找文档
func (r *DocumentRepository) FindByUploadID(uploadID string) (*model.Document, error) {
	var doc model.Document
	err := model.DB.Where("upload_id = ? AND status = ?", uploadID, model.DocumentStatusNormal).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// Delete 删除文档（软删除）
func (r *DocumentRepository) Delete(id string) error {
	return model.DB.Model(&model.Document{}).Where("id = ?", id).Update("status", model.DocumentStatusDeleted).Error
}

// List 分页查询文档列表
func (r *DocumentRepository) List(page, pageSize int, keyword string, appIdentifier string) ([]*model.Document, int64, error) {
	var docs []*model.Document
	var total int64

	query := model.DB.Model(&model.Document{}).Where("status = ?", model.DocumentStatusNormal)

	if keyword != "" {
		query = query.Where("file_name LIKE ?", "%"+keyword+"%")
	}

	// 按应用标识筛选
	if appIdentifier != "" {
		// JOIN applications 表筛选应用标识
		query = query.Joins("INNER JOIN applications ON documents.app_id = applications.id").
			Where("applications.app_identifier = ? AND documents.uploaded_by = ?", appIdentifier, model.UploadedByApp)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&docs).Error
	if err != nil {
		return nil, 0, err
	}

	return docs, total, nil
}
