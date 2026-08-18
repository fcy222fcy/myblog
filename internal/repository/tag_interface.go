package repository

import "blog/internal/model/entity"

// TagRepository 标签数据访问接口
type TagRepository interface {
	// FindByID 根据 ID 查找标签
	FindByID(id uint) (*entity.Tag, error)

	// FindByName 根据名称查找标签
	FindByName(name string) (*entity.Tag, error)

	// FindBySlug 根据 slug 查找标签
	FindBySlug(slug string) (*entity.Tag, error)

	// Create 创建标签
	Create(tag *entity.Tag) error

	// Update 更新标签
	Update(tag *entity.Tag) error

	// Delete 删除标签（硬删除，清理关联）
	Delete(id uint) error

	// List 标签列表
	List() ([]*entity.Tag, error)

	// Count 统计标签数量
	Count() (int64, error)

	// GetTagArticleCount 获取标签文章数量
	GetTagArticleCount(tagID uint) (int64, error)

	// FindByIDs 根据ID列表查找标签
	FindByIDs(ids []uint) ([]*entity.Tag, error)
}
