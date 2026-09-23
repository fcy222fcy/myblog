package middleware

import (
	"blog/pkg/logger"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequireBlogger 博主身份校验中间件，必须挂在 Auth 之后。
//
// Auth 只校验 token 签名，不校验身份，任何人（包括前台注册用户）拿到合法 token
// 都能通过。所有后台路由组必须同时挂载 Auth 与 RequireBlogger，否则等同于开放后台。
//
// bloggerUserID 取 config.Blogger.UserID。博主登录时签发的 token，其 user_id 恒等于
// 该配置值，因此不会把博主自己挡在门外。
func RequireBlogger(bloggerUserID uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		// UserID 为 0 说明博主账号未配置。此时绝不能放行：未登录请求 GetUserID 也返回 0，
		// 放行会让所有人（包括游客）都以博主身份通过。
		if bloggerUserID == 0 {
			logger.Errorf("博主 UserID 未配置，拒绝访问后台: path=%s", c.FullPath())
			response.Forbidden(c, "博主账号未配置")
			c.Abort()
			return
		}

		userID := GetUserID(c)
		if userID != bloggerUserID {
			logger.Warnf("非博主用户访问后台被拦截: user_id=%d path=%s", userID, c.FullPath())
			response.Forbidden(c, "无权访问")
			c.Abort()
			return
		}

		c.Next()
	}
}
