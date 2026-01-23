package services

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

type MockServicesRepository struct {
	Data map[string]*domain.Services
}

func (m *MockServicesRepository) List(ctx context.Context, limit, offset int) ([]*domain.Services, error) {
	var results []*domain.Services
	for _, v := range m.Data {
		results = append(results, v)
	}
	
	// Simple slicing for mock pagination
	if offset >= len(results) {
		return []*domain.Services{}, nil
	}
	end := offset + limit
	if end > len(results) {
		end = len(results)
	}
	return results[offset:end], nil
}

func (m *MockServicesRepository) Get(ctx context.Context, id string) (*domain.Services, error) {
	if val, ok := m.Data[id]; ok {
		return val, nil
	}
	return nil, nil
}

func (m *MockServicesRepository) Create(ctx context.Context, model *domain.Services) (string, error) {
	id := "test-id"
	model.ID = id
	if m.Data == nil {
		m.Data = make(map[string]*domain.Services)
	}
	m.Data[id] = model
	return id, nil
}

func (m *MockServicesRepository) Update(ctx context.Context, id string, model *domain.Services) error {
	m.Data[id] = model
	return nil
}

func (m *MockServicesRepository) Delete(ctx context.Context, id string) error {
	delete(m.Data, id)
	return nil
}

func TestServicesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &MockServicesRepository{Data: make(map[string]*domain.Services)}
	handler := NewServicesHandler(repo)
	r := gin.Default()

	r.GET("/services", handler.List)
	r.POST("/services", handler.Create)

	t.Run("Create", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := domain.Services{}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(jsonBody))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("List", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/services?page=1&limit=10", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
