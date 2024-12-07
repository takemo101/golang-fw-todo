package ginapp

import (
	"app/repository"
	"app/shared"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	engine := gin.Default()
	setupRoutes(engine)

	t.Run("無効なトークンのテスト", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/todos?token=invalid", nil)
		engine.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
	})

	t.Run("Todo一覧取得のテスト", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/todos?token="+shared.Token, nil)
		engine.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("Todo詳細取得のテスト", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/todos/first?token="+shared.Token, nil)
		engine.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("Todo作成のテスト", func(t *testing.T) {
		todo := repository.TodoForCreate{
			Title:     "test",
			Completed: false,
		}

		data, _ := json.Marshal(todo)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/todos?token="+shared.Token, strings.NewReader(string(data)))
		req.Header.Set("Content-Type", "application/json")
		engine.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
		// レスポンスの body をパースして Todo を取得
		var body struct {
			Todo repository.Todo `json:"todo"`
		}
		json.NewDecoder(w.Body).Decode(&body)

		assert.Equal(t, todo.Title, body.Todo.Title)
		assert.Equal(t, todo.Completed, body.Todo.Completed)
	})
}
