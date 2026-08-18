package repository

import "blog/internal/model/entity"

// UserRepository 用户数据访问接口
type UserRepository interface {
	// FindByID 根据 ID 查找用户
	FindByID(id uint) (*entity.User, error)

	// FindByUsername 根据用户名查找用户
	FindByUsername(username string) (*entity.User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(email string) (*entity.User, error)

	// FindFirstAdmin 查找第一个管理员用户（博客主人）
	FindFirstAdmin() (*entity.User, error)

	// Create 创建用户
	Create(user *entity.User) error

	// Update 更新用户
	Update(user *entity.User) error

	// Delete 删除用户（硬删除）
	Delete(id uint) error

	// List 用户列表
	List(offset, limit int) ([]*entity.User, int64, error)
}
