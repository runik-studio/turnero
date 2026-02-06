package providers

import (
	"net/http"

	"ServiceBookingApp/internal/domain"
	"ServiceBookingApp/internal/utils"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

type ProvidersHandler struct {
	repo domain.ProvidersRepository
}

func NewProvidersHandler(repo domain.ProvidersRepository) *ProvidersHandler {
	return &ProvidersHandler{repo: repo}
}

func (h *ProvidersHandler) List(c *gin.Context) {
	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	token := u.(*auth.Token)

	provider, err := h.repo.GetByUserId(c.Request.Context(), token.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	results := []*domain.Providers{}
	if provider != nil {
		results = append(results, provider)
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
	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	token := u.(*auth.Token)

	existing, err := h.repo.GetByUserId(c.Request.Context(), token.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already has a provider"})
		return
	}

	var m domain.Providers
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m.UserId = token.UID

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
	
	existing, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
		return
	}
	
	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	token := u.(*auth.Token)
	if existing.UserId != token.UID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	
	var updates domain.Providers
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
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
	
	existing.UpdatedAt = utils.Now()
	
	if err := h.repo.Update(c.Request.Context(), id, existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, existing)
}

func (h *ProvidersHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	
	provider, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
		return
	}

	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	token := u.(*auth.Token)
	if provider.UserId != token.UID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	
	now := utils.Now()
	provider.DeletedAt = &now
	
	if err := h.repo.Update(c.Request.Context(), id, provider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
