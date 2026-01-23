package domain

import (
	"context"
	
)

type Services struct {
	ID string `json:"id" bson:"_id,omitempty"`
	
	Description interface{} `json:"description" bson:"description"`
	
	DurationMinutes int `json:"duration_minutes" bson:"duration_minutes"`
	
	IconUrl string `json:"icon_url" bson:"icon_url"`
	
	Title string `json:"title" bson:"title"`
	
	
}

type ServicesRepository interface {
	List(ctx context.Context, limit, offset int) ([]*Services, error)
	Get(ctx context.Context, id string) (*Services, error)
	Create(ctx context.Context, model *Services) (string, error)
	Update(ctx context.Context, id string, model *Services) error
	Delete(ctx context.Context, id string) error
}
