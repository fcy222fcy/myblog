package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

// userRepository 用户数据访问实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户数据访问
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// FindByID 根据 ID 查找用户
func (r *userRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *userRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindFirstAdmin 查找第一个管理员用户（博客主人）
func (r *userRepository) FindFirstAdmin() (*entity.User, error) {
	var user entity.User
	// 获取第一个启用的用户作为博客主人
	err := r.db.Where("status = ?", 1).Order("id ASC").First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户（硬删除）
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}

// List 用户列表
func (r *userRepository) List(offset, limit int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	err := r.db.Model(&entity.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
