package db

import (
	"context"
	"time"

	"ServiceBookingApp/internal/domain"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type AppointmentsRepository struct {
	client *FirestoreRepository
}

func NewAppointmentsRepository(client *FirestoreRepository) *AppointmentsRepository {
	return &AppointmentsRepository{client: client}
}

// GetByEmail is used for JWT auth

func (r *AppointmentsRepository) ListByDate(ctx context.Context, date time.Time) ([]*domain.Appointments, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	iter := r.client.client.Collection("appointments").
		Where("ScheduledAt", ">=", startOfDay).
		Where("ScheduledAt", "<", endOfDay).
		Documents(ctx)

	var results []*domain.Appointments
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var m domain.Appointments
		if err := doc.DataTo(&m); err != nil {
			return nil, err
		}
		m.ID = doc.Ref.ID
		results = append(results, &m)
	}
	return results, nil
}

func (r *AppointmentsRepository) List(ctx context.Context, limit, offset int, filterType string) ([]*domain.Appointments, error) {
	query := r.client.client.Collection("appointments").Query
	now := time.Now()

	if filterType == "upcoming" {
		query = query.Where("ScheduledAt", ">=", now).OrderBy("ScheduledAt", firestore.Asc)
	} else if filterType == "past" {
		query = query.Where("ScheduledAt", "<", now).OrderBy("ScheduledAt", firestore.Desc)
	} else {
		// Default ordering if no type specified
		query = query.OrderBy("ScheduledAt", firestore.Desc)
	}

	iter := query.Offset(offset).Limit(limit).Documents(ctx)
	var results []*domain.Appointments
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var m domain.Appointments
		if err := doc.DataTo(&m); err != nil {
			return nil, err
		}
		m.ID = doc.Ref.ID
		results = append(results, &m)
	}
	return results, nil
}

func (r *AppointmentsRepository) Get(ctx context.Context, id string) (*domain.Appointments, error) {
	doc, err := r.client.client.Collection("appointments").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var m domain.Appointments
	if err := doc.DataTo(&m); err != nil {
		return nil, err
	}
	m.ID = doc.Ref.ID
	return &m, nil
}

func (r *AppointmentsRepository) Create(ctx context.Context, model *domain.Appointments) (string, error) {
	ref, _, err := r.client.client.Collection("appointments").Add(ctx, model)
	if err != nil {
		return "", err
	}
	return ref.ID, nil
}

func (r *AppointmentsRepository) Update(ctx context.Context, id string, m *domain.Appointments) error {
	_, err := r.client.client.Collection("appointments").Doc(id).Set(ctx, m)
	return err
}

func (r *AppointmentsRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.client.Collection("appointments").Doc(id).Delete(ctx)
	return err
}
