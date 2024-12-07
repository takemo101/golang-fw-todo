package chiapp

import (
	"app/repository"
	"app/shared"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// contract.Serverを実装したサーバーの構造体
type ChiServer struct {
	mux *chi.Mux
}

// Jsonレスポンスの型
type ChiMap map[string]interface{}

// Todo作成リクエストの型
type CreateTodoRequest struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// render.Bidnerインタフェースの実装
func (todo *CreateTodoRequest) Bind(r *http.Request) error {
	return nil
}

// ルーティングのセットアップ
func setupRoutes(router *chi.Mux) {

	// api ルートのグループ
	router.Route("/api/v1", func(r chi.Router) {

		// クエリパラメーターによる認証
		r.Use(func(next http.Handler) http.Handler {
			// 他のフレームワークと違い、Contextがインタフェースで提供されている
			return http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {

					token := r.URL.Query().Get("token")

					// トークンが一致しない場合は 401 を返す
					if shared.IsInvalidToken(token) {
						render.Status(r, http.StatusUnauthorized)
						render.JSON(w, r, ChiMap{
							"message": "Unauthorized",
						})
					}

					next.ServeHTTP(w, r)
				},
			)
		})

		// todo ルートのグループ
		r.Route("/todos", func(todo chi.Router) {
			// 一覧取得
			todo.Get("/", func(w http.ResponseWriter, r *http.Request) {

				todos := repository.GetTodos()

				render.JSON(w, r, ChiMap{
					"todos": todos,
				})
			})

			// 詳細取得
			todo.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
				// パスパラメーターから id を取得
				id := chi.URLParam(r, "id")

				todo, ok := repository.GetTodoById(id)

				if !ok {
					render.Status(r, http.StatusNotFound)
					render.JSON(w, r, ChiMap{
						"message": "Not Found",
					})
				}

				render.JSON(w, r, ChiMap{
					"todo": todo,
				})
			})

			// 作成
			todo.Post("/", func(w http.ResponseWriter, r *http.Request) {
				todo := &CreateTodoRequest{}

				// リクエストボディをパース
				if err := render.Bind(r, todo); err != nil {
					render.Status(r, http.StatusBadRequest)
					render.JSON(w, r, ChiMap{
						"message": "Bad Request",
					})
				}

				created := repository.CreateTodo(repository.TodoForCreate{
					Title:     todo.Title,
					Completed: todo.Completed,
				})

				render.Status(r, http.StatusCreated)
				render.JSON(w, r, ChiMap{
					"todo": created,
				})
			})
		})
	})
}

// サーバーのインスタンスを生成
func NewChiServer() shared.Server {
	mux := chi.NewMux()

	setupRoutes(mux)

	return &ChiServer{
		mux,
	}
}

// サーバーを起動
func (s *ChiServer) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, s.mux))
}
