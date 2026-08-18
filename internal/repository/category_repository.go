package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

// categoryRepository 分类数据访问实现
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类数据访问
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// FindByID 根据 ID 查找分类
func (r *categoryRepository) FindByID(id uint) (*entity.Category, error) {
	var category entity.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// FindByName 根据名称查找分类
func (r *categoryRepository) FindByName(name string) (*entity.Category, error) {
	var category entity.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// FindBySlug 根据 slug 查找分类
func (r *categoryRepository) FindBySlug(slug string) (*entity.Category, error) {
	var category entity.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// Create 创建分类
func (r *categoryRepository) Create(category *entity.Category) error {
	return r.db.Create(category).Error
}

// Update 更新分类
func (r *categoryRepository) Update(category *entity.Category) error {
	return r.db.Save(category).Error
}

// Delete 删除分类（硬删除）
func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Category{}, id).Error
}

// List 分类列表
func (r *categoryRepository) List() ([]*entity.Category, error) {
	var categories []*entity.Category
	err := r.db.Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

// Count 统计分类数量
func (r *categoryRepository) Count() (int64, error) {
	var total int64
	err := r.db.Model(&entity.Category{}).Count(&total).Error
	return total, err
}

// GetCategoryArticleCount 获取分类文章数量
func (r *categoryRepository) GetCategoryArticleCount(categoryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Article{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}
