package domain

import (
	"context"
	"time"
)

type Appointments struct {
	ID string `json:"id" firestore:"-"`

	Notes interface{} `json:"notes" firestore:"Notes"`

	ScheduledAt time.Time `json:"scheduled_at" firestore:"ScheduledAt"`

	Status string `json:"status" firestore:"Status"`

	Provider string `json:"provider" firestore:"Provider"`

	Service string `json:"service" firestore:"Service"`

	User string `json:"user" firestore:"User"`
}

type AppointmentsRepository interface {
	List(ctx context.Context, limit, offset int, filterType string) ([]*Appointments, error)
	Get(ctx context.Context, id string) (*Appointments, error)
	ListByDate(ctx context.Context, date time.Time) ([]*Appointments, error)
	Create(ctx context.Context, model *Appointments) (string, error)
	Update(ctx context.Context, id string, model *Appointments) error
	Delete(ctx context.Context, id string) error
}
