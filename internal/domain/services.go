package domain

import (
	"context"
	"time"
)

type Services struct {
	ID string `json:"id" firestore:"-"`

	ProviderId string `json:"provider_id" firestore:"ProviderId"`

	Description interface{} `json:"description" firestore:"Description"`

	DurationMinutes int `json:"duration_minutes" firestore:"DurationMinutes"`

	IconUrl string `json:"icon_url" firestore:"IconUrl"`

	Price float64 `json:"price" firestore:"Price"`

	Color string `json:"color" firestore:"Color"`

	Title string `json:"title" firestore:"Title"`

	CreatedAt time.Time  `json:"created_at" firestore:"CreatedAt"`
	UpdatedAt time.Time  `json:"updated_at" firestore:"UpdatedAt"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" firestore:"DeletedAt,omitempty"`
}

type ServicesRepository interface {
	List(ctx context.Context, limit, offset int) ([]*Services, error)
	Get(ctx context.Context, id string) (*Services, error)
	Create(ctx context.Context, model *Services) (string, error)
	Update(ctx context.Context, id string, model *Services) error
	Delete(ctx context.Context, id string) error
}
