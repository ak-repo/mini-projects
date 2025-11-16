package router

import (
	"time"

	"github.com/ak-repo/code-board/services/trello-service/internals/app"
	"github.com/ak-repo/code-board/services/trello-service/internals/handler"
	"github.com/ak-repo/code-board/services/trello-service/internals/repo"
	"github.com/ak-repo/code-board/services/trello-service/internals/service"
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
	boardRoutes(api, appCtx)
	listRoutes(api, appCtx)

	return r
}

func boardRoutes(r *gin.RouterGroup, appCtx *app.AppContext) {

	// --------------------------
	// INIT CORE REPOSITORIES
	// --------------------------
	boardRepo := repo.NewBoardRepo(appCtx.DB)

	// --------------------------
	// INIT SERVICES
	// -------------------------
	boardService := service.NewBoardService(boardRepo)

	// --------------------------
	// INIT HANDLERS
	// --------------------------
	boardHandler := handler.NewBoardHandler(boardService)

	board := r.Group("/board")

	// Board
	board.POST("/create", boardHandler.CreateBoard)
	board.GET("/:id", boardHandler.GetBoardByID)

}

func listRoutes(r *gin.RouterGroup, appCtx *app.AppContext) {

	// --------------------------
	// INIT CORE REPOSITORIES
	// --------------------------
	listRepo := repo.NewListRepo(appCtx.DB)

	// --------------------------
	// INIT SERVICES
	// -------------------------
	listService := service.NewListService(listRepo)

	// --------------------------
	// INIT HANDLERS
	// --------------------------
	listHandler := handler.NewListHandler(listService)

	list := r.Group("/list")

	// Board
	list.POST("/create", listHandler.CreateList)
	list.GET("/:id", listHandler.GetListsByBoard)

}
