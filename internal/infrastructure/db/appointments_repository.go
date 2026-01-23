package db

import (
	"context"
	
	"ServiceBookingApp/internal/domain"
	"google.golang.org/api/iterator"
)

type AppointmentsRepository struct {
	client *FirestoreRepository
}

func NewAppointmentsRepository(client *FirestoreRepository) *AppointmentsRepository {
	return &AppointmentsRepository{client: client}
}

// GetByEmail is used for JWT auth


func (r *AppointmentsRepository) List(ctx context.Context, limit, offset int) ([]*domain.Appointments, error) {
	iter := r.client.client.Collection("appointments").Offset(offset).Limit(limit).Documents(ctx)
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
