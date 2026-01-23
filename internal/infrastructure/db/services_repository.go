package db

import (
	"context"
	
	"ServiceBookingApp/internal/domain"
	"google.golang.org/api/iterator"
)

type ServicesRepository struct {
	client *FirestoreRepository
}

func NewServicesRepository(client *FirestoreRepository) *ServicesRepository {
	return &ServicesRepository{client: client}
}

// GetByEmail is used for JWT auth


func (r *ServicesRepository) List(ctx context.Context, limit, offset int) ([]*domain.Services, error) {
	iter := r.client.client.Collection("services").Offset(offset).Limit(limit).Documents(ctx)
	var results []*domain.Services
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var m domain.Services
		if err := doc.DataTo(&m); err != nil {
			return nil, err
		}
		m.ID = doc.Ref.ID
		results = append(results, &m)
	}
	return results, nil
}

func (r *ServicesRepository) Get(ctx context.Context, id string) (*domain.Services, error) {
	doc, err := r.client.client.Collection("services").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var m domain.Services
	if err := doc.DataTo(&m); err != nil {
		return nil, err
	}
	m.ID = doc.Ref.ID
	return &m, nil
}

func (r *ServicesRepository) Create(ctx context.Context, model *domain.Services) (string, error) {
	ref, _, err := r.client.client.Collection("services").Add(ctx, model)
	if err != nil {
		return "", err
	}
	return ref.ID, nil
}



func (r *ServicesRepository) Update(ctx context.Context, id string, m *domain.Services) error {
	_, err := r.client.client.Collection("services").Doc(id).Set(ctx, m)
	return err
}

func (r *ServicesRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.client.Collection("services").Doc(id).Delete(ctx)
	return err
}
