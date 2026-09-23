package request

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// LoginRequest 登录请求：支持邮箱或用户名 + 密码
type LoginRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

// ValidateLogin 自定义校验：邮箱和用户名至少填一个
func (r *LoginRequest) ValidateLogin() bool {
	hasEmail := strings.TrimSpace(r.Email) != ""
	hasUsername := strings.TrimSpace(r.Username) != ""
	return hasEmail || hasUsername
}

// RegisterRequest 注册请求：邮箱必填
type RegisterRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Nickname string `json:"nickname" binding:"required,max=50"`
	Email    string `json:"email" binding:"required,email"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Avatar   string `json:"avatar" binding:"max=500"`
	Bio      string `json:"bio" binding:"max=500"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// RegisterCustomValidations 注册自定义验证器（在 routes.go 中调用）
func RegisterCustomValidations(v *validator.Validate) {
	// 预留：如需 binding 级别的自定义校验，可在此注册
}
