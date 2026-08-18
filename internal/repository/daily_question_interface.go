package repository

import "blog/internal/model/entity"

// DailyQuestionRepository 每日一问数据访问接口
type DailyQuestionRepository interface {
	// FindByID 根据 ID 查找问题
	FindByID(id uint) (*entity.DailyQuestion, error)

	// GetAllPublished 获取所有已发布问题
	GetAllPublished() ([]*entity.DailyQuestion, error)

	// FindByDate 根据日期查找问题
	FindByDate(date string) (*entity.DailyQuestion, error)

	// Create 创建问题
	Create(question *entity.DailyQuestion) error

	// Update 更新问题
	Update(question *entity.DailyQuestion) error

	// Delete 删除问题（硬删除，清理点赞记录）
	Delete(id uint) error

	// List 问题列表（后台）
	List(offset, limit int, status int) ([]*entity.DailyQuestion, int64, error)

	// GetLatest 获取最新问题
	GetLatest() (*entity.DailyQuestion, error)

	// GetPrevious 获取前一天的问题
	GetPrevious(date string) (*entity.DailyQuestion, error)

	// GetNext 获取后一天的问题
	GetNext(date string) (*entity.DailyQuestion, error)

	// Count 统计问题数量
	Count(status int) (int64, error)

	// SumViewCount 统计总浏览量
	SumViewCount() (int64, error)

	// SumLikeCount 统计总点赞数
	SumLikeCount() (int64, error)

	// IncrementViewCount 增加浏览量
	IncrementViewCount(id uint) error

	// IncrementLikeCount 增加点赞数
	IncrementLikeCount(id uint) (int64, error)

	// LikeQuestionWithLog 点赞问题（事务 + IP 唯一约束防重复）
	LikeQuestionWithLog(questionID uint, visitorIP string) (int64, error)

	// BatchUpdateStatus 批量更新状态
	BatchUpdateStatus(ids []uint, status int) error

	// BatchDelete 批量删除
	PublishScheduledQuestions(today string) (int64, error)

	BatchDelete(ids []uint) error
}
