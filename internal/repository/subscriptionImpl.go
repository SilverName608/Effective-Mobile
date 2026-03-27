package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/SilverName608/Effective-Mobile/internal/model"
)

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepositoryI {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(ctx context.Context, sub *model.Subscription) error {
	return nil
}

func (r *subscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	return nil, nil
}

func (r *subscriptionRepository) List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*model.Subscription, error) {
	return nil, nil
}

func (r *subscriptionRepository) Update(ctx context.Context, sub *model.Subscription) error {
	return nil
}

func (r *subscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r *subscriptionRepository) TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) (int, error) {
	return 0, nil
}
