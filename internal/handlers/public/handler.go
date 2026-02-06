package public

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"ServiceBookingApp/internal/domain"
	"ServiceBookingApp/internal/utils"

	"github.com/gin-gonic/gin"
)

type PublicHandler struct {
	servicesRepo     domain.ServicesRepository
	schedulesRepo    domain.SchedulesRepository
	appointmentsRepo domain.AppointmentsRepository
	providersRepo    domain.ProvidersRepository
}

func NewPublicHandler(servicesRepo domain.ServicesRepository, schedulesRepo domain.SchedulesRepository, appointmentsRepo domain.AppointmentsRepository, providersRepo domain.ProvidersRepository) *PublicHandler {
	return &PublicHandler{
		servicesRepo:     servicesRepo,
		schedulesRepo:    schedulesRepo,
		appointmentsRepo: appointmentsRepo,
		providersRepo:    providersRepo,
	}
}

func (h *PublicHandler) GetServices(c *gin.Context) {
	providerId := c.Param("provider_id")
	if providerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider_id is required"})
		return
	}

	services, err := h.servicesRepo.List(c.Request.Context(), 100, 0, providerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

func (h *PublicHandler) GetAvailableSlots(c *gin.Context) {
	providerId := c.Param("provider_id")
	dateStr := c.Query("date")
	serviceID := c.Query("service")
	tzOffsetStr := c.Query("timezone_offset")

	if providerId == "" || dateStr == "" || serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider_id, date and service are required"})
		return
	}

	// Verify service belongs to provider
	service, err := h.servicesRepo.Get(c.Request.Context(), serviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}
	if service.ProviderId != providerId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service does not belong to provider"})
		return
	}

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

	schedule, err := h.schedulesRepo.GetByProvider(c.Request.Context(), providerId, domain.ScheduleTypeGlobal)
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

	// Fetch all appointments for the day to check availability
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	appointments, err := h.appointmentsRepo.ListByDate(c.Request.Context(), startOfDay, providerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type timeRange struct {
		Start time.Time
		End   time.Time
	}
	busySlots := []timeRange{}

	allServices, _ := h.servicesRepo.List(c.Request.Context(), 100, 0, providerId)
	serviceDurations := make(map[string]int)
	for _, s := range allServices {
		serviceDurations[s.ID] = s.DurationMinutes
	}

	for _, appt := range appointments {
		dur := serviceDurations[appt.ServiceId]
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
				// Check for overlap
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

	if availableSlots == nil {
		availableSlots = []string{}
	}

	c.JSON(http.StatusOK, availableSlots)
}

func (h *PublicHandler) CreateAppointment(c *gin.Context) {
	providerId := c.Param("provider_id")
	if providerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider_id is required"})
		return
	}

	var m domain.Appointments
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m.ProviderId = providerId
	// Verify service
	if m.ServiceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service_id is required"})
		return
	}
	service, err := h.servicesRepo.Get(c.Request.Context(), m.ServiceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service_id"})
		return
	}
	if service.ProviderId != providerId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service does not belong to provider"})
		return
	}

	now := utils.Now()
	// Allow creating slightly in the past (e.g. 5 minutes ago) to account for clock skew/delays?
	// But mostly future.
	if m.ScheduledAt.Before(now.Add(-5 * time.Minute)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot create appointment in the past"})
		return
	}

	// Prevent duplicates / Overlap check
	// ListByDate fetches for the whole day.
	date := m.ScheduledAt
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, waitingLocation(date.Location()))
	
	appointments, err := h.appointmentsRepo.ListByDate(c.Request.Context(), startOfDay, providerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check availability"})
		return
	}

	newApptStart := m.ScheduledAt
	dur := service.DurationMinutes
	if dur == 0 {
		dur = 30
	}
	newApptEnd := newApptStart.Add(time.Duration(dur) * time.Minute)

	// We need durations of other appointments to check overlap correctly
	// Optimization: we could cache or just fetch all services once outside
	// Ideally we already have them if we want to be super precise or just fetch them if needed.
	// For now, let's fetch all services for the provider again or optimize?
	// Let's optimize: fetch list of services only if we haven't already?
	// Actually, fetching all services is okay since providers don't have thousands.
	
	allServices, _ := h.servicesRepo.List(c.Request.Context(), 100, 0, providerId)
	serviceDurations := make(map[string]int)
	for _, s := range allServices {
		serviceDurations[s.ID] = s.DurationMinutes
	}

	for _, existing := range appointments {
		existingDur := serviceDurations[existing.ServiceId]
		if existingDur == 0 {
			existingDur = 60 // Default or unknown
		}
		existingStart := existing.ScheduledAt
		existingEnd := existingStart.Add(time.Duration(existingDur) * time.Minute)

		// Check overlap
		// (StartA < EndB) and (EndA > StartB)
		if newApptStart.Before(existingEnd) && newApptEnd.After(existingStart) {
			c.JSON(http.StatusConflict, gin.H{"error": "time slot is already booked"})
			return
		}
	}
	
	m.Status = "confirmed"
	
	id, err := h.appointmentsRepo.Create(c.Request.Context(), &m)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	m.ID = id
	c.JSON(http.StatusCreated, m)
}

func waitingLocation(loc *time.Location) *time.Location {
	if loc == nil {
		return time.UTC
	}
	return loc
}
