package appointments

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"ServiceBookingApp/internal/domain"

	"github.com/gin-gonic/gin"
)

type AppointmentsHandler struct {
	repo          domain.AppointmentsRepository
	servicesRepo  domain.ServicesRepository
	providersRepo domain.ProvidersRepository
	schedulesRepo domain.SchedulesRepository
}

func NewAppointmentsHandler(repo domain.AppointmentsRepository, servicesRepo domain.ServicesRepository, providersRepo domain.ProvidersRepository, schedulesRepo domain.SchedulesRepository) *AppointmentsHandler {
	return &AppointmentsHandler{
		repo:          repo,
		servicesRepo:  servicesRepo,
		providersRepo: providersRepo,
		schedulesRepo: schedulesRepo,
	}
}

func (h *AppointmentsHandler) List(c *gin.Context) {
	limit := 20
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	offset := 0
	if o := c.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	} else {
		page := 1
		if p := c.Query("page"); p != "" {
			if val, err := strconv.Atoi(p); err == nil && val > 0 {
				page = val
			}
		}
		offset = (page - 1) * limit
	}

	filterType := c.Query("type")

	results, err := h.repo.List(c.Request.Context(), limit, offset, filterType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *AppointmentsHandler) Get(c *gin.Context) {
	id := c.Param("id")
	result, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AppointmentsHandler) Create(c *gin.Context) {
	var m domain.Appointments
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

func (h *AppointmentsHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var m domain.Appointments
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Update(c.Request.Context(), id, &m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *AppointmentsHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *AppointmentsHandler) GetAvailableSlots(c *gin.Context) {
	dateStr := c.Query("date")
	serviceID := c.Query("service")
	tzOffsetStr := c.Query("timezone_offset")

	if dateStr == "" || serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date and service are required"})
		return
	}

	// Default to UTC if no offset provided
	loc := time.UTC
	if tzOffsetStr != "" {
		offset, err := strconv.Atoi(tzOffsetStr)
		if err == nil {
			loc = time.FixedZone("Client", offset*60)
		}
	}

	date, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
			return
		}
		loc = date.Location()
	}

	service, err := h.servicesRepo.Get(c.Request.Context(), serviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}

	// Select first provider as default for single-tenant assumption
	providers, err := h.providersRepo.List(c.Request.Context(), 1, 0)
	if err != nil || len(providers) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no providers available"})
		return
	}
	provider := providers[0]

	schedule, err := h.schedulesRepo.GetByProvider(c.Request.Context(), provider.ID, domain.ScheduleTypeGlobal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch schedule"})
		return
	}

	if schedule == nil || len(schedule.Days) == 0 {
		c.JSON(http.StatusOK, []string{})
		return
	}

	daysOfWeek := []string{"sun", "mon", "tue", "wed", "thu", "fri", "sat"}
	dayOfWeekStr := daysOfWeek[date.Weekday()]

	daySchedule, ok := schedule.Days[dayOfWeekStr]

	if !ok || !daySchedule.Enabled {
		c.JSON(http.StatusOK, []string{})
		return
	}

	appointments, err := h.repo.ListByDate(c.Request.Context(), date)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type timeRange struct {
		Start time.Time
		End   time.Time
	}
	busySlots := []timeRange{}

	// Pre-fetch all services to get durations
	allServices, _ := h.servicesRepo.List(c.Request.Context(), 100, 0)
	serviceDurations := make(map[string]int)
	for _, s := range allServices {
		serviceDurations[s.ID] = s.DurationMinutes
	}

	for _, appt := range appointments {
		dur := serviceDurations[appt.Service]
		if dur == 0 {
			dur = 60
		}
		end := appt.ScheduledAt.Add(time.Duration(dur) * time.Minute)
		busySlots = append(busySlots, timeRange{Start: appt.ScheduledAt, End: end})
	}

	var availableSlots []string

	for _, workRange := range daySchedule.Ranges {
		var startH, startM, endH, endM int
		fmt.Sscanf(workRange.Start, "%d:%d", &startH, &startM)
		fmt.Sscanf(workRange.End, "%d:%d", &endH, &endM)

		currentSlot := time.Date(date.Year(), date.Month(), date.Day(), startH, startM, 0, 0, date.Location())
		rangeEnd := time.Date(date.Year(), date.Month(), date.Day(), endH, endM, 0, 0, date.Location())

		for currentSlot.Before(rangeEnd) {
			dur := service.DurationMinutes
			if dur <= 0 {
				dur = 30
			}
			slotEnd := currentSlot.Add(time.Duration(dur) * time.Minute)

			if slotEnd.After(rangeEnd) {
				break
			}

			isBusy := false
			for _, busy := range busySlots {
				// Overlap if (StartA < EndB) and (EndA > StartB)
				if currentSlot.Before(busy.End) && slotEnd.After(busy.Start) {
					isBusy = true
					break
				}
			}

			if !isBusy {
				availableSlots = append(availableSlots, currentSlot.Format("15:04"))
			}

			currentSlot = currentSlot.Add(30 * time.Minute)
		}
	}

	// Ensure returning empty array instead of null for JSON serialization
	if availableSlots == nil {
		availableSlots = []string{}
	}

	c.JSON(http.StatusOK, availableSlots)
}
