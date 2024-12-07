package fiberapp

import (
	"app/repository"
	"app/shared"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

// contract.Serverを実装したサーバーの構造体
type FiberServer struct {
	app *fiber.App
}

// ルーティングのセットアップ
func setupRoutes(router *fiber.App) {

	// api ルートのグループ
	v1 := router.Group("/api/v1").Name("api.")
	{
		// クエリパラメーターによる認証
		v1.Use(func(ctx fiber.Ctx) error {
			// Fiberのv3ではContextがインタフェースで提供されている

			token := fiber.Query[string](ctx, "token")

			// トークンが一致しない場合は 401 を返す
			if shared.IsInvalidToken(token) {
				return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"message": "Unauthorized",
				})
			}

			return ctx.Next()
		})

		// todo ルートのグループ
		todo := v1.Group("/todos").Name("todos.")
		{
			// 一覧取得
			todo.Get("", func(ctx fiber.Ctx) error {

				todos := repository.GetTodos()

				return ctx.JSON(fiber.Map{
					"todos": todos,
				})
			}).Name("index") // ルートに名付けできる！便利ー！

			// 詳細取得
			todo.Get("/:id", func(ctx fiber.Ctx) error {
				// パスパラメーターから id を取得
				id := fiber.Params[string](ctx, "id")

				todo, ok := repository.GetTodoById(id)

				if !ok {
					return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
						"message": "Not Found",
					})
				}

				return ctx.JSON(fiber.Map{
					"todo": todo,
				})
			}).Name("show")

			// 作成
			todo.Post("", func(ctx fiber.Ctx) error {
				todo := &repository.TodoForCreate{}

				// リクエストボディをパース
				if err := ctx.Bind().Body(todo); err != nil {
					return err
				}

				created := repository.CreateTodo(*todo)

				return ctx.Status(http.StatusCreated).JSON(fiber.Map{
					"todo": created,
				})
			}).Name("store")
		}
	}
}

// サーバーのインスタンスを生成
func NewFiberServer() shared.Server {
	app := fiber.New()

	setupRoutes(app)

	return &FiberServer{
		app,
	}
}

// サーバーを起動
func (s *FiberServer) Run(addr string) {
	log.Fatal(s.app.Listen(addr))
}
