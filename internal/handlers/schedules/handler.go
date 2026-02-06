package schedules

import (
	"net/http"

	"ServiceBookingApp/internal/domain"
	
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

type SchedulesHandler struct {
	repo          domain.SchedulesRepository
	providersRepo domain.ProvidersRepository
}

func NewSchedulesHandler(repo domain.SchedulesRepository, providersRepo domain.ProvidersRepository) *SchedulesHandler {
	return &SchedulesHandler{
		repo:          repo,
		providersRepo: providersRepo,
	}
}

func (h *SchedulesHandler) GetByProvider(c *gin.Context) {
	providerID := c.Query("provider_id")
	if providerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider_id is required"})
		return
	}

	scheduleType := c.Query("type")
	if scheduleType == "" {
		scheduleType = string(domain.ScheduleTypeGlobal)
	}

	schedule, err := h.repo.GetByProvider(c.Request.Context(), providerID, domain.ScheduleType(scheduleType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if schedule == nil {
		c.JSON(http.StatusOK, domain.Schedule{
			ProviderId: providerID,
			Type:       domain.ScheduleType(scheduleType),
			Days:       make(map[string]domain.DaySchedule),
		})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func (h *SchedulesHandler) Upsert(c *gin.Context) {
	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	token := u.(*auth.Token)

	provider, err := h.providersRepo.GetByUserId(c.Request.Context(), token.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if provider == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a provider"})
		return
	}

	var schedule domain.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule.ProviderId = provider.ID

	if err := h.repo.Upsert(c.Request.Context(), &schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}
