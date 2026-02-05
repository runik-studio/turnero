package providers

import (
	"net/http"
	"strconv"
	"ServiceBookingApp/internal/domain"
	"ServiceBookingApp/internal/utils"
	"github.com/gin-gonic/gin"
)

type ProvidersHandler struct {
	repo domain.ProvidersRepository
}

func NewProvidersHandler(repo domain.ProvidersRepository) *ProvidersHandler {
	return &ProvidersHandler{repo: repo}
}

func (h *ProvidersHandler) List(c *gin.Context) {
	limit := 20
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}
	page := 1
	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	offset := (page - 1) * limit

	results, err := h.repo.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *ProvidersHandler) Get(c *gin.Context) {
	id := c.Param("id")
	result, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ProvidersHandler) Create(c *gin.Context) {
	var m domain.Providers
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.repo.Create(c.Request.Context(), &m)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	m.ID = id
	c.JSON(http.StatusCreated, m)
}

func (h *ProvidersHandler) Update(c *gin.Context) {
	id := c.Param("id")
	
	// Get the existing provider first
	existing, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
		return
	}
	
	// Parse the update request
	var updates domain.Providers
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Update only the fields that are provided (non-zero values)
	if updates.Phone != "" {
		existing.Phone = updates.Phone
	}
	if updates.Address != "" {
		existing.Address = updates.Address
	}
	if updates.AvatarUrl != "" {
		existing.AvatarUrl = updates.AvatarUrl
	}
	if updates.EstablishmentName != "" {
		existing.EstablishmentName = updates.EstablishmentName
	}
	
	// Update UpdatedAt timestamp
	existing.UpdatedAt = utils.Now()
	
	if err := h.repo.Update(c.Request.Context(), id, existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Return the updated provider object
	c.JSON(http.StatusOK, existing)
}

func (h *ProvidersHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	
	// Get the provider first
	provider, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
		return
	}
	
	// Set DeletedAt to current time (soft delete)
	now := utils.Now()
	provider.DeletedAt = &now
	
	// Update the provider with DeletedAt set
	if err := h.repo.Update(c.Request.Context(), id, provider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
