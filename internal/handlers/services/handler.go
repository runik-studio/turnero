package services

import (
	"fmt"
	"net/http"
	"strconv"

	"ServiceBookingApp/internal/domain"
	"ServiceBookingApp/internal/utils"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

type ServicesHandler struct {
	repo          domain.ServicesRepository
	providersRepo domain.ProvidersRepository
}

func NewServicesHandler(repo domain.ServicesRepository, providersRepo domain.ProvidersRepository) *ServicesHandler {
	return &ServicesHandler{
		repo:          repo,
		providersRepo: providersRepo,
	}
}

func (h *ServicesHandler) getProviderID(c *gin.Context) (string, error) {
	u, exists := c.Get("user")
	if !exists {
		return "", fmt.Errorf("user not found in context")
	}
	token := u.(*auth.Token)
	
	provider, err := h.providersRepo.GetByUserId(c.Request.Context(), token.UID)
	if err != nil {
		return "", err
	}
	if provider == nil {
		return "", fmt.Errorf("user is not a provider")
	}
	return provider.ID, nil
}

func (h *ServicesHandler) List(c *gin.Context) {
	providerId := c.Query("provider_id")
	if providerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider_id is required"})
		return
	}

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

	results, err := h.repo.List(c.Request.Context(), limit, offset, providerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *ServicesHandler) Get(c *gin.Context) {
	id := c.Param("id")
	result, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	// TODO: Checks regarding ownership?
	// Publicly readable? Usually yes for booking.
	c.JSON(http.StatusOK, result)
}

func (h *ServicesHandler) Create(c *gin.Context) {
	providerId, err := h.getProviderID(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "must be a provider to create services"})
		return
	}

	var m domain.Services
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m.ProviderId = providerId

	id, err := h.repo.Create(c.Request.Context(), &m)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	m.ID = id
	c.JSON(http.StatusCreated, m)
}

func (h *ServicesHandler) Update(c *gin.Context) {
	id := c.Param("id")
	
	providerId, err := h.getProviderID(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "must be a provider"})
		return
	}

	existing, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}
	
	if existing.ProviderId != providerId {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var updates domain.Services
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if updates.Title != "" {
		existing.Title = updates.Title
	}
	if updates.Description != nil {
		existing.Description = updates.Description
	}
	if updates.DurationMinutes != 0 {
		existing.DurationMinutes = updates.DurationMinutes
	}
	if updates.Price != 0 {
		existing.Price = updates.Price
	}
	if updates.IconUrl != "" {
		existing.IconUrl = updates.IconUrl
	}
	if updates.Color != "" {
		existing.Color = updates.Color
	}
	
	existing.UpdatedAt = utils.Now()
	
	if err := h.repo.Update(c.Request.Context(), id, existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, existing)
}

func (h *ServicesHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	
	providerId, err := h.getProviderID(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "must be a provider"})
		return
	}

	service, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}
	
	if service.ProviderId != providerId {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	now := utils.Now()
	service.DeletedAt = &now
	
	if err := h.repo.Update(c.Request.Context(), id, service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
