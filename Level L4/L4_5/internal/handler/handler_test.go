package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api_optimization/internal/entity"
	"api_optimization/internal/usecase"

	"github.com/gin-gonic/gin"
)

func BenchmarkSumHandler(b *testing.B) {
	gin.SetMode(gin.TestMode)
	h := New(usecase.New())

	// Два разных входа
	inputs := []entity.Input{
		{A: 1, B: 2},
		{A: 2, B: 3},
	}

	// Преобразуем их в JSON заранее
	bodies := make([][]byte, len(inputs))
	for i, in := range inputs {
		// Добавим "нагрузку" в A и B
		for j := 0; j < 1000; j++ {
			in.A += 1
			in.B += 2
		}
		body, _ := json.Marshal(in)
		bodies[i] = body
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		body := bodies[i%len(bodies)]
		req := httptest.NewRequest(http.MethodPost, "/sum", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		h.Sum(ctx)
	}
}
