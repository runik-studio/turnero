package domain

import (
	"context"
	"time"
)

type Users struct {
	ID string `json:"id" bson:"_id,omitempty" firestore:"-"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`

	Email string `json:"email" bson:"email"`

	Name string `json:"name" bson:"name"`

	Picture string `json:"picture" bson:"picture"`

	RoleId string `json:"role_id" bson:"role_id"`

	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type UsersRepository interface {
	List(ctx context.Context, limit, offset int) ([]*Users, error)
	Get(ctx context.Context, id string) (*Users, error)
	Create(ctx context.Context, model *Users) (string, error)
	Update(ctx context.Context, id string, model *Users) error
	Delete(ctx context.Context, id string) error
}
