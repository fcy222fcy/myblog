package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strconv"
	"strings"

	"blog/internal/service"

	"github.com/gin-gonic/gin"
)

// sensitiveKeys 需要脱敏的请求体字段（小写匹配）
var sensitiveKeys = map[string]bool{
	"password":      true,
	"old_password":  true,
	"new_password":  true,
	"token":         true,
	"authorization": true,
	"secret":        true,
	"api_key":       true,
	"apikey":        true,
	"access_token":  true,
	"refresh_token": true,
}

// sanitizeBody 脱敏请求体中的敏感字段，防止明文密码/令牌落库
func sanitizeBody(bodyBytes []byte) []byte {
	var body map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		return bodyBytes
	}
	changed := false
	for key := range body {
		if sensitiveKeys[strings.ToLower(strings.TrimSpace(key))] {
			body[key] = "***"
			changed = true
		}
	}
	if !changed {
		return bodyBytes
	}
	sanitized, err := json.Marshal(body)
	if err != nil {
		return bodyBytes
	}
	return sanitized
}

// Audit 审计日志中间件
// 自动记录所有 POST/PUT/DELETE 操作的审计日志
func Audit(auditLogSvc service.AuditLogService) gin.HandlerFunc {
	// 路径匹配正则：/api/v1/admin/(模块)/(id或动作)
	reTarget := regexp.MustCompile(`/admin/(\w[\w-]*)/?`)

	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.Request.URL.Path

		// 只记录写操作
		action := mapHTTPMethodToAction(method)
		if action == "" {
			c.Next()
			return
		}

		// 只记录 admin 路径下的操作
		if !strings.Contains(path, "/admin/") {
			c.Next()
			return
		}

		// 提取请求信息（不依赖 Auth）
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()
		targetType := mapPathToTargetType(c.Request.URL.Path, reTarget)
		targetID := parseTargetIDFromPath(c)

		// 读取请求体（不消耗原始流）
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		targetTitle := extractTitle(bodyBytes)

		c.Next()

		// Auth 中间件在 c.Next() 之后已设置用户信息
		operatorID := GetUserID(c)
		operatorName := GetUsername(c)

		// 仅记录成功操作（2xx）
		status := c.Writer.Status()
		if status >= 200 && status < 300 {
			// 如果 targetID 仍为 0，尝试从响应中提取（创建操作）
			if targetID == 0 && method == "POST" {
				targetID = extractIDFromResponse(c)
			}

			detail := truncateStr(string(sanitizeBody(bodyBytes)), 500)
			// 异步写入，不阻塞主流程（defer recover 防止 goroutine panic 拖垮进程）
			go func() {
				defer func() { _ = recover() }()
				auditLogSvc.Create(
					operatorID,
					operatorName,
					action,
					targetType,
					targetID,
					targetTitle,
					detail,
					ip,
					userAgent,
				)
			}()
		}
	}
}

// mapHTTPMethodToAction 将 HTTP 方法映射为操作类型
func mapHTTPMethodToAction(method string) string {
	switch method {
	case "POST":
		return "create"
	case "PUT":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return ""
	}
}

// mapPathToTargetType 将 URL 路径映射为目标类型
func mapPathToTargetType(path string, re *regexp.Regexp) string {
	matches := re.FindStringSubmatch(path)
	if len(matches) < 2 {
		return "unknown"
	}
	module := matches[1]

	// 特殊映射
	switch module {
	case "articles":
		return "article"
	case "categories":
		return "category"
	case "tags":
		return "tag"
	case "comments":
		return "comment"
	case "daily-questions":
		return "daily_question"
	case "about":
		return "about_page"
	case "dashboard":
		return "dashboard"
	default:
		return module
	}
}

// parseTargetIDFromPath 从 URL 路径中提取目标 ID
func parseTargetIDFromPath(c *gin.Context) uint {
	// 尝试 :id 参数
	if idStr := c.Param("id"); idStr != "" {
		if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			return uint(id)
		}
	}
	return 0
}

// extractTitle 从请求体中提取标题/名称
func extractTitle(bodyBytes []byte) string {
	if len(bodyBytes) == 0 {
		return ""
	}
	var body map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		return ""
	}
	// 按优先级提取标题
	for _, key := range []string{"title", "name", "nickname"} {
		if v, ok := body[key]; ok {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
		}
	}
	return ""
}

// extractIDFromResponse 从响应中提取 ID（用于创建操作）
func extractIDFromResponse(c *gin.Context) uint {
	// 尝试从 gin context 中获取写入的响应
	// 注意：这依赖于控制器在 c.Set 中存储 ID，或通过其他方式
	if idVal, exists := c.Get("audit_created_id"); exists {
		if id, ok := idVal.(uint); ok {
			return id
		}
	}
	return 0
}

// truncateStr 截断字符串
func truncateStr(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen])
	}
	return s
}
