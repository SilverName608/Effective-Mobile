package repository

import (
	"context"
	"errors"
	"fmt"
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
	err := r.db.WithContext(ctx).Create(sub).Error
	if err != nil {
		return fmt.Errorf("create subscription: %w", err)
	}
	return nil
}

func (r *subscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	var sub model.Subscription
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&sub).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}
	return &sub, nil
}

func (r *subscriptionRepository) List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*model.Subscription, error) {
	var subs []*model.Subscription

	q := r.db.WithContext(ctx).Model(&model.Subscription{})

	if userID != nil {
		q = q.Where("user_id = ?", *userID)
	}
	if serviceName != nil {
		q = q.Where("service_name ILIKE ?", *serviceName)
	}

	if err := q.Order("created_at DESC").Find(&subs).Error; err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}

	return subs, nil
}

func (r *subscriptionRepository) Update(ctx context.Context, sub *model.Subscription) error {
	err := r.db.WithContext(ctx).Save(sub).Error
	if err != nil {
		return fmt.Errorf("update subscription: %w", err)
	}
	return nil
}

func (r *subscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Subscription{}).Error; err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}
	return nil
}

func (r *subscriptionRepository) TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) (int, error) {
	q := r.db.WithContext(ctx).Model(&model.Subscription{}).
		Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", to, from)

	if userID != nil {
		q = q.Where("user_id = ?", *userID)
	}
	if serviceName != nil {
		q = q.Where("service_name ILIKE ?", *serviceName)
	}

	var total int
	err := q.Select("COALESCE(SUM(price), 0)").Scan(&total).Error
	if err != nil {
		return 0, fmt.Errorf("total cost: %w", err)
	}
	return total, nil
}
