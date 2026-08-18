package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

// tagRepository 标签数据访问实现
type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建标签数据访问
func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

// FindByID 根据 ID 查找标签
func (r *tagRepository) FindByID(id uint) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

// FindByName 根据名称查找标签
func (r *tagRepository) FindByName(name string) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

// FindBySlug 根据 slug 查找标签
func (r *tagRepository) FindBySlug(slug string) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

// Create 创建标签
func (r *tagRepository) Create(tag *entity.Tag) error {
	return r.db.Create(tag).Error
}

// Update 更新标签
func (r *tagRepository) Update(tag *entity.Tag) error {
	return r.db.Save(tag).Error
}

// Delete 删除标签（硬删除，清理 article_tags 关联行）
func (r *tagRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 清理文章-标签关联
		if err := tx.Where("tag_id = ?", id).Delete(&entity.ArticleTag{}).Error; err != nil {
			return err
		}
		return tx.Delete(&entity.Tag{}, id).Error
	})
}

// List 标签列表
func (r *tagRepository) List() ([]*entity.Tag, error) {
	var tags []*entity.Tag
	err := r.db.Order("created_at DESC").Find(&tags).Error
	return tags, err
}

// Count 统计标签数量
func (r *tagRepository) Count() (int64, error) {
	var total int64
	err := r.db.Model(&entity.Tag{}).Count(&total).Error
	return total, err
}

// GetTagArticleCount 获取标签文章数量
func (r *tagRepository) GetTagArticleCount(tagID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Article{}).
		Joins("JOIN article_tags ON article_tags.article_id = articles.id").
		Where("article_tags.tag_id = ?", tagID).
		Count(&count).Error
	return count, err
}

// FindByIDs 根据ID列表查找标签
func (r *tagRepository) FindByIDs(ids []uint) ([]*entity.Tag, error) {
	var tags []*entity.Tag
	err := r.db.Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}
