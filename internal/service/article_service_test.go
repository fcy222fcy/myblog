package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/entity"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// mockArticleRepository 文章仓库的模拟实现
type mockArticleRepository struct {
	articles map[uint]*entity.Article
	nextID   uint
}

func newMockArticleRepo() *mockArticleRepository {
	return &mockArticleRepository{
		articles: make(map[uint]*entity.Article),
		nextID:   1,
	}
}

func (m *mockArticleRepository) FindByID(id uint) (*entity.Article, error) {
	if article, ok := m.articles[id]; ok {
		return article, nil
	}
	return nil, fmt.Errorf("article not found")
}

func (m *mockArticleRepository) FindBySlug(slug string) (*entity.Article, error) {
	for _, article := range m.articles {
		if article.Slug == slug {
			return article, nil
		}
	}
	return nil, fmt.Errorf("article not found")
}

func (m *mockArticleRepository) Create(article *entity.Article) error {
	article.ID = m.nextID
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	m.articles[article.ID] = article
	m.nextID++
	return nil
}

func (m *mockArticleRepository) Update(article *entity.Article) error {
	m.articles[article.ID] = article
	return nil
}

func (m *mockArticleRepository) Delete(id uint) error {
	delete(m.articles, id)
	return nil
}

func (m *mockArticleRepository) ListPublished(offset, limit int, categoryId, tagId uint, keyword string) ([]*entity.Article, int64, error) {
	var result []*entity.Article
	for _, article := range m.articles {
		if article.Status == entity.ArticleStatusPublished {
			result = append(result, article)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockArticleRepository) ListAll(offset, limit int, status, keyword string, categoryID uint) ([]*entity.Article, int64, error) {
	var result []*entity.Article
	for _, article := range m.articles {
		result = append(result, article)
	}
	return result, int64(len(result)), nil
}

func (m *mockArticleRepository) IncrementViewCount(id uint) error {
	if article, ok := m.articles[id]; ok {
		article.ViewCount++
	}
	return nil
}

func (m *mockArticleRepository) UpdateCommentCount(articleID uint, delta int64) error {
	if article, ok := m.articles[articleID]; ok {
		article.CommentCount += delta
		if article.CommentCount < 0 {
			article.CommentCount = 0
		}
	}
	return nil
}

func (m *mockArticleRepository) BatchDelete(ids []uint) error {
	for _, id := range ids {
		delete(m.articles, id)
	}
	return nil
}

func (m *mockArticleRepository) Count(status string) (int64, error) {
	var count int64
	for _, article := range m.articles {
		if status == "" || article.Status == status {
			count++
		}
	}
	return count, nil
}

func (m *mockArticleRepository) SumViewCount() (int64, error) {
	var total int64
	for _, article := range m.articles {
		total += article.ViewCount
	}
	return total, nil
}

func (m *mockArticleRepository) FindByCategoryID(categoryID uint, offset, limit int) ([]*entity.Article, int64, error) {
	return nil, 0, nil
}

func (m *mockArticleRepository) FindByTagID(tagID uint, offset, limit int) ([]*entity.Article, int64, error) {
	return nil, 0, nil
}

func (m *mockArticleRepository) GetArchives() ([]map[string][]*entity.Article, error) {
	return nil, nil
}

func (m *mockArticleRepository) GetRecent(limit int) ([]entity.Article, error) {
	return nil, nil
}

func (m *mockArticleRepository) Search(keyword string, offset, limit int) ([]*entity.Article, int64, error) {
	return nil, 0, nil
}

func (m *mockArticleRepository) UpdateTags(article *entity.Article, tags []entity.Tag) error {
	article.Tags = tags
	return nil
}

func (m *mockArticleRepository) ListScheduledAfter(now time.Time) ([]*entity.Article, error) {
	var result []*entity.Article
	for _, article := range m.articles {
		if article.Status == entity.ArticleStatusScheduled && article.ScheduledAt != nil && article.ScheduledAt.After(now) {
			result = append(result, article)
		}
	}
	return result, nil
}

func (m *mockArticleRepository) PublishScheduledArticle(id uint, now time.Time) (bool, error) {
	article, ok := m.articles[id]
	if !ok || article.Status != entity.ArticleStatusScheduled || article.ScheduledAt == nil || article.ScheduledAt.After(now) {
		return false, nil
	}
	article.Status = entity.ArticleStatusPublished
	article.ScheduledAt = nil
	return true, nil
}

func (m *mockArticleRepository) PublishDueScheduledArticles(now time.Time) (int64, error) {
	var count int64
	for _, article := range m.articles {
		if article.Status == entity.ArticleStatusScheduled && article.ScheduledAt != nil && !article.ScheduledAt.After(now) {
			article.Status = entity.ArticleStatusPublished
			article.ScheduledAt = nil
			count++
		}
	}
	return count, nil
}

func (m *mockArticleRepository) GetDB() *gorm.DB {
	return nil
}

// calculateReadingTime 测试
func TestCalculateReadingTime(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "空内容",
			content:  "",
			expected: 1,
		},
		{
			name:     "短内容",
			content:  "Hello World",
			expected: 1,
		},
		{
			name:     "中文内容约400字",
			content:  generateChineseText(400),
			expected: 1,
		},
		{
			name:     "中文内容约800字",
			content:  generateChineseText(800),
			expected: 2,
		},
		{
			name:     "Markdown格式",
			content:  "# 标题\n\n这是一段测试内容，用来验证Markdown格式下的阅读时间计算。\n\n## 小标题\n\n更多内容在这里。",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateReadingTime(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestArticleServiceGetDetailDoesNotIncrementViews(t *testing.T) {
	articleRepo := newMockArticleRepo()
	article := &entity.Article{
		Title:     "只读详情",
		Slug:      "read-only-detail",
		Status:    entity.ArticleStatusPublished,
		ViewCount: 6,
	}
	require.NoError(t, articleRepo.Create(article))
	svc := NewArticleService(articleRepo, nil, nil, nil)

	first, err := svc.GetArticleDetail(article.Slug)
	require.NoError(t, err)
	second, err := svc.GetArticleDetail(article.Slug)
	require.NoError(t, err)

	require.Equal(t, int64(6), first.ViewCount)
	require.Equal(t, int64(6), second.ViewCount)
	require.Equal(t, int64(6), articleRepo.articles[article.ID].ViewCount)
}

// generateChineseText 生成指定字数的中文文本
func generateChineseText(count int) string {
	text := ""
	for i := 0; i < count; i++ {
		text += "测"
	}
	return text
}

// TestArticleService_Create 测试创建文章
func TestArticleService_Create(t *testing.T) {
	articleRepo := newMockArticleRepo()
	svc := NewArticleService(articleRepo, nil, nil, nil)

	t.Run("成功创建文章", func(t *testing.T) {
		req := &request.CreateArticleRequest{
			Title:      "测试文章",
			Content:    "这是测试内容，至少需要一些文字来验证。",
			Summary:    "测试摘要",
			CategoryID: 1,
			Status:     "draft",
		}

		id, err := svc.CreateArticle(req)
		assert.NoError(t, err)
		assert.Greater(t, id, uint(0))

		article, err := articleRepo.FindByID(id)
		assert.NoError(t, err)
		assert.Equal(t, "测试文章", article.Title)
		assert.Equal(t, "draft", article.Status)
	})

	t.Run("创建文章时计算阅读时间", func(t *testing.T) {
		req := &request.CreateArticleRequest{
			Title:      "阅读时间测试",
			Content:    generateChineseText(800),
			CategoryID: 1,
		}

		id, err := svc.CreateArticle(req)
		assert.NoError(t, err)

		article, err := articleRepo.FindByID(id)
		assert.NoError(t, err)
		assert.Equal(t, 2, article.ReadingTime)
	})

	t.Run("默认状态为草稿", func(t *testing.T) {
		req := &request.CreateArticleRequest{
			Title:      "默认状态测试",
			Content:    "内容",
			CategoryID: 1,
		}

		id, err := svc.CreateArticle(req)
		assert.NoError(t, err)

		article, err := articleRepo.FindByID(id)
		assert.NoError(t, err)
		assert.Equal(t, entity.ArticleStatusDraft, article.Status)
	})
}

// TestArticleService_Update 测试更新文章
func TestArticleService_Update(t *testing.T) {
	articleRepo := newMockArticleRepo()
	svc := NewArticleService(articleRepo, nil, nil, nil)

	article := &entity.Article{
		Title:      "原始标题",
		Content:    "原始内容",
		CategoryID: 1,
		Status:     "draft",
	}
	articleRepo.Create(article)

	t.Run("更新标题", func(t *testing.T) {
		req := &request.UpdateArticleRequest{
			Title: "新标题",
		}

		err := svc.UpdateArticle(article.ID, req)
		assert.NoError(t, err)

		updated, _ := articleRepo.FindByID(article.ID)
		assert.Equal(t, "新标题", updated.Title)
	})

	t.Run("更新内容时重新计算阅读时间", func(t *testing.T) {
		req := &request.UpdateArticleRequest{
			Content: generateChineseText(1200),
		}

		err := svc.UpdateArticle(article.ID, req)
		assert.NoError(t, err)

		updated, _ := articleRepo.FindByID(article.ID)
		assert.Equal(t, 3, updated.ReadingTime)
	})

	t.Run("更新不存在的文章", func(t *testing.T) {
		req := &request.UpdateArticleRequest{
			Title: "不存在",
		}

		err := svc.UpdateArticle(9999, req)
		assert.Error(t, err)
	})
}

// TestArticleService_Delete 测试删除文章
func TestArticleService_Delete(t *testing.T) {
	articleRepo := newMockArticleRepo()
	svc := NewArticleService(articleRepo, nil, nil, nil)

	article := &entity.Article{
		Title:      "待删除文章",
		Content:    "内容",
		CategoryID: 1,
	}
	articleRepo.Create(article)

	t.Run("成功删除", func(t *testing.T) {
		err := svc.DeleteArticle(article.ID)
		assert.NoError(t, err)

		_, err = articleRepo.FindByID(article.ID)
		assert.Error(t, err)
	})

	t.Run("删除不存在的文章", func(t *testing.T) {
		err := svc.DeleteArticle(9999)
		assert.Error(t, err)
	})
}

// TestArticleService_BatchDelete 测试批量删除
func TestArticleService_BatchDelete(t *testing.T) {
	articleRepo := newMockArticleRepo()
	svc := NewArticleService(articleRepo, nil, nil, nil)

	for i := 0; i < 3; i++ {
		articleRepo.Create(&entity.Article{
			Title:      fmt.Sprintf("文章%d", i+1),
			Content:    "内容",
			CategoryID: 1,
		})
	}

	err := svc.BatchDeleteArticles([]uint{1, 2, 3})
	assert.NoError(t, err)
	assert.Len(t, articleRepo.articles, 0)
}
