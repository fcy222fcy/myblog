package request

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUpdateArticleRequest_PartialUpdateSemantics 回归：更新文章的请求体里，
// 「字段未出现」必须能与「显式传零值」区分开。否则前端未提交 is_top /
// scheduled_at 时，后端会把零值当成用户意图，静默覆盖掉原有数据。
func TestUpdateArticleRequest_PartialUpdateSemantics(t *testing.T) {
	t.Run("请求体未含 is_top 时保持 nil（表示不改动）", func(t *testing.T) {
		var req UpdateArticleRequest
		require.NoError(t, json.Unmarshal([]byte(`{"title":"新标题"}`), &req))

		assert.Nil(t, req.IsTop, "未提交的 is_top 不应被解析成 false")
		assert.False(t, req.ScheduledAt.Set, "未提交的 scheduled_at 不应被标记为已设置")
	})

	t.Run("显式提交 is_top=false 时可与未提交区分", func(t *testing.T) {
		var req UpdateArticleRequest
		require.NoError(t, json.Unmarshal([]byte(`{"is_top":false}`), &req))

		require.NotNil(t, req.IsTop)
		assert.False(t, *req.IsTop)
	})

	t.Run("显式提交 is_top=true 时正确解析", func(t *testing.T) {
		var req UpdateArticleRequest
		require.NoError(t, json.Unmarshal([]byte(`{"is_top":true}`), &req))

		require.NotNil(t, req.IsTop)
		assert.True(t, *req.IsTop)
	})

	t.Run("显式提交 scheduled_at=null 时标记为已设置且值为空", func(t *testing.T) {
		var req UpdateArticleRequest
		require.NoError(t, json.Unmarshal([]byte(`{"scheduled_at":null}`), &req))

		assert.True(t, req.ScheduledAt.Set)
		assert.Nil(t, req.ScheduledAt.Value)
	})

	t.Run("显式提交 scheduled_at 时间时正确解析", func(t *testing.T) {
		var req UpdateArticleRequest
		require.NoError(t, json.Unmarshal([]byte(`{"scheduled_at":"2026-09-11T10:00:00Z"}`), &req))

		require.True(t, req.ScheduledAt.Set)
		require.NotNil(t, req.ScheduledAt.Value)
		assert.True(t, time.Date(2026, 9, 11, 10, 0, 0, 0, time.UTC).Equal(*req.ScheduledAt.Value))
	})

	t.Run("scheduled_at 非法格式返回错误", func(t *testing.T) {
		var req UpdateArticleRequest
		assert.Error(t, json.Unmarshal([]byte(`{"scheduled_at":"not-a-time"}`), &req))
	})
}
