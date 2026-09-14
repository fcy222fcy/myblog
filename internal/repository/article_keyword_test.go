package repository

import (
	"fmt"
	"testing"

	"blog/internal/model/entity"
)

// TestEscapeLikePattern escapeLikePattern 必须先转义反斜杠，再转义 % 与 _，
// 否则后两步引入的反斜杠会被 MySQL 当成转义符吞掉。
func TestEscapeLikePattern(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "普通字符串原样返回", in: "golang", want: "golang"},
		{name: "百分号转义", in: "50%", want: `50\%`},
		{name: "下划线转义", in: "a_b", want: `a\_b`},
		{name: "反斜杠自身转义", in: `a\b`, want: `a\\b`},
		{name: "反斜杠加通配符", in: `a\%`, want: `a\\\%`},
		{name: "混合", in: `100%_x\y`, want: `100\%\_x\\y`},
		{name: "空串", in: "", want: ""},
	}

	for _, c := range cases {
		if got := escapeLikePattern(c.in); got != c.want {
			t.Errorf("%s: escapeLikePattern(%q) = %q, 期望 %q", c.name, c.in, got, c.want)
		}
	}
}

// TestArticleRepository_KeywordIsEscaped 关键词必须走 escapeLikePattern，
// 否则用户输入的 % / _ 会被当成 LIKE 通配符，退化成全表命中。
func TestArticleRepository_KeywordIsEscaped(t *testing.T) {
	db := setupArchiveDB(t)
	repo := &articleRepository{db: db}

	// 标题刻意不含 % _ \，让「搜索通配符应命中 0 条」这个断言
	// 在 MySQL（\ 是 LIKE 转义符）与 sqlite（测试用）下结论一致。
	for i, title := range []string{"Go 入门", "MySQL 索引优化", "并发编程实践"} {
		a := entity.Article{
			Title:  title,
			Slug:   fmt.Sprintf("kw-%d", i),
			Status: entity.ArticleStatusPublished,
		}
		if err := db.Create(&a).Error; err != nil {
			t.Fatalf("seed article: %v", err)
		}
	}

	// 正常关键词仍要能命中
	got, total, err := repo.ListPublished(0, 10, 0, 0, "MySQL")
	if err != nil {
		t.Fatalf("ListPublished: %v", err)
	}
	if total != 1 || len(got) != 1 {
		t.Fatalf("搜索 MySQL 应命中 1 篇，实际 total=%d len=%d", total, len(got))
	}

	// 通配符不能被当成「匹配任意」，否则会退化成全表命中
	for _, kw := range []string{"%", "_", "%%"} {
		_, total, err := repo.ListPublished(0, 10, 0, 0, kw)
		if err != nil {
			t.Fatalf("ListPublished(%q): %v", kw, err)
		}
		if total != 0 {
			t.Errorf("搜索 %q 应命中 0 篇（通配符已被转义），实际 %d 篇", kw, total)
		}
	}

	// 后台列表 ListAll 走的是同一条过滤路径，同样验证
	_, allTotal, err := repo.ListAll(0, 10, "", "MySQL", 0)
	if err != nil {
		t.Fatalf("ListAll: %v", err)
	}
	if allTotal != 1 {
		t.Errorf("ListAll 搜索 MySQL 应命中 1 篇，实际 %d", allTotal)
	}

	_, pctTotal, err := repo.ListAll(0, 10, "", "%", 0)
	if err != nil {
		t.Fatalf("ListAll(%q): %v", "%", err)
	}
	if pctTotal != 0 {
		t.Errorf("ListAll 搜索 %% 应命中 0 篇，实际 %d 篇", pctTotal)
	}
}
