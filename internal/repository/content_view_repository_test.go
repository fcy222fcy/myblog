package repository

import (
	"blog/internal/model/entity"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newContentViewTestRepository(t *testing.T) (*contentViewRepository, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, db.AutoMigrate(&entity.Article{}, &entity.DailyQuestion{}, &entity.ContentViewEvent{}))
	return NewContentViewRepository(db).(*contentViewRepository), db
}

func TestContentViewRepositoryRecordsOneViewPerVisitorAndDay(t *testing.T) {
	repo, db := newContentViewTestRepository(t)
	article := entity.Article{Title: "published", Slug: "published", Status: entity.ArticleStatusPublished}
	require.NoError(t, db.Create(&article).Error)
	day := time.Date(2026, 8, 22, 0, 0, 0, 0, time.FixedZone("Asia/Shanghai", 8*60*60))

	counted, count, err := repo.RecordView(ContentTypeArticle, article.ID, "visitor-a", day)
	require.NoError(t, err)
	require.True(t, counted)
	require.Equal(t, int64(1), count)

	counted, count, err = repo.RecordView(ContentTypeArticle, article.ID, "visitor-a", day.Add(12*time.Hour))
	require.NoError(t, err)
	require.False(t, counted)
	require.Equal(t, int64(1), count)

	counted, count, err = repo.RecordView(ContentTypeArticle, article.ID, "visitor-a", day.Add(24*time.Hour))
	require.NoError(t, err)
	require.True(t, counted)
	require.Equal(t, int64(2), count)
}

func TestContentViewRepositorySupportsDailyQuestions(t *testing.T) {
	repo, db := newContentViewTestRepository(t)
	question := entity.DailyQuestion{Question: "Q", Date: "2026-08-22", Status: entity.DailyQuestionStatusPublished}
	require.NoError(t, db.Create(&question).Error)
	day := time.Date(2026, 8, 22, 0, 0, 0, 0, time.Local)

	counted, count, err := repo.RecordView(ContentTypeDailyQuestion, question.ID, "visitor-a", day)
	require.NoError(t, err)
	require.True(t, counted)
	require.Equal(t, int64(1), count)
}

func TestContentViewRepositoryRollsBackInvalidTargets(t *testing.T) {
	repo, db := newContentViewTestRepository(t)
	draft := entity.Article{Title: "draft", Slug: "draft", Status: entity.ArticleStatusDraft}
	require.NoError(t, db.Create(&draft).Error)

	_, _, err := repo.RecordView(ContentTypeArticle, draft.ID, "visitor-a", time.Now())
	require.ErrorIs(t, err, ErrViewTargetNotFound)

	var eventCount int64
	require.NoError(t, db.Model(&entity.ContentViewEvent{}).Count(&eventCount).Error)
	require.Zero(t, eventCount)
}

func TestContentViewRepositoryIsSafeForConcurrentDuplicates(t *testing.T) {
	repo, db := newContentViewTestRepository(t)
	article := entity.Article{Title: "published", Slug: "concurrent", Status: entity.ArticleStatusPublished}
	require.NoError(t, db.Create(&article).Error)
	day := time.Date(2026, 8, 22, 0, 0, 0, 0, time.Local)

	var countedTotal atomic.Int64
	var wg sync.WaitGroup
	errs := make(chan error, 20)
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counted, _, err := repo.RecordView(ContentTypeArticle, article.ID, "visitor-a", day)
			if err != nil {
				errs <- err
				return
			}
			if counted {
				countedTotal.Add(1)
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}
	require.Equal(t, int64(1), countedTotal.Load())

	var stored entity.Article
	require.NoError(t, db.First(&stored, article.ID).Error)
	require.Equal(t, int64(1), stored.ViewCount)
}
