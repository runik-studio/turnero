package domain

import (
	"context"
	"time"
)

type ScheduleType string

const (
	ScheduleTypeGlobal ScheduleType = "global"
	ScheduleTypeCustom ScheduleType = "custom"
)

type Schedule struct {
	ID         string                 `json:"id" bson:"_id,omitempty" firestore:"-"`
	ProviderID string                 `json:"provider_id" bson:"provider_id" firestore:"ProviderID"`
	Type       ScheduleType           `json:"type" bson:"type" firestore:"Type"`
	Days       map[string]DaySchedule `json:"days" bson:"days" firestore:"Days"`
	ValidFrom  *time.Time             `json:"valid_from,omitempty" bson:"valid_from,omitempty" firestore:"ValidFrom,omitempty"`
	ValidTo    *time.Time             `json:"valid_to,omitempty" bson:"valid_to,omitempty" firestore:"ValidTo,omitempty"`
	CreatedAt  time.Time              `json:"created_at" bson:"created_at" firestore:"CreatedAt"`
	UpdatedAt  time.Time              `json:"updated_at" bson:"updated_at" firestore:"UpdatedAt"`
}

type SchedulesRepository interface {
	GetByProvider(ctx context.Context, providerID string, scheduleType ScheduleType) (*Schedule, error)
	Upsert(ctx context.Context, schedule *Schedule) error
}
