package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/SilverName608/Effective-Mobile/internal/model"
	"github.com/SilverName608/Effective-Mobile/internal/repository"
)

type subscriptionService struct {
	repo repository.SubscriptionRepositoryI
	log  *logrus.Logger
}

func NewSubscriptionService(repo repository.SubscriptionRepositoryI, log *logrus.Logger) SubscriptionServiceI {
	return &subscriptionService{repo: repo, log: log}
}

func (s *subscriptionService) Create(ctx context.Context, sub *model.Subscription) error {
	err := s.repo.Create(ctx, sub)
	if err != nil {
		s.log.WithError(err).Error("create subscription failed")
		return fmt.Errorf("create subscription: %w", err)
	}
	s.log.WithField("id", sub.ID).Info("subscription created")
	return nil
}

func (s *subscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.WithError(err).Error("get subscription failed")
		return nil, fmt.Errorf("get subscription: %w", err)
	}
	return sub, nil
}

func (s *subscriptionService) List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*model.Subscription, error) {
	subs, err := s.repo.List(ctx, userID, serviceName)
	if err != nil {
		s.log.WithError(err).Error("list subscriptions failed")
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}
	return subs, nil
}

func (s *subscriptionService) Update(ctx context.Context, sub *model.Subscription) error {
	existing, err := s.repo.GetByID(ctx, sub.ID)
	if err != nil {
		return fmt.Errorf("get subscription: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("subscription not found")
	}

	err = s.repo.Update(ctx, sub)
	if err != nil {
		s.log.WithError(err).Error("update subscription failed")
		return fmt.Errorf("update subscription: %w", err)
	}
	s.log.WithField("id", sub.ID).Info("subscription updated")
	return nil
}

func (s *subscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get subscription: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("subscription not found")
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		s.log.WithError(err).Error("delete subscription failed")
		return fmt.Errorf("delete subscription: %w", err)
	}
	s.log.WithField("id", id).Info("subscription deleted")
	return nil
}

func (s *subscriptionService) TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) (int, error) {
	total, err := s.repo.TotalCost(ctx, userID, serviceName, from, to)
	if err != nil {
		s.log.WithError(err).Error("total cost failed")
		return 0, fmt.Errorf("total cost: %w", err)
	}
	return total, nil
}
