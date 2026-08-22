package database

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"blog/internal/model/entity"
	"blog/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database 数据库
type Database struct {
	DB *gorm.DB
}

// NewDatabase 创建数据库连接
func NewDatabase(cfg config.MySQLConfig) (*Database, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层 *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 连接池配置
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Database{DB: db}, nil
}

// AutoMigrate 自动迁移
func (d *Database) AutoMigrate() error {
	return d.DB.AutoMigrate(
		&entity.User{},
		&entity.Article{},
		&entity.Category{},
		&entity.Tag{},
		&entity.Comment{},
		&entity.CommentLikeLog{},
		&entity.DailyQuestion{},
		&entity.DailyQuestionLikeLog{},
		&entity.Media{},
		&entity.AboutPage{},
		&entity.VisitLog{},
		&entity.ContentViewEvent{},
		&entity.AuditLog{},
	)
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// SeedIfEmpty 如果数据库为空（articles 表 0 条数据），自动执行 init_data.sql 填充初始数据
// 返回值 seeded=true 表示实际执行了导入
func (d *Database) SeedIfEmpty(initSQLDir string) (bool, error) {
	var articleCount int64
	if err := d.DB.Model(&entity.Article{}).Count(&articleCount).Error; err != nil {
		return false, fmt.Errorf("检查 articles 表数据量失败: %w", err)
	}
	if articleCount > 0 {
		fmt.Printf("[Seed] 检测到 articles 表有 %d 条数据，跳过初始数据导入\n", articleCount)
		return false, nil
	}

	sqlPath := resolveSQLPath(initSQLDir)
	fmt.Printf("[Seed] 数据库为空，准备导入初始数据: %s\n", sqlPath)

	if err := d.execSQLFile(sqlPath); err != nil {
		return false, err
	}
	fmt.Println("[Seed] 初始数据导入成功")
	return true, nil
}

// resolveSQLPath 解析 init_data.sql 的实际路径，兼容开发和部署多种目录结构
func resolveSQLPath(customDir string) string {
	candidates := []string{}
	if customDir != "" {
		candidates = append(candidates, filepath.Join(customDir, "init_data.sql"))
	}
	candidates = append(candidates,
		filepath.Join("scripts", "init_data.sql"),
		filepath.Join("..", "scripts", "init_data.sql"),
		filepath.Join("/app", "scripts", "init_data.sql"),
	)
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			abs, _ := filepath.Abs(p)
			return abs
		}
	}
	return candidates[0]
}

// execSQLFile 读取并执行 SQL 脚本文件，按语句拆分执行，自动跳过注释和空语句
func (d *Database) execSQLFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("打开 SQL 文件失败: %w", err)
	}
	defer f.Close()

	statements := splitSQLStatements(f)
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}

	execCount := 0
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := sqlDB.Exec(stmt); err != nil {
			fmt.Printf("[Seed] 第 %d 条语句执行失败: %v（语句前 120 字: %s）\n",
				i+1, err, truncate(stmt, 120))
			return fmt.Errorf("执行 SQL 失败: %w", err)
		}
		execCount++
	}
	fmt.Printf("[Seed] 共执行 %d 条 SQL 语句\n", execCount)
	return nil
}

func findMatchingBracket(s string, open int) int {
	depth := 0
	inS, inD, inB := false, false, false
	for i := open; i < len(s); i++ {
		c := s[i]
		switch c {
		case '\'':
			if !inD && !inB {
				if inS {
					if i+1 < len(s) && s[i+1] == '\'' {
						i++
						continue
					}
					inS = false
					continue
				}
				inS = true
				continue
			}
		case '"':
			if !inS && !inB {
				if inD {
					if i+1 < len(s) && s[i+1] == '"' {
						i++
						continue
					}
					inD = false
					continue
				}
				inD = true
				continue
			}
		case '`':
			if !inS && !inD {
				inB = !inB
				continue
			}
		case '(':
			if !inS && !inD && !inB {
				depth++
			}
		case ')':
			if !inS && !inD && !inB {
				depth--
				if depth == 0 {
					return i
				}
			}
		case '\\':
			if (inS || inD) && i+1 < len(s) {
				i++
				continue
			}
		}
	}
	return -1
}

// splitSQLStatements 按分号拆分 SQL 语句，正确处理引号、反引号、MySQL ” 转义和注释
func splitSQLStatements(f *os.File) []string {
	var (
		stmts   []string
		cur     strings.Builder
		inS     bool
		inD     bool
		inBack  bool
		inLine  bool
		inBlock bool
	)

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 64*1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		runes := []rune(line)
		i := 0
		for i < len(runes) {
			r := runes[i]

			if inLine {
				i++
				continue
			}

			if inBlock {
				if i > 0 && runes[i-1] == '*' && r == '/' {
					inBlock = false
				}
				cur.WriteRune(r)
				i++
				continue
			}

			// 检测行注释 --
			if !inS && !inD && !inBack && i+1 < len(runes) && r == '-' && runes[i+1] == '-' {
				inLine = true
				i += 2
				continue
			}
			// 检测块注释 /* */
			if !inS && !inD && !inBack && i+1 < len(runes) && r == '/' && runes[i+1] == '*' {
				inBlock = true
				cur.WriteRune(r)
				cur.WriteRune(runes[i+1])
				i += 2
				continue
			}

			switch r {
			case '\'':
				if !inD && !inBack {
					cur.WriteRune(r)
					if inS {
						// MySQL 标准 '' 转义为单引号字符
						if i+1 < len(runes) && runes[i+1] == '\'' {
							cur.WriteRune('\'')
							i += 2
							continue
						}
						inS = false
						i++
						continue
					}
					inS = true
					i++
					continue
				}
			case '"':
				if !inS && !inBack {
					cur.WriteRune(r)
					if inD {
						if i+1 < len(runes) && runes[i+1] == '"' {
							cur.WriteRune('"')
							i += 2
							continue
						}
						inD = false
						i++
						continue
					}
					inD = true
					i++
					continue
				}
			case '`':
				if !inS && !inD {
					cur.WriteRune(r)
					inBack = !inBack
					i++
					continue
				}
			case ';':
				if !inS && !inD && !inBack {
					stmts = append(stmts, cur.String())
					cur.Reset()
					i++
					continue
				}
			case '\\':
				// 保留转义字符，避免改变传给 MySQL 的原始 SQL。
				if (inS || inD) && i+1 < len(runes) {
					cur.WriteRune(r)
					cur.WriteRune(runes[i+1])
					i += 2
					continue
				}
			}

			cur.WriteRune(r)
			i++
		}

		inLine = false
		cur.WriteByte('\n')
	}
	if cur.Len() > 0 {
		stmts = append(stmts, cur.String())
	}
	return stmts
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
