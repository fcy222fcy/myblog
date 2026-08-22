package service

import (
	"blog/internal/repository"
	"errors"
	"time"
)

type ContentViewResult struct {
	Counted   bool  `json:"counted"`
	ViewCount int64 `json:"view_count"`
}

type ContentViewService interface {
	Record(contentType string, contentID uint, visitorKey string) (*ContentViewResult, error)
	CountToday(contentType string) (int64, error)
}

type contentViewService struct {
	repo repository.ContentViewRepository
	now  func() time.Time
}

func NewContentViewService(repo repository.ContentViewRepository) ContentViewService {
	return &contentViewService{repo: repo, now: time.Now}
}

func (s *contentViewService) Record(contentType string, contentID uint, visitorKey string) (*ContentViewResult, error) {
	if contentType != repository.ContentTypeArticle && contentType != repository.ContentTypeDailyQuestion {
		return nil, repository.ErrUnsupportedContentType
	}
	if contentID == 0 || visitorKey == "" {
		return nil, errors.New("content id and visitor key are required")
	}

	counted, viewCount, err := s.repo.RecordView(contentType, contentID, visitorKey, shanghaiNow(s.now()))
	if err != nil {
		return nil, err
	}
	return &ContentViewResult{Counted: counted, ViewCount: viewCount}, nil
}

func (s *contentViewService) CountToday(contentType string) (int64, error) {
	return s.repo.CountByDate(contentType, shanghaiNow(s.now()))
}

func shanghaiNow(value time.Time) time.Time {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		location = time.FixedZone("Asia/Shanghai", 8*60*60)
	}
	return value.In(location)
}
