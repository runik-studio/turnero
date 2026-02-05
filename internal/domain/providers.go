package domain

import (
	"context"
	"time"
)

type Providers struct {
	ID string `json:"id" firestore:"-"`

	UserId string `json:"user_id" firestore:"UserId"`

	Address string `json:"address" firestore:"Address"`

	AvatarUrl string `json:"avatar_url" firestore:"AvatarUrl"`

	EstablishmentName string `json:"establishment_name" firestore:"EstablishmentName"`

	Phone string `json:"phone" firestore:"Phone"`

	CreatedAt time.Time `json:"created_at" firestore:"CreatedAt"`
	UpdatedAt time.Time `json:"updated_at" firestore:"UpdatedAt"`
}

type DaySchedule struct {
	Ranges  []TimeRange `json:"ranges" firestore:"Ranges"`
	Enabled bool        `json:"enabled" firestore:"Enabled"`
}

type TimeRange struct {
	Start string `json:"start" firestore:"Start"`
	End   string `json:"end" firestore:"End"`
}

type ProvidersRepository interface {
	List(ctx context.Context, limit, offset int) ([]*Providers, error)
	Get(ctx context.Context, id string) (*Providers, error)
	Create(ctx context.Context, model *Providers) (string, error)
	Update(ctx context.Context, id string, model *Providers) error
	Delete(ctx context.Context, id string) error
}
