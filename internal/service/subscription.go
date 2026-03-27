package service

import (
	"context"
	"time"

	"github.com/SilverName608/Effective-Mobile/internal/model"
	"github.com/google/uuid"
)

type SubscriptionServiceI interface {
	Create(ctx context.Context, sub *model.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*model.Subscription, error)
	Update(ctx context.Context, sub *model.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) (int, error)
}
