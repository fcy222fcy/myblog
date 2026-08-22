package service

import (
	"blog/internal/repository"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fakeContentViewRepository struct {
	contentType string
	contentID   uint
	visitorKey  string
	viewDate    time.Time
	counted     bool
	viewCount   int64
	err         error
}

func (f *fakeContentViewRepository) RecordView(contentType string, contentID uint, visitorKey string, viewDate time.Time) (bool, int64, error) {
	f.contentType = contentType
	f.contentID = contentID
	f.visitorKey = visitorKey
	f.viewDate = viewDate
	return f.counted, f.viewCount, f.err
}

func (f *fakeContentViewRepository) CountByDate(contentType string, viewDate time.Time) (int64, error) {
	f.contentType = contentType
	f.viewDate = viewDate
	return f.viewCount, f.err
}

func TestContentViewServiceRecordsUsingShanghaiCalendarDate(t *testing.T) {
	repo := &fakeContentViewRepository{counted: true, viewCount: 9}
	svc := NewContentViewService(repo).(*contentViewService)
	svc.now = func() time.Time {
		return time.Date(2026, 8, 21, 16, 30, 0, 0, time.UTC)
	}

	result, err := svc.Record(repository.ContentTypeArticle, 7, "visitor-key")
	require.NoError(t, err)
	require.True(t, result.Counted)
	require.Equal(t, int64(9), result.ViewCount)
	require.Equal(t, "2026-08-22", repo.viewDate.Format("2006-01-02"))
}

func TestContentViewServiceRejectsInvalidInput(t *testing.T) {
	svc := NewContentViewService(&fakeContentViewRepository{})

	_, err := svc.Record("unknown", 1, "visitor-key")
	require.ErrorIs(t, err, repository.ErrUnsupportedContentType)

	_, err = svc.Record(repository.ContentTypeArticle, 0, "visitor-key")
	require.Error(t, err)

	_, err = svc.Record(repository.ContentTypeArticle, 1, "")
	require.Error(t, err)
}

func TestContentViewServicePropagatesRepositoryErrors(t *testing.T) {
	expected := errors.New("database unavailable")
	svc := NewContentViewService(&fakeContentViewRepository{err: expected})

	_, err := svc.Record(repository.ContentTypeArticle, 1, "visitor-key")
	require.ErrorIs(t, err, expected)
}
