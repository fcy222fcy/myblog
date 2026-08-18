package repository

import "blog/internal/model/entity"

// CategoryRepository 分类数据访问接口
type CategoryRepository interface {
	// FindByID 根据 ID 查找分类
	FindByID(id uint) (*entity.Category, error)

	// FindByName 根据名称查找分类
	FindByName(name string) (*entity.Category, error)

	// FindBySlug 根据 slug 查找分类
	FindBySlug(slug string) (*entity.Category, error)

	// Create 创建分类
	Create(category *entity.Category) error

	// Update 更新分类
	Update(category *entity.Category) error

	// Delete 删除分类（硬删除）
	Delete(id uint) error

	// List 分类列表
	List() ([]*entity.Category, error)

	// Count 统计分类数量
	Count() (int64, error)

	// GetCategoryArticleCount 获取分类文章数量
	GetCategoryArticleCount(categoryID uint) (int64, error)
}
