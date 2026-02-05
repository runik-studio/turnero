package db

import (
	"context"
	"errors"

	"ServiceBookingApp/internal/domain"
	"ServiceBookingApp/internal/utils"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type SchedulesRepository struct {
	client *FirestoreRepository
}

func NewSchedulesRepository(client *FirestoreRepository) *SchedulesRepository {
	return &SchedulesRepository{client: client}
}

func (r *SchedulesRepository) GetByProvider(ctx context.Context, providerID string, scheduleType domain.ScheduleType) (*domain.Schedule, error) {
	// Query for schedule with matching provider_id and type
	iter := r.client.client.Collection("schedules").
		Where("ProviderId", "==", providerID).
		Where("Type", "==", scheduleType).
		Limit(1).
		Documents(ctx)

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}

	var s domain.Schedule
	if err := doc.DataTo(&s); err != nil {
		return nil, err
	}
	s.ID = doc.Ref.ID
	return &s, nil
}

func (r *SchedulesRepository) Upsert(ctx context.Context, schedule *domain.Schedule) error {
	if schedule.ProviderId == "" {
		return errors.New("provider_id is required")
	}

	collection := r.client.client.Collection("schedules")

	var docRef *firestore.DocumentRef

	if schedule.ID != "" {
		docRef = collection.Doc(schedule.ID)
	} else {
		// Try to find existing one to update by provider and type
		existing, err := r.GetByProvider(ctx, schedule.ProviderId, schedule.Type)
		if err != nil {
			return err
		}
		if existing != nil {
			docRef = collection.Doc(existing.ID)
			schedule.ID = existing.ID
			schedule.CreatedAt = existing.CreatedAt // Preserve created_at
		} else {
			// Create new
			docRef = collection.NewDoc()
			schedule.ID = docRef.ID
			schedule.CreatedAt = utils.Now()
		}
	}

	schedule.UpdatedAt = utils.Now()

	_, err := docRef.Set(ctx, schedule)
	return err
}
