package db

import (
	"context"
	
	"ServiceBookingApp/internal/domain"
	"google.golang.org/api/iterator"
)

type ProvidersRepository struct {
	client *FirestoreRepository
}

func NewProvidersRepository(client *FirestoreRepository) *ProvidersRepository {
	return &ProvidersRepository{client: client}
}

// GetByEmail is used for JWT auth


func (r *ProvidersRepository) List(ctx context.Context, limit, offset int) ([]*domain.Providers, error) {
	iter := r.client.client.Collection("providers").Offset(offset).Limit(limit).Documents(ctx)
	var results []*domain.Providers
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var m domain.Providers
		if err := doc.DataTo(&m); err != nil {
			return nil, err
		}
		m.ID = doc.Ref.ID
		results = append(results, &m)
	}
	return results, nil
}

func (r *ProvidersRepository) Get(ctx context.Context, id string) (*domain.Providers, error) {
	doc, err := r.client.client.Collection("providers").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var m domain.Providers
	if err := doc.DataTo(&m); err != nil {
		return nil, err
	}
	m.ID = doc.Ref.ID
	return &m, nil
}

func (r *ProvidersRepository) Create(ctx context.Context, model *domain.Providers) (string, error) {
	now := getNow()
	model.CreatedAt = now
	model.UpdatedAt = now
	ref, _, err := r.client.client.Collection("providers").Add(ctx, model)
	if err != nil {
		return "", err
	}
	return ref.ID, nil
}



func (r *ProvidersRepository) Update(ctx context.Context, id string, m *domain.Providers) error {
	m.UpdatedAt = getNow()
	_, err := r.client.client.Collection("providers").Doc(id).Set(ctx, m)
	return err
}

func (r *ProvidersRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.client.Collection("providers").Doc(id).Delete(ctx)
	return err
}
