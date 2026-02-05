package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ServiceBookingApp/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockProvidersRepository struct {
	Data map[string]*domain.Providers
}

func (m *MockProvidersRepository) List(ctx context.Context, limit, offset int) ([]*domain.Providers, error) {
	var results []*domain.Providers
	for _, v := range m.Data {
		results = append(results, v)
	}
	
	if offset >= len(results) {
		return []*domain.Providers{}, nil
	}
	end := offset + limit
	if end > len(results) {
		end = len(results)
	}
	return results[offset:end], nil
}

func (m *MockProvidersRepository) Get(ctx context.Context, id string) (*domain.Providers, error) {
	if val, ok := m.Data[id]; ok {
		return val, nil
	}
	return nil, nil
}

func (m *MockProvidersRepository) Create(ctx context.Context, model *domain.Providers) (string, error) {
	id := "test-id"
	model.ID = id
	if m.Data == nil {
		m.Data = make(map[string]*domain.Providers)
	}
	m.Data[id] = model
	return id, nil
}

func (m *MockProvidersRepository) Update(ctx context.Context, id string, model *domain.Providers) error {
	m.Data[id] = model
	return nil
}

func (m *MockProvidersRepository) Delete(ctx context.Context, id string) error {
	delete(m.Data, id)
	return nil
}

func TestProvidersHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &MockProvidersRepository{Data: make(map[string]*domain.Providers)}
	handler := NewProvidersHandler(repo)
	r := gin.Default()

	r.GET("/providers", handler.List)
	r.POST("/providers", handler.Create)

	t.Run("Create", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := domain.Providers{}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/providers", bytes.NewBuffer(jsonBody))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("List", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/providers?page=1&limit=10", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
