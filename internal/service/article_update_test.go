package service

import (
	"testing"
	"time"

	"blog/internal/model/dto/request"
	"blog/internal/model/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newUpdateTestSvc 构造一篇文章用于更新场景测试：
// 已置顶 + 已设置定时发布时间，便于验证「部分更新」不会误清这些字段。
func newUpdateTestSvc(t *testing.T) (*mockArticleRepository, ArticleService, uint) {
	t.Helper()

	scheduled := time.Date(2026, 9, 20, 8, 0, 0, 0, time.UTC)
	repo := newMockArticleRepo()
	svc := NewArticleService(repo, nil, nil, nil)

	article := &entity.Article{
		Title:       "原始标题",
		Slug:        "orig-slug",
		Content:     "原始内容",
		CategoryID:  1,
		Status:      entity.ArticleStatusScheduled,
		IsTop:       true,
		ScheduledAt: &scheduled,
	}
	require.NoError(t, repo.Create(article))

	return repo, svc, article.ID
}

// TestArticleService_Update_IsTopPartialUpdate 回归：
// 请求未提交 is_top 时必须保留原置顶状态，只有显式提交才允许改动。
func TestArticleService_Update_IsTopPartialUpdate(t *testing.T) {
	t.Run("未提交 is_top 时保留置顶", func(t *testing.T) {
		repo, svc, id := newUpdateTestSvc(t)

		require.NoError(t, svc.UpdateArticle(id, &request.UpdateArticleRequest{Title: "新标题"}))

		updated, err := repo.FindByID(id)
		require.NoError(t, err)
		assert.True(t, updated.IsTop, "未提交 is_top 不应取消置顶")
	})

	t.Run("显式提交 is_top=false 时取消置顶", func(t *testing.T) {
		repo, svc, id := newUpdateTestSvc(t)
		notTop := false

		require.NoError(t, svc.UpdateArticle(id, &request.UpdateArticleRequest{IsTop: &notTop}))

		updated, err := repo.FindByID(id)
		require.NoError(t, err)
		assert.False(t, updated.IsTop)
	})

	t.Run("显式提交 is_top=true 时保持置顶", func(t *testing.T) {
		repo, svc, id := newUpdateTestSvc(t)
		top := true

		require.NoError(t, svc.UpdateArticle(id, &request.UpdateArticleRequest{IsTop: &top}))

		updated, err := repo.FindByID(id)
		require.NoError(t, err)
		assert.True(t, updated.IsTop)
	})
}

// TestArticleService_Update_ScheduledAtPartialUpdate 回归：
// 请求未提交 scheduled_at 时必须保留原定时发布时间；
// 显式传 null 时才允许清空。
func TestArticleService_Update_ScheduledAtPartialUpdate(t *testing.T) {
	original := time.Date(2026, 9, 20, 8, 0, 0, 0, time.UTC)

	t.Run("未提交 scheduled_at 时保留原值", func(t *testing.T) {
		repo, svc, id := newUpdateTestSvc(t)

		require.NoError(t, svc.UpdateArticle(id, &request.UpdateArticleRequest{Title: "新标题"}))

		updated, err := repo.FindByID(id)
		require.NoError(t, err)
		require.NotNil(t, updated.ScheduledAt, "未提交 scheduled_at 不应清空定时发布时间")
		assert.True(t, original.Equal(*updated.ScheduledAt))
	})

	t.Run("显式提交 scheduled_at=null 时清空", func(t *testing.T) {
		repo, svc, id := newUpdateTestSvc(t)

		req := &request.UpdateArticleRequest{
			ScheduledAt: request.OptionalTime{Set: true, Value: nil},
		}
		require.NoError(t, svc.UpdateArticle(id, req))

		updated, err := repo.FindByID(id)
		require.NoError(t, err)
		assert.Nil(t, updated.ScheduledAt)
	})

	t.Run("显式提交新的 scheduled_at 时更新", func(t *testing.T) {
		repo, svc, id := newUpdateTestSvc(t)
		next := time.Date(2026, 10, 1, 9, 30, 0, 0, time.UTC)

		req := &request.UpdateArticleRequest{
			ScheduledAt: request.OptionalTime{Set: true, Value: &next},
		}
		require.NoError(t, svc.UpdateArticle(id, req))

		updated, err := repo.FindByID(id)
		require.NoError(t, err)
		require.NotNil(t, updated.ScheduledAt)
		assert.True(t, next.Equal(*updated.ScheduledAt))
	})
}

// TestArticleCacheKeys 回归：文章更新需要失效的缓存 key 计算正确。
// 改 slug 时旧、新地址都要清，否则旧链接会继续命中陈旧缓存。
func TestArticleCacheKeys(t *testing.T) {
	t.Run("改 slug 时同时包含新旧详情 key 与列表 key", func(t *testing.T) {
		keys := articleCacheKeys("old-slug", "new-slug")

		assert.Contains(t, keys, "article:detail:old-slug")
		assert.Contains(t, keys, "article:detail:new-slug")
		assert.Contains(t, keys, "category:list")
		assert.Contains(t, keys, "tag:list")
		assert.Contains(t, keys, "article:archives")
	})

	t.Run("slug 未变化时不重复生成详情 key", func(t *testing.T) {
		keys := articleCacheKeys("same-slug", "same-slug")

		count := 0
		for _, k := range keys {
			if k == "article:detail:same-slug" {
				count++
			}
		}
		assert.Equal(t, 1, count)
	})

	t.Run("空 slug 被忽略且不影响列表 key", func(t *testing.T) {
		keys := articleCacheKeys("", "")

		for _, k := range keys {
			assert.NotContains(t, k, "article:detail:")
		}
		assert.Contains(t, keys, "category:list")
		assert.Contains(t, keys, "article:archives")
	})

	t.Run("单个 slug 时行为与旧实现一致", func(t *testing.T) {
		keys := articleCacheKeys("only-slug")

		assert.Contains(t, keys, "article:detail:only-slug")
		assert.Len(t, keys, 4)
	})
}
