package repository

import (
	"blog/internal/model/entity"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	ContentTypeArticle       = "article"
	ContentTypeDailyQuestion = "daily_question"
)

var (
	ErrUnsupportedContentType = errors.New("unsupported content type")
	ErrViewTargetNotFound     = errors.New("view target not found or not public")
)

type ContentViewRepository interface {
	RecordView(contentType string, contentID uint, visitorKey string, viewDate time.Time) (bool, int64, error)
	CountByDate(contentType string, viewDate time.Time) (int64, error)
}

type contentViewRepository struct {
	db *gorm.DB
}

func NewContentViewRepository(db *gorm.DB) ContentViewRepository {
	return &contentViewRepository{db: db}
}

func (r *contentViewRepository) RecordView(contentType string, contentID uint, visitorKey string, viewDate time.Time) (bool, int64, error) {
	if contentID == 0 || visitorKey == "" {
		return false, 0, ErrViewTargetNotFound
	}
	if contentType != ContentTypeArticle && contentType != ContentTypeDailyQuestion {
		return false, 0, ErrUnsupportedContentType
	}

	viewDate = calendarDate(viewDate)
	var counted bool
	var viewCount int64
	err := r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&entity.ContentViewEvent{
			ContentType: contentType,
			ContentID:   contentID,
			VisitorKey:  visitorKey,
			ViewDate:    viewDate,
		})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			current, err := currentViewCount(tx, contentType, contentID)
			if err != nil {
				return err
			}
			viewCount = current
			return nil
		}

		if err := incrementPublicViewCount(tx, contentType, contentID, viewDate); err != nil {
			return err
		}
		current, err := currentViewCount(tx, contentType, contentID)
		if err != nil {
			return err
		}
		counted = true
		viewCount = current
		return nil
	})
	return counted, viewCount, err
}

func (r *contentViewRepository) CountByDate(contentType string, viewDate time.Time) (int64, error) {
	if contentType != ContentTypeArticle && contentType != ContentTypeDailyQuestion {
		return 0, ErrUnsupportedContentType
	}
	var count int64
	err := r.db.Model(&entity.ContentViewEvent{}).
		Where("content_type = ? AND view_date = ?", contentType, calendarDate(viewDate)).
		Count(&count).Error
	return count, err
}

func incrementPublicViewCount(tx *gorm.DB, contentType string, contentID uint, viewDate time.Time) error {
	var result *gorm.DB
	switch contentType {
	case ContentTypeArticle:
		result = tx.Model(&entity.Article{}).
			Where("id = ? AND status = ?", contentID, entity.ArticleStatusPublished).
			UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	case ContentTypeDailyQuestion:
		result = tx.Model(&entity.DailyQuestion{}).
			Where("id = ? AND status = ? AND date <= ?", contentID, entity.DailyQuestionStatusPublished, viewDate.Format("2006-01-02")).
			UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	default:
		return ErrUnsupportedContentType
	}
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return ErrViewTargetNotFound
	}
	return nil
}

func currentViewCount(tx *gorm.DB, contentType string, contentID uint) (int64, error) {
	var viewCount int64
	var result *gorm.DB
	switch contentType {
	case ContentTypeArticle:
		result = tx.Model(&entity.Article{}).Select("view_count").Where("id = ?", contentID).Scan(&viewCount)
	case ContentTypeDailyQuestion:
		result = tx.Model(&entity.DailyQuestion{}).Select("view_count").Where("id = ?", contentID).Scan(&viewCount)
	default:
		return 0, ErrUnsupportedContentType
	}
	return viewCount, result.Error
}

func calendarDate(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, value.Location())
}
