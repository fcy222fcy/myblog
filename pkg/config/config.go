package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	MySQL   MySQLConfig   `mapstructure:"mysql"`
	Redis   RedisConfig   `mapstructure:"redis"`
	JWT     JWTConfig     `mapstructure:"jwt"`
	Log     LogConfig     `mapstructure:"log"`
	Email   EmailConfig   `mapstructure:"email"`
	Blogger BloggerConfig `mapstructure:"blogger"`
	App     AppConfig     `mapstructure:"app"`
}

// AppConfig 应用通用配置
type AppConfig struct {
	UploadDir  string `mapstructure:"upload_dir"`
	SeedData   bool   `mapstructure:"seed_data"`    // 启动时若库为空，自动填充初始数据
	InitSQLDir string `mapstructure:"init_sql_dir"` // init_data.sql 所在目录，默认 scripts
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port        int      `mapstructure:"port"`
	CORSOrigins []string `mapstructure:"cors_origins"` // 允许的跨域来源，逗号分隔
}

// MySQLConfig 数据库配置
type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

// DSN 获取数据库连接字符串，默认启用 multiStatements 以便批量执行 SQL
func (m *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		m.Username, m.Password, m.Host, m.Port, m.DBName)
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireHour int    `mapstructure:"expire_hour"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	From      string `mapstructure:"from"`
	FromEmail string `mapstructure:"from_email"`
}

// BloggerConfig 博主账号配置（写死配置，不走用户表）
type BloggerConfig struct {
	UserID   uint   `mapstructure:"user_id"`  // 博主虚拟 UserID，固定值用于判断博主身份
	Username string `mapstructure:"username"` // 登录用户名
	Password string `mapstructure:"password"` // 登录密码（明文，注意生产环境应加密或环境变量覆盖）
	Nickname string `mapstructure:"nickname"` // 显示昵称
	Email    string `mapstructure:"email"`    // 邮箱（用于头像/通知，与 from_email 保持一致）
	Avatar   string `mapstructure:"avatar"`   // 头像URL（可选，空则用邮箱Gravatar）
}

// Load 加载配置
func Load() (*Config, error) {
	loadDotEnv(".env")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.AddConfigPath(".")

	// 设置默认值（未配置时的兜底）
	viper.SetDefault("app.seed_data", true)
	viper.SetDefault("app.init_sql_dir", "")
	viper.SetDefault("app.upload_dir", "uploads")
	viper.SetDefault("server.port", 9090)
	viper.SetDefault("server.cors_origins", []string{
		"http://localhost:3000",
		"http://localhost:9090",
		"http://localhost:5173",
		"http://localhost:5174",
	})

	// 读取环境变量
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 环境变量覆盖敏感配置
	applyEnvOverrides(&config)

	return &config, nil
}

func applyEnvOverrides(config *Config) {
	overrideInt(&config.Server.Port, "SERVER_PORT", "BACKEND_PORT")
	// CORS 来源：CORS_ORIGINS 逗号分隔
	if value, ok := firstEnv("CORS_ORIGINS"); ok {
		if origins := parseCSV(value); len(origins) > 0 {
			config.Server.CORSOrigins = origins
		}
	}

	overrideString(&config.MySQL.Host, "DB_HOST", "MYSQL_HOST")
	overrideInt(&config.MySQL.Port, "DB_PORT", "MYSQL_PORT")
	overrideString(&config.MySQL.Username, "DB_USER", "MYSQL_USER", "MYSQL_USERNAME")
	overrideString(&config.MySQL.Password, "DB_PASSWORD", "MYSQL_PASSWORD")
	overrideString(&config.MySQL.DBName, "DB_NAME", "MYSQL_DATABASE", "MYSQL_DBNAME")

	overrideString(&config.Redis.Host, "REDIS_HOST")
	overrideInt(&config.Redis.Port, "REDIS_PORT")
	overrideString(&config.Redis.Password, "REDIS_PASSWORD")
	overrideInt(&config.Redis.DB, "REDIS_DB")

	overrideString(&config.JWT.Secret, "JWT_SECRET")
	overrideInt(&config.JWT.ExpireHour, "JWT_EXPIRE_HOURS", "JWT_EXPIRE_HOUR")

	overrideString(&config.Log.Level, "LOG_LEVEL")
	overrideString(&config.Log.Filename, "LOG_FILENAME")
	overrideInt(&config.Log.MaxSize, "LOG_MAX_SIZE")
	overrideInt(&config.Log.MaxBackups, "LOG_MAX_BACKUPS")
	overrideInt(&config.Log.MaxAge, "LOG_MAX_AGE")

	overrideString(&config.Email.Host, "EMAIL_HOST")
	overrideInt(&config.Email.Port, "EMAIL_PORT")
	overrideString(&config.Email.Username, "EMAIL_USER", "EMAIL_USERNAME")
	overrideString(&config.Email.Password, "EMAIL_PASSWORD")
	overrideString(&config.Email.From, "EMAIL_FROM")
	overrideString(&config.Email.FromEmail, "EMAIL_FROM_EMAIL")

	overrideUint(&config.Blogger.UserID, "BLOGGER_USER_ID")
	overrideString(&config.Blogger.Username, "BLOGGER_USERNAME")
	overrideString(&config.Blogger.Password, "BLOGGER_PASSWORD")
	overrideString(&config.Blogger.Nickname, "BLOGGER_NICKNAME")
	overrideString(&config.Blogger.Email, "BLOGGER_EMAIL")
	overrideString(&config.Blogger.Avatar, "BLOGGER_AVATAR")

	overrideBool(&config.App.SeedData, "APP_SEED_DATA")
	overrideString(&config.App.InitSQLDir, "APP_INIT_SQL_DIR")
	overrideString(&config.App.UploadDir, "APP_UPLOAD_DIR", "UPLOAD_DIR")
}

func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key, value, ok := parseDotEnvLine(scanner.Text())
		if !ok {
			continue
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		_ = os.Setenv(key, value)
	}
}

func parseDotEnvLine(line string) (string, string, bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}
	line = strings.TrimPrefix(line, "export ")
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if key == "" {
		return "", "", false
	}
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}
	}
	return key, value, true
}

func overrideString(target *string, names ...string) {
	if value, ok := firstEnv(names...); ok {
		*target = value
	}
}

func overrideInt(target *int, names ...string) {
	if value, ok := firstEnv(names...); ok {
		parsed, err := strconv.Atoi(value)
		if err == nil {
			*target = parsed
		}
	}
}

func overrideUint(target *uint, names ...string) {
	if value, ok := firstEnv(names...); ok {
		parsed, err := strconv.ParseUint(value, 10, 0)
		if err == nil {
			*target = uint(parsed)
		}
	}
}

func overrideBool(target *bool, names ...string) {
	if value, ok := firstEnv(names...); ok {
		parsed, err := strconv.ParseBool(value)
		if err == nil {
			*target = parsed
		}
	}
}

func firstEnv(names ...string) (string, bool) {
	for _, name := range names {
		if value, ok := os.LookupEnv(name); ok {
			return value, true
		}
	}
	return "", false
}

// parseCSV 解析逗号分隔的配置值，去除空白与空项
func parseCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			result = append(result, t)
		}
	}
	return result
}
