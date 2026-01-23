package domain

import (
	"context"
	"time"
)

type Appointments struct {
	ID string `json:"id" bson:"_id,omitempty"`
	
	Notes interface{} `json:"notes" bson:"notes"`
	
	ScheduledAt time.Time `json:"scheduled_at" bson:"scheduled_at"`
	
	Status string `json:"status" bson:"status"`
	
	
	Provider string `json:"provider" bson:"provider"`
	
	Service string `json:"service" bson:"service"`
	
	User string `json:"user" bson:"user"`
	
}

type AppointmentsRepository interface {
	List(ctx context.Context, limit, offset int) ([]*Appointments, error)
	Get(ctx context.Context, id string) (*Appointments, error)
	Create(ctx context.Context, model *Appointments) (string, error)
	Update(ctx context.Context, id string, model *Appointments) error
	Delete(ctx context.Context, id string) error
}
