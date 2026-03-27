package service

import (
	"context"
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
	return nil
}

func (s *subscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	return nil, nil
}

func (s *subscriptionService) List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*model.Subscription, error) {
	return nil, nil
}

func (s *subscriptionService) Update(ctx context.Context, sub *model.Subscription) error {
	return nil
}

func (s *subscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (s *subscriptionService) TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) (int, error) {
	return 0, nil
}
