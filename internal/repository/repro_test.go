package repository

import (
	"testing"

	"blog/internal/model/entity"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupArticleDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&entity.Category{}, &entity.Tag{}, &entity.Article{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	cats := []entity.Category{{Name: "Tech", Slug: "tech"}, {Name: "Life", Slug: "life"}}
	if err := db.Create(&cats).Error; err != nil {
		t.Fatalf("seed category: %v", err)
	}
	if err := db.Create(&entity.Article{Title: "orig", Slug: "orig", CategoryID: cats[0].ID}).Error; err != nil {
		t.Fatalf("seed article: %v", err)
	}
	return db
}

func rawCategoryID(t *testing.T, db *gorm.DB, id uint) uint {
	var cid uint
	if err := db.Raw("SELECT category_id FROM articles WHERE id = ?", id).Scan(&cid).Error; err != nil {
		t.Fatalf("raw scan: %v", err)
	}
	return cid
}

// TestArticleRepository_Update_CategoryChange 回归：更新文章时修改 category_id 必须生效，
// 不能因预加载的 Category 关联被 Save 同步回旧值。
func TestArticleRepository_Update_CategoryChange(t *testing.T) {
	db := setupArticleDB(t)
	repo := &articleRepository{db: db}

	var article entity.Article
	if err := db.Preload("Category").Preload("Tags").First(&article, 1).Error; err != nil {
		t.Fatalf("find: %v", err)
	}

	article.CategoryID = 2 // 期望从 Tech(1) 改为 Life(2)
	article.Title = "updated"
	article.Summary = "new-sum"
	if err := repo.Update(&article); err != nil {
		t.Fatalf("update: %v", err)
	}

	if got := rawCategoryID(t, db, 1); got != 2 {
		t.Fatalf("category_id 未更新，期望 2，实际 %d", got)
	}
	var saved entity.Article
	db.First(&saved, 1)
	if saved.Title != "updated" {
		t.Fatalf("Title 未更新，实际 %q", saved.Title)
	}
	if saved.Summary != "new-sum" {
		t.Fatalf("Summary 未更新，实际 %q", saved.Summary)
	}
}

// TestArticleRepository_Update_TagsPreserved 回归：未改动标签时，Update 不应清空多对多关联。
func TestArticleRepository_Update_TagsPreserved(t *testing.T) {
	db := setupArticleDB(t)
	repo := &articleRepository{db: db}

	tag := entity.Tag{Name: "go", Slug: "go"}
	if err := db.Create(&tag).Error; err != nil {
		t.Fatalf("seed tag: %v", err)
	}
	var a entity.Article
	db.First(&a, 1)
	if err := db.Model(&a).Association("Tags").Append(&tag); err != nil {
		t.Fatalf("append tag: %v", err)
	}

	var article entity.Article
	db.Preload("Category").Preload("Tags").First(&article, 1)
	article.Title = "touch" // 只改标题，不动标签
	if err := repo.Update(&article); err != nil {
		t.Fatalf("update: %v", err)
	}

	var after entity.Article
	db.Preload("Tags").First(&after, 1)
	if len(after.Tags) != 1 {
		t.Fatalf("标签关联被错误清空，期望保留 1 个，实际 %d", len(after.Tags))
	}
}
