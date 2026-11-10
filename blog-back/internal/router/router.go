package router

import (
	"database/sql"
	"os"
	"path/filepath"

	"blog-back/internal/config"
	"blog-back/internal/handler"
	"blog-back/internal/middleware"
	"blog-back/internal/store"
	"blog-back/internal/upload"
	"blog-back/internal/version"

	"github.com/gin-gonic/gin"
)

func Setup(cfg config.Config, db *sql.DB) *gin.Engine {
	tagStore := store.NewTagStore(db)
	articleH := handler.NewArticleHandler(store.NewArticleStore(db, tagStore))
	tagH := handler.NewTagHandler(tagStore)
	userH := handler.NewUserHandler(store.NewUserStore(db))
	hollowH := handler.NewHollowHandler(store.NewHollowStore(db))
	uploadH := handler.NewUploadHandler(upload.NewService(cfg.UploadPath, cfg.UploadURLPrefix))
	chatH := handler.NewChatHandler(cfg)

	r := gin.New()
	proxies := cfg.TrustedProxies
	if proxies == nil {
		proxies = []string{"127.0.0.1", "::1"}
	}
	if err := r.SetTrustedProxies(proxies); err != nil {
		panic("invalid TRUSTED_PROXIES configuration")
	}

	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.CORSOrigins))

	imagesDir := filepath.Join(cfg.UploadPath, "images")
	r.Static("/upload/images", imagesDir)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "ok",
			"version":  version.Info,
			"features": []string{"chat", "article", "upload"},
		})
	})

	r.GET("/site/summary", handler.SiteSummary(store.NewSummaryStore(db)))
	registerAPI(r, articleH, tagH, userH, hollowH, uploadH, chatH, cfg)
	return r
}

func registerAPI(
	r *gin.Engine,
	articleH *handler.ArticleHandler,
	tagH *handler.TagHandler,
	userH *handler.UserHandler,
	hollowH *handler.HollowHandler,
	uploadH *handler.UploadHandler,
	chatH *handler.ChatHandler,
	cfg config.Config,
) {
	// 公开读取接口
	r.GET("/article/list", articleH.List)
	r.GET("/article/user/:userId", articleH.ByUser)
	r.GET("/article/tag/:tagId", articleH.ByTag)
	r.GET("/article/:id", articleH.GetByID)
	r.POST("/article/view/:id", articleH.View)

	r.GET("/tag", tagH.List)
	r.GET("/tag/:id", tagH.GetByID)

	r.GET("/user", userH.List)
	r.GET("/user/:id", userH.GetByID)

	r.GET("/hollow", hollowH.List)

	// 需要管理员权限的写操作
	sessions := middleware.NewAdminSessions(os.Getenv("ADMIN_TOKEN"))
	r.POST("/auth/login", sessions.Login)
	r.POST("/auth/logout", sessions.Logout)
	r.GET("/auth/session", sessions.Require(), func(c *gin.Context) { c.JSON(200, gin.H{"code": 200}) })
	admin := r.Group("/", sessions.Require())
	{
		admin.POST("/article", articleH.Save)
		admin.POST("/article/save", articleH.Save)
		admin.PUT("/article/update", articleH.Update)
		admin.DELETE("/article/delete/:id", articleH.Delete)

		admin.POST("/tag", tagH.Save)
		admin.PUT("/tag", tagH.Update)
		admin.DELETE("/tag/:id", tagH.Delete)

		admin.POST("/user/create", userH.Save)
		admin.PUT("/user/update", userH.Update)
		admin.DELETE("/user/delete/:id", userH.Delete)

		admin.POST("/hollow", hollowH.Save)

		admin.POST("/upload/image", uploadH.Image)
	}
	// 兼容 Nginx 是否剥离 /api 前缀的两种反代方式
	r.GET("/chat/status", chatH.Status)
	r.GET("/api/chat/status", chatH.Status)
	r.POST("/chat", chatH.Chat)
	r.POST("/api/chat", chatH.Chat)
}
