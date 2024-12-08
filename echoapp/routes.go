package echoapp

import (
	"app/repository"
	"app/shared"
	"net/http"

	"github.com/labstack/echo/v4"
)

// contract.Serverを実装したサーバーの構造体
type EchoServer struct {
	app *echo.Echo
}

// ルーティングのセットアップ
func setupRoutes(router *echo.Echo) {
	// api ルートのグループ
	v1 := router.Group("/api/v1")
	{
		// クエリパラメーターによる認証
		v1.Use(TokenAuth)

		// todo ルートのグループ
		todo := v1.Group("/todos")
		{
			// 一覧取得
			todo.GET("", getTodos).Name = "api.todos.index" // ルートに名付けできるけどグループに名前をつけることはできない
			// 詳細取得
			todo.GET("/:id", getTodoById).Name = "api.todos.show"
			// 作成
			todo.POST("", createTodo).Name = "api.todos.store"
		}
	}
}

func TokenAuth(next echo.HandlerFunc) echo.HandlerFunc {
	// EchoではContextがインタフェースで提供されている
	return func(ctx echo.Context) error {

		token := ctx.QueryParam("token")

		// トークンが一致しない場合は 401 を返す
		if shared.IsInvalidToken(token) {
			return ctx.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Unauthorized",
			})
		}

		return next(ctx)
	}
}

func getTodos(ctx echo.Context) error {
	todos := repository.GetTodos()

	return ctx.JSON(http.StatusOK, echo.Map{
		"todos": todos,
	})
}

func getTodoById(ctx echo.Context) error {
	// パスパラメーターから id を取得
	id := ctx.Param("id")

	todo, ok := repository.GetTodoById(id)

	if !ok {
		return ctx.JSON(http.StatusNotFound, echo.Map{
			"message": "Not Found",
		})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"todo": todo,
	})
}

func createTodo(ctx echo.Context) error {
	todo := &repository.TodoForCreate{}

	// リクエストボディをパース
	if err := ctx.Bind(todo); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"message": "Bad Request",
		})
	}

	created := repository.CreateTodo(*todo)

	return ctx.JSON(http.StatusCreated, echo.Map{
		"todo": created,
	})
}

// サーバーのインスタンスを生成
func NewEchoServer() shared.Server {
	app := echo.New()

	setupRoutes(app)

	return &EchoServer{
		app,
	}
}

// サーバーを起動
func (s *EchoServer) Run(addr string) {
	s.app.Logger.Fatal(s.app.Start(addr))
}
