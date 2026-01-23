package domain

import (
	"context"
	
)

type Providers struct {
	ID string `json:"id" bson:"_id,omitempty"`
	
	Address string `json:"address" bson:"address"`
	
	AvatarUrl string `json:"avatar_url" bson:"avatar_url"`
	
	EstablishmentName string `json:"establishment_name" bson:"establishment_name"`
	
	FullName string `json:"full_name" bson:"full_name"`
	
	
}

type ProvidersRepository interface {
	List(ctx context.Context, limit, offset int) ([]*Providers, error)
	Get(ctx context.Context, id string) (*Providers, error)
	Create(ctx context.Context, model *Providers) (string, error)
	Update(ctx context.Context, id string, model *Providers) error
	Delete(ctx context.Context, id string) error
}
