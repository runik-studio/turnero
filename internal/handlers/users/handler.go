package users

import (
	"net/http"
	"strconv"
	"ServiceBookingApp/internal/domain"
	"ServiceBookingApp/internal/utils"
	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	repo domain.UsersRepository
}

func NewUsersHandler(repo domain.UsersRepository) *UsersHandler {
	return &UsersHandler{repo: repo}
}

func (h *UsersHandler) List(c *gin.Context) {
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

func (h *UsersHandler) Get(c *gin.Context) {
	id := c.Param("id")
	result, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *UsersHandler) Create(c *gin.Context) {
	var m domain.Users
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

func (h *UsersHandler) Update(c *gin.Context) {
	id := c.Param("id")
	
	// Get the existing user first
	existing, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	
	// Parse the update request
	var updates domain.Users
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Update only the fields that are provided (non-zero values)
	if updates.Name != "" {
		existing.Name = updates.Name
	}
	if updates.Email != "" {
		existing.Email = updates.Email
	}
	if updates.Picture != "" {
		existing.Picture = updates.Picture
	}
	if updates.RoleId != "" {
		existing.RoleId = updates.RoleId
	}
	
	if err := h.repo.Update(c.Request.Context(), id, existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Return the updated user object
	c.JSON(http.StatusOK, existing)
}

func (h *UsersHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	
	// Get the user first
	user, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	
	// Set DeletedAt to current time (soft delete)
	now := utils.Now()
	user.DeletedAt = &now
	
	// Update the user with DeletedAt set
	if err := h.repo.Update(c.Request.Context(), id, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
