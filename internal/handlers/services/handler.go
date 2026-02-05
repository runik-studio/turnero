package services

import (
	"net/http"
	"strconv"
	"ServiceBookingApp/internal/domain"
	"ServiceBookingApp/internal/utils"
	"github.com/gin-gonic/gin"
)

type ServicesHandler struct {
	repo domain.ServicesRepository
}

func NewServicesHandler(repo domain.ServicesRepository) *ServicesHandler {
	return &ServicesHandler{repo: repo}
}

func (h *ServicesHandler) List(c *gin.Context) {
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

func (h *ServicesHandler) Get(c *gin.Context) {
	id := c.Param("id")
	result, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ServicesHandler) Create(c *gin.Context) {
	var m domain.Services
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

func (h *ServicesHandler) Update(c *gin.Context) {
	id := c.Param("id")
	
	existing, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
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
	
	service, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
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
