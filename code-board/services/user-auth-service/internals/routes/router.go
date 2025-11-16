package router

import (
	"time"

	"github.com/ak-repo/code-board/services/user-auth-service/internals/app"
	"github.com/ak-repo/code-board/services/user-auth-service/internals/handler"
	"github.com/ak-repo/code-board/services/user-auth-service/internals/repo"
	"github.com/ak-repo/code-board/services/user-auth-service/internals/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(appCtx *app.AppContext) *gin.Engine {
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // for local react
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// api
	api := r.Group("/api/v1")

	// auth
	AuthRoutes(appCtx, api)

	// user
	UserRoutes(appCtx, api)

	// admin
	AdminRoutes(appCtx, api)

	return r
}

func AdminRoutes(appCtx *app.AppContext, r *gin.RouterGroup) {
	// Users
	userGroup := r.Group("/users")
	{
		userGroup.GET("/")
		userGroup.GET("/:id")
		userGroup.PUT("/:id")
		userGroup.DELETE("/:id")
	}

}

func UserRoutes(appCtx *app.AppContext, r *gin.RouterGroup) {
	// Password & Email
	passwordGroup := r.Group("/password")
	{
		passwordGroup.POST("/forgot")
		passwordGroup.POST("/reset")
	}
	emailGroup := r.Group("/email")
	{
		emailGroup.POST("/verify")
		emailGroup.GET("/verify/:token")
	}

}

func AuthRoutes(appCtx *app.AppContext, r *gin.RouterGroup) {
	// --------------------------
	// INIT CORE REPOSITORIES
	// --------------------------
	authRepo := repo.NewAuthRepo(appCtx.DB)

	// --------------------------
	// INIT SERVICES
	// -------------------------
	authService := service.NewAuthService(authRepo, appCtx.Config)

	// --------------------------
	// INIT HANDLERS
	// --------------------------
	authHandler := handler.NewAuthHandler(authService)

	// Authentication
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.UserRegister)
		authGroup.POST("/login", authHandler.UserLogin)
		authGroup.POST("/refresh")
		authGroup.POST("/logout")
		authGroup.GET("/me")
	}

}
