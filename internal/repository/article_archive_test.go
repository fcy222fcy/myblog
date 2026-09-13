package repository

import (
	"testing"

	"blog/internal/model/entity"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupArchiveDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&entity.Category{}, &entity.Tag{}, &entity.Article{},
		&entity.Comment{}, &entity.CommentLikeLog{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// TestArticleRepository_GetArchives_TrimmedColumns 归档查询必须裁掉 longtext 正文与关联预加载，
// 只返回 id/title/slug/created_at，且仅包含已发布文章。
func TestArticleRepository_GetArchives_TrimmedColumns(t *testing.T) {
	db := setupArchiveDB(t)
	repo := &articleRepository{db: db}

	cat := entity.Category{Name: "Tech", Slug: "tech"}
	if err := db.Create(&cat).Error; err != nil {
		t.Fatalf("seed category: %v", err)
	}
	tag := entity.Tag{Name: "go", Slug: "go"}
	if err := db.Create(&tag).Error; err != nil {
		t.Fatalf("seed tag: %v", err)
	}

	publishedA := entity.Article{
		Title: "A", Slug: "a", Content: "FULL-CONTENT-A", Summary: "sum-a",
		CategoryID: cat.ID, Status: entity.ArticleStatusPublished,
	}
	publishedB := entity.Article{
		Title: "B", Slug: "b", Content: "FULL-CONTENT-B",
		CategoryID: cat.ID, Status: entity.ArticleStatusPublished,
	}
	draft := entity.Article{Title: "D", Slug: "d", Content: "DRAFT", Status: entity.ArticleStatusDraft}
	for _, a := range []*entity.Article{&publishedA, &publishedB, &draft} {
		if err := db.Create(a).Error; err != nil {
			t.Fatalf("seed article: %v", err)
		}
	}
	if err := db.Model(&publishedA).Association("Tags").Append(&tag); err != nil {
		t.Fatalf("append tag: %v", err)
	}

	articles, err := repo.GetArchives()
	if err != nil {
		t.Fatalf("GetArchives: %v", err)
	}

	if len(articles) != 2 {
		t.Fatalf("归档应只返回 2 篇已发布文章，实际 %d", len(articles))
	}
	for _, a := range articles {
		if a.Content != "" {
			t.Errorf("归档查询不应读取正文，id=%d 实际 %q", a.ID, a.Content)
		}
		if a.Category.ID != 0 || len(a.Tags) != 0 {
			t.Errorf("归档查询不应预加载关联，id=%d category=%d tags=%d", a.ID, a.Category.ID, len(a.Tags))
		}
		if a.Title == "" || a.Slug == "" || a.CreatedAt.IsZero() {
			t.Errorf("归档查询缺少必要字段，id=%d title=%q slug=%q", a.ID, a.Title, a.Slug)
		}
	}
}

// TestArticleRepository_BatchDelete_CascadesInOneCall 批量删除应在一次调用内级联清理
// 评论、评论点赞、标签关联与文章本身。
func TestArticleRepository_BatchDelete_CascadesInOneCall(t *testing.T) {
	db := setupArchiveDB(t)
	repo := &articleRepository{db: db}

	articles := []*entity.Article{
		{Title: "A", Slug: "a", Status: entity.ArticleStatusPublished},
		{Title: "B", Slug: "b", Status: entity.ArticleStatusPublished},
		{Title: "C", Slug: "c", Status: entity.ArticleStatusPublished},
	}
	for _, a := range articles {
		if err := db.Create(a).Error; err != nil {
			t.Fatalf("seed article: %v", err)
		}
	}
	tag := entity.Tag{Name: "go", Slug: "go"}
	if err := db.Create(&tag).Error; err != nil {
		t.Fatalf("seed tag: %v", err)
	}

	comments := []*entity.Comment{
		{Content: "c1", ArticleID: articles[0].ID, IP: "1.1.1.1"},
		{Content: "c2", ArticleID: articles[1].ID, IP: "1.1.1.1"},
		{Content: "c3", ArticleID: articles[2].ID, IP: "1.1.1.1"},
	}
	for _, c := range comments {
		if err := db.Create(c).Error; err != nil {
			t.Fatalf("seed comment: %v", err)
		}
	}
	likes := []*entity.CommentLikeLog{
		{CommentID: comments[0].ID, VisitorIP: "1.1.1.1"},
		{CommentID: comments[1].ID, VisitorIP: "1.1.1.1"},
	}
	for _, l := range likes {
		if err := db.Create(l).Error; err != nil {
			t.Fatalf("seed like: %v", err)
		}
	}
	if err := db.Create(&entity.ArticleTag{ArticleID: articles[0].ID, TagID: tag.ID}).Error; err != nil {
		t.Fatalf("seed article tag: %v", err)
	}
	if err := db.Create(&entity.ArticleTag{ArticleID: articles[1].ID, TagID: tag.ID}).Error; err != nil {
		t.Fatalf("seed article tag: %v", err)
	}

	if err := repo.BatchDelete([]uint{articles[0].ID, articles[1].ID}); err != nil {
		t.Fatalf("BatchDelete: %v", err)
	}

	var articleCount, commentCount, likeCount, articleTagCount int64
	db.Model(&entity.Article{}).Count(&articleCount)
	db.Model(&entity.Comment{}).Count(&commentCount)
	db.Model(&entity.CommentLikeLog{}).Count(&likeCount)
	db.Model(&entity.ArticleTag{}).Count(&articleTagCount)

	if articleCount != 1 {
		t.Errorf("文章应剩 1 篇，实际 %d", articleCount)
	}
	if commentCount != 1 {
		t.Errorf("评论应剩 1 条（仅第三篇的），实际 %d", commentCount)
	}
	if likeCount != 0 {
		t.Errorf("被删文章的评论点赞应清零，实际 %d", likeCount)
	}
	if articleTagCount != 0 {
		t.Errorf("被删文章的标签关联应清零，实际 %d", articleTagCount)
	}
}

// TestArticleRepository_BatchDelete_EmptyIDs 空 ID 列表应直接返回，不报错。
func TestArticleRepository_BatchDelete_EmptyIDs(t *testing.T) {
	db := setupArchiveDB(t)
	repo := &articleRepository{db: db}

	if err := repo.BatchDelete(nil); err != nil {
		t.Fatalf("空列表批量删除不应报错，实际: %v", err)
	}
}
