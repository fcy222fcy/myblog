package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// newBloggerTestRouter 构造「先写 user_id，再走 RequireBlogger」的路由，
// 复现真实链路里 Auth 在前、RequireBlogger 在后的挂载顺序。
// setUserID 为 false 时模拟未挂载 Auth（上下文中没有 user_id）的情况。
func newBloggerTestRouter(setUserID bool, userID uint, bloggerUserID uint) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		if setUserID {
			c.Set("user_id", userID)
		}
		c.Next()
	})
	router.Use(RequireBlogger(bloggerUserID))
	router.GET("/admin/articles", func(c *gin.Context) { c.Status(http.StatusOK) })
	return router
}

func TestRequireBloggerAllowsBlogger(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/articles", nil)

	newBloggerTestRouter(true, 1, 1).ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestRequireBloggerRejectsRegisteredUser(t *testing.T) {
	// 前台注册用户持有合法 token，但身份不是博主
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/articles", nil)

	newBloggerTestRouter(true, 2, 1).ServeHTTP(recorder, req)

	require.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestRequireBloggerRejectsWhenBloggerIDUnset(t *testing.T) {
	// 博主 ID 未配置时必须拒绝。若放行，未登录请求（user_id 同为 0）会被当成博主放进后台。
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/articles", nil)

	newBloggerTestRouter(true, 0, 0).ServeHTTP(recorder, req)

	require.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestRequireBloggerRejectsWhenAuthNotMounted(t *testing.T) {
	// 漏挂 Auth 时上下文没有 user_id，GetUserID 返回 0，同样不能放行
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/articles", nil)

	newBloggerTestRouter(false, 0, 1).ServeHTTP(recorder, req)

	require.Equal(t, http.StatusForbidden, recorder.Code)
}
