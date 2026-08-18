package errors

// 错误码定义
// 0       成功
// 400-499 通用客户端错误
// 500-599 服务端错误
// 1000+   用户/认证业务错误
// 2000+   参数错误
// 3000+   资源错误

const (
	// 成功
	CodeSuccess = 0

	// 通用错误
	CodeBadRequest     = 400
	CodeUnauthorized   = 401
	CodeForbidden      = 403
	CodeNotFound       = 404
	CodeInternalServer = 500

	// 用户/认证业务错误 (1000+)
	CodeUserNotFound      = 1001
	CodeUserAlreadyExists = 1002
	CodePasswordIncorrect = 1003
	CodeTokenExpired      = 1004
	CodeTokenInvalid      = 1005
	CodeAccountDisabled   = 1006

	// 参数错误 (2000+)
	CodeInvalidParams    = 2001
	CodeMissingRequired  = 2002
	CodeInvalidEmail     = 2003
	CodeInvalidPhone     = 2004
	CodePasswordTooShort = 2005
	CodeUsernameTooShort = 2006

	// 资源错误 (3000+)
	CodeArticleNotFound       = 3001
	CodeCategoryNotFound      = 3002
	CodeTagNotFound           = 3003
	CodeCommentNotFound       = 3004
	CodeDailyQuestionNotFound = 3006
	CodeMediaNotFound         = 3007
	CodePageNotFound          = 3008
	CodeRecordNotFound        = 3009

	// 业务错误 (4000+)
	CodeCategoryNameExists      = 4001
	CodeCategoryHasArticles     = 4002
	CodeTagNameExists           = 4003
	CodeTagHasArticles          = 4004
	CodeDailyQuestionDateExists = 4005
	CodeCommentAlreadyLiked     = 4006
	CodeCommentNotLiked         = 4007
	CodeAlreadyLiked            = 4008
)

// 错误码消息映射
var CodeMessages = map[int]string{
	CodeSuccess:               "成功",
	CodeBadRequest:            "请求错误",
	CodeUnauthorized:          "未授权",
	CodeForbidden:             "禁止访问",
	CodeNotFound:              "资源不存在",
	CodeInternalServer:        "服务器内部错误",
	CodeUserNotFound:          "用户不存在",
	CodeUserAlreadyExists:     "用户已存在",
	CodePasswordIncorrect:     "密码错误",
	CodeTokenExpired:          "Token 已过期",
	CodeTokenInvalid:          "Token 无效",
	CodeAccountDisabled:       "账户已禁用",
	CodeInvalidParams:         "参数错误",
	CodeMissingRequired:       "缺少必填参数",
	CodeInvalidEmail:          "邮箱格式错误",
	CodeInvalidPhone:          "手机号格式错误",
	CodePasswordTooShort:      "密码长度不足",
	CodeUsernameTooShort:      "用户名长度不足",
	CodeArticleNotFound:       "文章不存在",
	CodeCategoryNotFound:      "分类不存在",
	CodeTagNotFound:           "标签不存在",
	CodeCommentNotFound:       "评论不存在",
	CodeDailyQuestionNotFound: "每日一问不存在",
	CodeMediaNotFound:         "媒体文件不存在",
	CodePageNotFound:          "页面不存在",
	CodeRecordNotFound:        "记录不存在",
	CodeCommentAlreadyLiked:   "您已经点赞过了",
	CodeCommentNotLiked:       "您还没有点赞过",
	CodeAlreadyLiked:          "您已经点赞过了",
}

// GetMessage 获取错误码对应的消息
func GetMessage(code int) string {
	if msg, ok := CodeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
