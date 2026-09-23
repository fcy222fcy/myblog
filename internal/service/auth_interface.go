package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
)

// AuthService 认证服务接口
type AuthService interface {
	// Login 用户登录
	Login(req *request.LoginRequest) (*response.LoginResponse, error)

	// Register 用户注册
	Register(req *request.RegisterRequest) error

	// GetProfile 获取用户信息
	GetProfile(userID uint) (*response.UserProfileResponse, error)

	// ChangePassword 修改密码
	ChangePassword(userID uint, req *request.ChangePasswordRequest) error
}
