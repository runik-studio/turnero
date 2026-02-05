package db

import (
	"context"
	
	"ServiceBookingApp/internal/domain"
	"google.golang.org/api/iterator"
)

type UsersRepository struct {
	client *FirestoreRepository
}

func NewUsersRepository(client *FirestoreRepository) *UsersRepository {
	return &UsersRepository{client: client}
}

func (r *UsersRepository) List(ctx context.Context, limit, offset int) ([]*domain.Users, error) {
	iter := r.client.client.Collection("users").Offset(offset).Limit(limit).Documents(ctx)
	var results []*domain.Users
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var m domain.Users
		if err := doc.DataTo(&m); err != nil {
			return nil, err
		}
		m.ID = doc.Ref.ID
		results = append(results, &m)
	}
	return results, nil
}

func (r *UsersRepository) Get(ctx context.Context, id string) (*domain.Users, error) {
	doc, err := r.client.client.Collection("users").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var m domain.Users
	if err := doc.DataTo(&m); err != nil {
		return nil, err
	}
	m.ID = doc.Ref.ID
	return &m, nil
}

func (r *UsersRepository) Create(ctx context.Context, model *domain.Users) (string, error) {
	ref, _, err := r.client.client.Collection("users").Add(ctx, model)
	if err != nil {
		return "", err
	}
	return ref.ID, nil
}

func (r *UsersRepository) Update(ctx context.Context, id string, m *domain.Users) error {
	_, err := r.client.client.Collection("users").Doc(id).Set(ctx, m)
	return err
}

func (r *UsersRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.client.Collection("users").Doc(id).Delete(ctx)
	return err
}
