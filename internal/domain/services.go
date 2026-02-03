package domain

import (
	"context"
)

type Services struct {
	ID string `json:"id" firestore:"-"`

	Description interface{} `json:"description" firestore:"Description"`

	DurationMinutes int `json:"duration_minutes" firestore:"DurationMinutes"`

	IconUrl string `json:"icon_url" firestore:"IconUrl"`

	Price float64 `json:"price" firestore:"Price"`

	Color string `json:"color" firestore:"Color"`

	Title string `json:"title" firestore:"Title"`
}

type ServicesRepository interface {
	List(ctx context.Context, limit, offset int) ([]*Services, error)
	Get(ctx context.Context, id string) (*Services, error)
	Create(ctx context.Context, model *Services) (string, error)
	Update(ctx context.Context, id string, model *Services) error
	Delete(ctx context.Context, id string) error
}
