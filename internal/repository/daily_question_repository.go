package repository

import (
	"blog/internal/model/entity"
	bizerrors "blog/pkg/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// dailyQuestionRepository 每日一问数据访问实现
type dailyQuestionRepository struct {
	db *gorm.DB
}

// NewDailyQuestionRepository 创建每日一问数据访问
func NewDailyQuestionRepository(db *gorm.DB) DailyQuestionRepository {
	return &dailyQuestionRepository{db: db}
}

// FindByID 根据 ID 查找问题
func (r *dailyQuestionRepository) FindByID(id uint) (*entity.DailyQuestion, error) {
	var question entity.DailyQuestion
	err := r.db.First(&question, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &question, nil
}

// GetAllPublished 获取所有已发布问题
func (r *dailyQuestionRepository) GetAllPublished() ([]*entity.DailyQuestion, error) {
	var questions []*entity.DailyQuestion
	err := r.db.Where("status = ?", entity.DailyQuestionStatusPublished).Order("date DESC").Find(&questions).Error
	return questions, err
}

// FindByDate 根据日期查找问题
func (r *dailyQuestionRepository) FindByDate(date string) (*entity.DailyQuestion, error) {
	var question entity.DailyQuestion
	err := r.db.Where("date = ?", date).First(&question).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &question, nil
}

// Create 创建问题
func (r *dailyQuestionRepository) Create(question *entity.DailyQuestion) error {
	return r.db.Create(question).Error
}

// Update 更新问题
func (r *dailyQuestionRepository) Update(question *entity.DailyQuestion) error {
	return r.db.Save(question).Error
}

// Delete 删除问题（硬删除，清理点赞记录）
func (r *dailyQuestionRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 清理点赞记录
		if err := tx.Where("question_id = ?", id).
			Delete(&entity.DailyQuestionLikeLog{}).Error; err != nil {
			return err
		}
		return tx.Delete(&entity.DailyQuestion{}, id).Error
	})
}

// List 问题列表（后台）
func (r *dailyQuestionRepository) List(offset, limit int, status int) ([]*entity.DailyQuestion, int64, error) {
	var questions []*entity.DailyQuestion
	var total int64

	query := r.db.Model(&entity.DailyQuestion{})
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).
		Order("date DESC").
		Find(&questions).Error
	return questions, total, err
}

// GetLatest 获取最新问题
func (r *dailyQuestionRepository) GetLatest() (*entity.DailyQuestion, error) {
	var question entity.DailyQuestion
	err := r.db.Where("status = ?", entity.DailyQuestionStatusPublished).Order("date DESC").First(&question).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &question, nil
}

// GetPrevious 获取前一天的问题
func (r *dailyQuestionRepository) GetPrevious(date string) (*entity.DailyQuestion, error) {
	var question entity.DailyQuestion
	err := r.db.Where("status = ? AND date < ?", entity.DailyQuestionStatusPublished, date).
		Order("date DESC").First(&question).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &question, nil
}

// GetNext 获取后一天的问题
func (r *dailyQuestionRepository) GetNext(date string) (*entity.DailyQuestion, error) {
	var question entity.DailyQuestion
	err := r.db.Where("status = ? AND date > ?", entity.DailyQuestionStatusPublished, date).
		Order("date ASC").First(&question).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &question, nil
}

// Count 统计问题数量
func (r *dailyQuestionRepository) Count(status int) (int64, error) {
	var total int64
	query := r.db.Model(&entity.DailyQuestion{})
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	err := query.Count(&total).Error
	return total, err
}

// SumViewCount 统计总浏览量
func (r *dailyQuestionRepository) SumViewCount() (int64, error) {
	var sum int64
	err := r.db.Model(&entity.DailyQuestion{}).Select("COALESCE(SUM(view_count), 0)").Scan(&sum).Error
	return sum, err
}

// SumLikeCount 统计总点赞数
func (r *dailyQuestionRepository) SumLikeCount() (int64, error) {
	var sum int64
	err := r.db.Model(&entity.DailyQuestion{}).Select("COALESCE(SUM(like_count), 0)").Scan(&sum).Error
	return sum, err
}

// IncrementViewCount 增加浏览量
func (r *dailyQuestionRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&entity.DailyQuestion{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// IncrementLikeCount 增加点赞数
func (r *dailyQuestionRepository) IncrementLikeCount(id uint) (int64, error) {
	err := r.db.Model(&entity.DailyQuestion{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
	if err != nil {
		return 0, err
	}
	var question entity.DailyQuestion
	err = r.db.Select("like_count").First(&question, id).Error
	return question.LikeCount, err
}

// LikeQuestionWithLog 点赞问题（事务 + IP 唯一约束防重复，并发安全）
// 依赖 daily_question_like_logs 表的 (question_id, visitor_ip) 唯一索引：
// 并发请求下只有一个能成功插入日志，其余被 OnConflict 忽略返回"已点赞"
func (r *dailyQuestionRepository) LikeQuestionWithLog(questionID uint, visitorIP string) (int64, error) {
	var likeCount int64
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 先插入点赞日志（唯一索引冲突时 DoNothing）
		result := tx.Clauses(clause.OnConflict{DoNothing: true}).
			Create(&entity.DailyQuestionLikeLog{
				QuestionID: questionID,
				VisitorIP:  visitorIP,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return bizerrors.New(bizerrors.CodeAlreadyLiked, "已经点过赞了")
		}

		// 2. 增加点赞数
		if err := tx.Model(&entity.DailyQuestion{}).Where("id = ?", questionID).
			UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
			return err
		}

		// 3. 返回最新点赞数
		var question entity.DailyQuestion
		if err := tx.Select("like_count").First(&question, questionID).Error; err != nil {
			return err
		}
		likeCount = question.LikeCount
		return nil
	})
	return likeCount, err
}

// BatchUpdateStatus 批量更新状态
func (r *dailyQuestionRepository) BatchUpdateStatus(ids []uint, status int) error {
	return r.db.Model(&entity.DailyQuestion{}).Where("id IN ?", ids).
		Update("status", status).Error
}

// BatchDelete 批量删除
func (r *dailyQuestionRepository) PublishScheduledQuestions(today string) (int64, error) {
	result := r.db.Model(&entity.DailyQuestion{}).
		Where("status = ? AND date <= ?", entity.DailyQuestionStatusScheduled, today).
		Update("status", entity.DailyQuestionStatusPublished)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (r *dailyQuestionRepository) BatchDelete(ids []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 清理点赞记录
		if err := tx.Where("question_id IN ?", ids).
			Delete(&entity.DailyQuestionLikeLog{}).Error; err != nil {
			return err
		}
		return tx.Where("id IN ?", ids).
			Delete(&entity.DailyQuestion{}).Error
	})
}
