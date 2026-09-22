package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	bizcrypt "blog/pkg/bcrypt"
	"blog/pkg/config"
	bizerrors "blog/pkg/errors"
	"blog/pkg/jwt"
	"blog/pkg/logger"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// validateInput 校验输入（仅做基础长度限制）
// 注：SQL 注入已由 GORM 参数化查询防护，这里不再使用黑名单拦截（黑名单会误伤合法输入）
func validateInput(input string) bool {
	return len(strings.TrimSpace(input)) <= 200
}

// authService 认证服务实现
type authService struct {
	userRepo repository.UserRepository
	jwt      *jwt.JWT
	config   *config.Config
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo repository.UserRepository, jwt *jwt.JWT, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		jwt:      jwt,
		config:   cfg,
	}
}

// isBlogger 判断 userID 是否为博主虚拟账号
func (s *authService) isBlogger(userID uint) bool {
	if s.config == nil {
		return false
	}
	return userID == s.config.Blogger.UserID
}

// isBloggerLogin 判断用户名密码是否匹配博主账号
func (s *authService) isBloggerLogin(username, password string) bool {
	if s.config == nil {
		return false
	}
	b := s.config.Blogger
	return strings.EqualFold(strings.TrimSpace(username), strings.TrimSpace(b.Username)) &&
		strings.TrimSpace(password) == strings.TrimSpace(b.Password)
}

// Login 用户登录
func (s *authService) Login(req *request.LoginRequest) (*response.LoginResponse, error) {
	logger.Infof("用户登录尝试, username=%s, email=%s", req.Username, req.Email)

	// 危险字符校验（密码一律不做危险字符校验，防止误伤合法密码）
	if req.Username != "" && !validateInput(req.Username) {
		return nil, bizerrors.New(bizerrors.CodeInvalidParams, "输入包含非法字符")
	}
	if req.Email != "" && !validateInput(req.Email) {
		return nil, bizerrors.New(bizerrors.CodeInvalidParams, "输入包含非法字符")
	}

	// 1. 优先匹配博主硬编码账号（用 username 匹配；若用户只填 email 且与博主邮箱相同，也兼容）
	bloggerByUsername := strings.TrimSpace(req.Username) != "" && s.isBloggerLogin(req.Username, req.Password)
	bloggerByEmail := false
	if !bloggerByUsername && strings.TrimSpace(req.Email) != "" && s.config != nil {
		b := s.config.Blogger
		if strings.EqualFold(strings.TrimSpace(req.Email), strings.TrimSpace(b.Email)) &&
			strings.TrimSpace(req.Password) == strings.TrimSpace(b.Password) {
			bloggerByEmail = true
		}
	}
	if bloggerByUsername || bloggerByEmail {
		b := s.config.Blogger
		token, expiresAt, err := s.jwt.GenerateToken(b.UserID, b.Username)
		if err != nil {
			logger.Error("生成博主Token失败", zap.Error(err))
			return nil, fmt.Errorf("生成Token失败, %w", err)
		}
		nickname := ""
		avatar := ""
		email := ""
		if u, dbErr := s.userRepo.FindByID(b.UserID); dbErr == nil && u != nil {
			nickname = u.Nickname
			avatar = u.Avatar
			email = u.Email
		}
		if nickname == "" {
			nickname = b.Nickname
		}
		if avatar == "" {
			avatar = b.Avatar
		}
		if email == "" {
			email = b.Email
		}
		logger.Infof("博主登录成功, username: %s", b.Username)
		return &response.LoginResponse{
			Token:     token,
			ExpiresAt: expiresAt,
			User: response.UserInfo{
				ID:       b.UserID,
				Username: b.Username,
				Nickname: nickname,
				Avatar:   avatar,
				Email:    email,
			},
		}, nil
	}

	// 2. 非博主账号：优先按邮箱查，其次按用户名查
	var user *entity.User
	var err error
	loginID := ""
	if strings.TrimSpace(req.Email) != "" {
		loginID = strings.TrimSpace(req.Email)
		user, err = s.userRepo.FindByEmail(loginID)
		if err != nil {
			logger.Error("按邮箱查询用户失败", zap.Error(err))
			return nil, fmt.Errorf("查询用户失败, %w", err)
		}
	}
	if user == nil && strings.TrimSpace(req.Username) != "" {
		loginID = strings.TrimSpace(req.Username)
		user, err = s.userRepo.FindByUsername(loginID)
		if err != nil {
			logger.Error("按用户名查询用户失败", zap.Error(err))
			return nil, fmt.Errorf("查询用户失败, %w", err)
		}
	}
	if user == nil {
		logger.Warn("用户不存在", zap.String("loginID", loginID))
		return nil, bizerrors.New(bizerrors.CodeUserNotFound, "账号或密码错误")
	}

	// 验证密码
	if !bizcrypt.CheckPassword(req.Password, user.Password) {
		return nil, bizerrors.New(bizerrors.CodePasswordIncorrect, "账号或密码错误")
	}

	// 检查账号状态
	if user.Status != 1 {
		return nil, bizerrors.New(bizerrors.CodeAccountDisabled, bizerrors.GetMessage(bizerrors.CodeAccountDisabled))
	}

	// 生成 Token
	token, expiresAt, err := s.jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		logger.Error("生成Token失败", zap.Error(err))
		return nil, fmt.Errorf("生成Token失败, %w", err)
	}

	logger.Infof("用户登录成功, username=%s, email=%s", user.Username, user.Email)
	return &response.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: response.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Email:    user.Email,
		},
	}, nil
}

// GetProfile 获取用户信息
func (s *authService) GetProfile(userID uint) (*response.UserProfileResponse, error) {
	logger.Infof("获取用户信息, userID: %d", userID)

	// 博主账号：与主页 /api/v1/user/info 保持同源，从 user 表读取展示信息
	if s.isBlogger(userID) {
		b := s.config.Blogger
		nickname := b.Nickname
		avatar := b.Avatar
		email := b.Email
		bio := ""
		if u, dbErr := s.userRepo.FindByID(userID); dbErr == nil && u != nil {
			if u.Nickname != "" {
				nickname = u.Nickname
			}
			if u.Avatar != "" {
				avatar = u.Avatar
			}
			if u.Email != "" {
				email = u.Email
			}
			if u.Bio != "" {
				bio = u.Bio
			}
		}
		return &response.UserProfileResponse{
			ID:       b.UserID,
			Username: b.Username,
			Nickname: nickname,
			Email:    email,
			Avatar:   avatar,
			Bio:      bio,
		}, nil
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.Error("查询用户失败", zap.Error(err))
		return nil, fmt.Errorf("查询用户失败, %w", err)
	}
	if user == nil {
		logger.Warn("用户不存在", zap.Uint("userID", userID))
		return nil, bizerrors.New(bizerrors.CodeUserNotFound, bizerrors.GetMessage(bizerrors.CodeUserNotFound))
	}

	return &response.UserProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
	}, nil
}

// ChangePassword 修改密码
func (s *authService) ChangePassword(userID uint, req *request.ChangePasswordRequest) error {
	logger.Infof("修改密码, userID: %d", userID)

	// 博主账号：校验旧密码（明文），然后更新配置中的博主密码
	if s.isBlogger(userID) {
		b := &s.config.Blogger
		if strings.TrimSpace(req.OldPassword) != strings.TrimSpace(b.Password) {
			return bizerrors.New(bizerrors.CodePasswordIncorrect, bizerrors.GetMessage(bizerrors.CodePasswordIncorrect))
		}
		b.Password = strings.TrimSpace(req.NewPassword)
		logger.Infof("博主密码修改成功, userID: %d", userID)
		return nil
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.Error("查询用户失败", zap.Error(err))
		return fmt.Errorf("查询用户失败, %w", err)
	}
	if user == nil {
		logger.Warn("用户不存在", zap.Uint("userID", userID))
		return bizerrors.New(bizerrors.CodeUserNotFound, bizerrors.GetMessage(bizerrors.CodeUserNotFound))
	}

	// 验证旧密码
	if !bizcrypt.CheckPassword(req.OldPassword, user.Password) {
		return bizerrors.New(bizerrors.CodePasswordIncorrect, bizerrors.GetMessage(bizerrors.CodePasswordIncorrect))
	}

	// 加密新密码
	hashedPassword, err := bizcrypt.HashPassword(req.NewPassword)
	if err != nil {
		return bizerrors.New(bizerrors.CodeInternalServer, "密码加密失败")
	}

	user.Password = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		logger.Error("更新密码失败", zap.Error(err))
		return fmt.Errorf("更新密码失败, %w", err)
	}

	logger.Infof("修改密码成功, userID: %d", userID)
	return nil
}
