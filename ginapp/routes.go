package ginapp

import (
	"app/repository"
	"app/shared"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// contract.Serverを実装したサーバーの構造体
type GinServer struct {
	engine *gin.Engine
}

// サーバーのセットアップ
func setupRoutes(router *gin.Engine) {
	// api ルートのグループ
	v1 := router.Group("/api/v1")
	{
		// クエリパラメーターによる認証
		v1.Use(tokenAuth)

		// todo ルートのグループ
		todo := v1.Group("/todos")
		{
			// 一覧取得
			todo.GET("", getTodos)
			// 詳細取得
			todo.GET("/:id", getTodoById)
			// 作成
			todo.POST("", createTodo)
		}
	}
}

func tokenAuth(ctx *gin.Context) {
	token := ctx.Query("token")

	// トークンが一致しない場合は 401 を返す
	if shared.IsInvalidToken(token) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})

		ctx.Abort()
	}

	ctx.Next()
}

func getTodos(ctx *gin.Context) {
	todos := repository.GetTodos()

	ctx.JSON(http.StatusOK, gin.H{
		"todos": todos,
	})
}

func getTodoById(ctx *gin.Context) {
	// パスパラメーターから id を取得
	id := ctx.Param("id")

	todo, ok := repository.GetTodoById(id)

	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"todo": todo,
	})
}

func createTodo(ctx *gin.Context) {
	todo := &repository.TodoForCreate{}

	// リクエストボディをパース
	ctx.BindJSON(todo)

	created := repository.CreateTodo(*todo)

	ctx.JSON(http.StatusCreated, gin.H{
		"todo": created,
	})
}

// サーバーのインスタンスを生成
func NewGinServer() shared.Server {
	engine := gin.Default()

	setupRoutes(engine)

	return &GinServer{
		engine,
	}
}

// サーバーを起動
func (s *GinServer) Run(addr string) {
	log.Fatal(s.engine.Run(addr))
}
