package schedules

import (
	"net/http"

	"ServiceBookingApp/internal/domain"
	"github.com/gin-gonic/gin"
)

type SchedulesHandler struct {
	repo domain.SchedulesRepository
}

func NewSchedulesHandler(repo domain.SchedulesRepository) *SchedulesHandler {
	return &SchedulesHandler{repo: repo}
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
		// Return empty schedule structure if not found
		c.JSON(http.StatusOK, domain.Schedule{
			ProviderID: providerID,
			Type:       domain.ScheduleType(scheduleType),
			Days:       make(map[string]domain.DaySchedule),
		})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func (h *SchedulesHandler) Upsert(c *gin.Context) {
	var schedule domain.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if schedule.ProviderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider_id is required"})
		return
	}

	if err := h.repo.Upsert(c.Request.Context(), &schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}
