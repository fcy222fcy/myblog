package api

import (
	"blog/internal/api/v1/about_page"
	"blog/internal/api/v1/article"
	"blog/internal/api/v1/audit_log"
	"blog/internal/api/v1/auth"
	"blog/internal/api/v1/category"
	"blog/internal/api/v1/comment"
	"blog/internal/api/v1/content_view"
	"blog/internal/api/v1/daily_question"
	"blog/internal/api/v1/media"
	"blog/internal/api/v1/rss"
	"blog/internal/api/v1/sitemap"
	"blog/internal/api/v1/tag"
	"blog/internal/api/v1/user"
	"blog/internal/middleware"
	"blog/internal/model/dto/request"
	"blog/internal/repository"
	"blog/internal/service"
	"blog/pkg/config"
	blogjwt "blog/pkg/jwt"
	"blog/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// Router 路由器
type Router struct {
	engine *gin.Engine
	config *config.Config

	// JWT 实例（预创建，复用）
	jwtInstance *blogjwt.JWT

	// 各模块 controller
	authController          *auth.Controller
	userController          *user.Controller
	articleController       *article.Controller
	categoryController      *category.Controller
	tagController           *tag.Controller
	commentController       *comment.Controller
	dailyQuestionController *daily_question.Controller
	contentViewController   *content_view.Controller
	aboutPageController     *about_page.Controller
	mediaController         *media.Controller
	auditLogController      *audit_log.Controller
	rssHandler              *rss.Handler
	sitemapHandler          *sitemap.Handler

	// 仓库（用于仪表盘统计）
	articleRepo    repository.ArticleRepository
	commentRepo    repository.CommentRepository
	contentViewSvc service.ContentViewService

	// 审计日志服务（供中间件使用）
	auditLogSvc service.AuditLogService
}

// NewRouter 创建路由器
func NewRouter(
	authSvc service.AuthService,
	userSvc service.UserService,
	articleSvc service.ArticleService,
	categorySvc service.CategoryService,
	tagSvc service.TagService,
	commentSvc service.CommentService,
	dailyQuestionSvc service.DailyQuestionService,
	contentViewSvc service.ContentViewService,
	aboutPageSvc service.AboutPageService,
	auditLogSvc service.AuditLogService,
	articleRepo repository.ArticleRepository,
	commentRepo repository.CommentRepository,
	config *config.Config,
) *Router {
	engine := gin.Default()
	if err := configureTrustedProxies(engine, config.Server.TrustedProxies); err != nil {
		logger.Warn("可信代理配置无效，已禁用代理请求头解析", zap.Error(err))
		_ = engine.SetTrustedProxies(nil)
	}

	// 注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("datetime", request.ValidateDate)
	}

	// 预创建 JWT 实例，避免每次请求都创建
	jwtInstance := blogjwt.NewJWT(config.JWT)

	return &Router{
		engine:                  engine,
		config:                  config,
		jwtInstance:             jwtInstance,
		authController:          auth.NewController(authSvc),
		userController:          user.NewController(userSvc),
		articleController:       article.NewController(articleSvc),
		categoryController:      category.NewController(categorySvc),
		tagController:           tag.NewController(tagSvc),
		commentController:       comment.NewController(commentSvc),
		dailyQuestionController: daily_question.NewController(dailyQuestionSvc),
		contentViewController:   content_view.NewController(contentViewSvc, config.JWT.Secret),
		aboutPageController:     about_page.NewController(aboutPageSvc),
		mediaController:         media.NewController(config.App.UploadDir),
		auditLogController:      audit_log.NewController(auditLogSvc),
		rssHandler:              rss.NewHandler(articleRepo),
		sitemapHandler:          sitemap.NewHandler(articleRepo),
		articleRepo:             articleRepo,
		commentRepo:             commentRepo,
		contentViewSvc:          contentViewSvc,
		auditLogSvc:             auditLogSvc,
	}
}

func configureTrustedProxies(engine *gin.Engine, trustedProxies []string) error {
	return engine.SetTrustedProxies(trustedProxies)
}

// Setup 设置路由
func (r *Router) Setup() *gin.Engine {
	// 全局中间件
	r.engine.Use(middleware.Recovery())
	r.engine.Use(middleware.Logger())

	// 使用配置的 CORS 来源（无配置时默认本地开发端口）
	origins := r.config.Server.CORSOrigins
	if len(origins) == 0 {
		origins = []string{
			"http://localhost:3000",
			"http://localhost:9090",
			"http://localhost:5173",
			"http://localhost:5174",
		}
	}
	r.engine.Use(middleware.CORS(origins...))

	// API v1 路由组
	apiV1 := r.engine.Group("/api/v1")

	// 审计日志中间件（仅记录 admin 写操作）
	apiV1.Use(middleware.Audit(r.auditLogSvc))

	// 注册各模块路由（传入预创建的 JWT 实例）
	auth.RegisterRoutes(apiV1, r.authController, r.jwtInstance)
	article.RegisterRoutes(apiV1, r.articleController, r.jwtInstance)
	category.RegisterRoutes(apiV1, r.categoryController, r.jwtInstance)
	tag.RegisterRoutes(apiV1, r.tagController, r.jwtInstance)
	comment.RegisterRoutes(apiV1, r.commentController, r.jwtInstance)
	daily_question.RegisterRoutes(apiV1, r.dailyQuestionController, r.jwtInstance)
	content_view.RegisterRoutes(apiV1, r.contentViewController)
	user.RegisterRoutes(apiV1, r.userController, r.jwtInstance)
	media.RegisterRoutes(apiV1, r.mediaController, r.jwtInstance)
	audit_log.RegisterRoutes(apiV1, r.auditLogController, r.jwtInstance)
	rss.RegisterRoutes(apiV1, r.rssHandler)
	sitemap.RegisterRoutes(apiV1, r.sitemapHandler)

	// 注册仪表盘路由（需要登录）
	r.registerDashboardRoutes(apiV1)

	// 注册关于页面路由
	r.registerAboutPageRoutes(apiV1)

	// 静态文件服务 - 上传的文件
	r.engine.Static("/uploads", r.mediaController.UploadDir())

	return r.engine
}

// registerDashboardRoutes 注册仪表盘路由
func (r *Router) registerDashboardRoutes(rg *gin.RouterGroup) {
	// 需要登录的路由
	protected := rg.Group("")
	protected.Use(middleware.Auth(r.jwtInstance))
	{
		dashboard := protected.Group("/admin/dashboard")
		{
			dashboard.GET("/stats", r.getDashboardStats)
			dashboard.GET("/recent-articles", r.getRecentArticles)
		}
	}
}

// registerAboutPageRoutes 注册关于页面路由
func (r *Router) registerAboutPageRoutes(rg *gin.RouterGroup) {
	// 公开路由（无需登录）
	about := rg.Group("/about")
	{
		about.GET("", r.aboutPageController.GetAboutPage)
	}

	// 需要登录的路由
	protected := rg.Group("")
	protected.Use(middleware.Auth(r.jwtInstance))
	{
		protected.PUT("/admin/about", r.aboutPageController.UpdateAboutPage)
	}
}

// getDashboardStats 获取仪表盘统计
func (r *Router) getDashboardStats(c *gin.Context) {
	// 统计文章数量
	articleCount, err := r.articleRepo.Count("")
	if err != nil {
		logger.Warn("统计文章数量失败", zap.Error(err))
	}

	// 统计已发布文章数量
	publishedCount, err := r.articleRepo.Count("published")
	if err != nil {
		logger.Warn("统计已发布文章数量失败", zap.Error(err))
	}

	// 统计定时发布文章数量
	scheduledCount, err := r.articleRepo.Count("scheduled")
	if err != nil {
		logger.Warn("统计定时发布文章数量失败", zap.Error(err))
	}

	// 统计总浏览量
	totalViews, err := r.articleRepo.SumViewCount()
	if err != nil {
		logger.Warn("统计总浏览量失败", zap.Error(err))
	}

	// 统计今日浏览量
	todayViews, err := r.contentViewSvc.CountToday(repository.ContentTypeArticle)
	if err != nil {
		logger.Warn("统计今日浏览量失败", zap.Error(err))
	}

	// 统计评论数量
	commentCount, err := r.commentRepo.Count("")
	if err != nil {
		logger.Warn("统计评论数量失败", zap.Error(err))
	}

	// 统计待审核评论数量
	pendingCommentCount, err := r.commentRepo.Count("pending")
	if err != nil {
		logger.Warn("统计待审核评论数量失败", zap.Error(err))
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"article_count":   articleCount,
			"published_count": publishedCount,
			"draft_count":     articleCount - publishedCount - scheduledCount,
			"scheduled_count": scheduledCount,
			"total_views":     totalViews,
			"today_views":     todayViews,
			"comment_count":   commentCount,
			"pending_count":   pendingCommentCount,
		},
	})
}

// getRecentArticles 获取最近文章
func (r *Router) getRecentArticles(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 20 {
		limit = 5
	}

	articles, err := r.articleRepo.GetRecent(limit)
	if err != nil {
		logger.Warn("获取最近文章失败", zap.Error(err))
		c.JSON(500, gin.H{
			"code":    500,
			"message": "获取最近文章失败",
		})
		return
	}

	// 转换为简化响应格式
	type recentArticle struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Summary   string `json:"summary"`
		Cover     string `json:"cover"`
		ViewCount int64  `json:"view_count"`
		CreatedAt string `json:"created_at"`
	}

	var result []recentArticle
	for _, article := range articles {
		result = append(result, recentArticle{
			ID:        article.ID,
			Title:     article.Title,
			Summary:   article.Summary,
			Cover:     article.Cover,
			ViewCount: article.ViewCount,
			CreatedAt: article.CreatedAt.Format("2006-01-02"),
		})
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data":    result,
	})
}
