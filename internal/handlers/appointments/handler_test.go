package appointments

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

type MockAppointmentsRepository struct {
	Data map[string]*domain.Appointments
}

func (m *MockAppointmentsRepository) List(ctx context.Context, limit, offset int) ([]*domain.Appointments, error) {
	var results []*domain.Appointments
	for _, v := range m.Data {
		results = append(results, v)
	}
	
	// Simple slicing for mock pagination
	if offset >= len(results) {
		return []*domain.Appointments{}, nil
	}
	end := offset + limit
	if end > len(results) {
		end = len(results)
	}
	return results[offset:end], nil
}

func (m *MockAppointmentsRepository) Get(ctx context.Context, id string) (*domain.Appointments, error) {
	if val, ok := m.Data[id]; ok {
		return val, nil
	}
	return nil, nil
}

func (m *MockAppointmentsRepository) Create(ctx context.Context, model *domain.Appointments) (string, error) {
	id := "test-id"
	model.ID = id
	if m.Data == nil {
		m.Data = make(map[string]*domain.Appointments)
	}
	m.Data[id] = model
	return id, nil
}

func (m *MockAppointmentsRepository) Update(ctx context.Context, id string, model *domain.Appointments) error {
	m.Data[id] = model
	return nil
}

func (m *MockAppointmentsRepository) Delete(ctx context.Context, id string) error {
	delete(m.Data, id)
	return nil
}

func TestAppointmentsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &MockAppointmentsRepository{Data: make(map[string]*domain.Appointments)}
	handler := NewAppointmentsHandler(repo)
	r := gin.Default()

	r.GET("/appointments", handler.List)
	r.POST("/appointments", handler.Create)

	t.Run("Create", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := domain.Appointments{}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/appointments", bytes.NewBuffer(jsonBody))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("List", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/appointments?page=1&limit=10", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
