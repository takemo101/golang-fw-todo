package fiberapp

import (
	"app/repository"
	"app/shared"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	app := fiber.New()
	setupRoutes(app)

	t.Run("無効なトークンのテスト", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/todos?token=invalid", nil)
		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 401, res.StatusCode)
	})

	t.Run("Todo一覧取得のテスト", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/todos?token="+shared.Token, nil)
		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)
	})

	t.Run("Todo詳細取得のテスト", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/todos/first?token="+shared.Token, nil)
		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)
	})

	t.Run("Todo作成のテスト", func(t *testing.T) {
		todo := repository.TodoForCreate{
			Title:     "test",
			Completed: false,
		}

		data, _ := json.Marshal(todo)

		req, _ := http.NewRequest("POST", "/api/v1/todos?token="+shared.Token, strings.NewReader(string(data)))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 201, res.StatusCode)
		// レスポンスの body をパースして Todo を取得
		var body struct {
			Todo repository.Todo `json:"todo"`
		}
		json.NewDecoder(res.Body).Decode(&body)

		assert.Equal(t, todo.Title, body.Todo.Title)
		assert.Equal(t, todo.Completed, body.Todo.Completed)
	})
}
