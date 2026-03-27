package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/SilverName608/Effective-Mobile/internal/model"
	"github.com/SilverName608/Effective-Mobile/internal/service"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Create(ctx context.Context, sub *model.Subscription) error {
	return m.Called(ctx, sub).Error(0)
}

func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Subscription), args.Error(1)
}

func (m *mockRepo) List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*model.Subscription, error) {
	args := m.Called(ctx, userID, serviceName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Subscription), args.Error(1)
}

func (m *mockRepo) Update(ctx context.Context, sub *model.Subscription) error {
	return m.Called(ctx, sub).Error(0)
}

func (m *mockRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockRepo) TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) (int, error) {
	args := m.Called(ctx, userID, serviceName, from, to)
	return args.Int(0), args.Error(1)
}

func newTestService(repo *mockRepo) service.SubscriptionServiceI {
	log := logrus.New()
	log.SetLevel(logrus.PanicLevel)
	return service.NewSubscriptionService(repo, log)
}

func newSub() *model.Subscription {
	return &model.Subscription{
		ID:          uuid.New(),
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New(),
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

// --- Create ---

func TestCreate_Success(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	sub := newSub()

	repo.On("Create", mock.Anything, sub).Return(nil)

	err := svc.Create(context.Background(), sub)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCreate_RepoError(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	sub := newSub()

	repo.On("Create", mock.Anything, sub).Return(errors.New("db error"))

	err := svc.Create(context.Background(), sub)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

// --- GetByID ---

func TestGetByID_Success(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	sub := newSub()

	repo.On("GetByID", mock.Anything, sub.ID).Return(sub, nil)

	result, err := svc.GetByID(context.Background(), sub.ID)
	assert.NoError(t, err)
	assert.Equal(t, sub, result)
}

func TestGetByID_NotFound(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	id := uuid.New()

	repo.On("GetByID", mock.Anything, id).Return(nil, nil)

	result, err := svc.GetByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestGetByID_RepoError(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	id := uuid.New()

	repo.On("GetByID", mock.Anything, id).Return(nil, errors.New("db error"))

	result, err := svc.GetByID(context.Background(), id)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// --- List ---

func TestList_NoFilters(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	subs := []*model.Subscription{newSub(), newSub()}

	repo.On("List", mock.Anything, (*uuid.UUID)(nil), (*string)(nil)).Return(subs, nil)

	result, err := svc.List(context.Background(), nil, nil)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestList_WithUserID(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	sub := newSub()
	uid := sub.UserID

	repo.On("List", mock.Anything, &uid, (*string)(nil)).Return([]*model.Subscription{sub}, nil)

	result, err := svc.List(context.Background(), &uid, nil)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestList_WithServiceName(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	name := "Yandex Plus"
	subs := []*model.Subscription{newSub()}

	repo.On("List", mock.Anything, (*uuid.UUID)(nil), &name).Return(subs, nil)

	result, err := svc.List(context.Background(), nil, &name)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestList_RepoError(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)

	repo.On("List", mock.Anything, (*uuid.UUID)(nil), (*string)(nil)).Return(nil, errors.New("db error"))

	result, err := svc.List(context.Background(), nil, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// --- Update ---

func TestUpdate_Success(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	sub := newSub()

	repo.On("GetByID", mock.Anything, sub.ID).Return(sub, nil)
	repo.On("Update", mock.Anything, sub).Return(nil)

	err := svc.Update(context.Background(), sub)
	assert.NoError(t, err)
}

func TestUpdate_NotFound(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	sub := newSub()

	repo.On("GetByID", mock.Anything, sub.ID).Return(nil, nil)

	err := svc.Update(context.Background(), sub)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestUpdate_RepoError(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	sub := newSub()

	repo.On("GetByID", mock.Anything, sub.ID).Return(sub, nil)
	repo.On("Update", mock.Anything, sub).Return(errors.New("db error"))

	err := svc.Update(context.Background(), sub)
	assert.Error(t, err)
}

// --- Delete ---

func TestDelete_Success(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	sub := newSub()

	repo.On("GetByID", mock.Anything, sub.ID).Return(sub, nil)
	repo.On("Delete", mock.Anything, sub.ID).Return(nil)

	err := svc.Delete(context.Background(), sub.ID)
	assert.NoError(t, err)
}

func TestDelete_NotFound(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	id := uuid.New()

	repo.On("GetByID", mock.Anything, id).Return(nil, nil)

	err := svc.Delete(context.Background(), id)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDelete_RepoError(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	sub := newSub()

	repo.On("GetByID", mock.Anything, sub.ID).Return(sub, nil)
	repo.On("Delete", mock.Anything, sub.ID).Return(errors.New("db error"))

	err := svc.Delete(context.Background(), sub.ID)
	assert.Error(t, err)
}

// --- TotalCost ---

func TestTotalCost_Success(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)

	repo.On("TotalCost", mock.Anything, (*uuid.UUID)(nil), (*string)(nil), from, to).Return(1200, nil)

	total, err := svc.TotalCost(context.Background(), nil, nil, from, to)
	assert.NoError(t, err)
	assert.Equal(t, 1200, total)
}

func TestTotalCost_WithFilters(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	uid := uuid.New()
	name := "Yandex Plus"

	repo.On("TotalCost", mock.Anything, &uid, &name, from, to).Return(400, nil)

	total, err := svc.TotalCost(context.Background(), &uid, &name, from, to)
	assert.NoError(t, err)
	assert.Equal(t, 400, total)
}

func TestTotalCost_RepoError(t *testing.T) {
	repo := &mockRepo{}
	svc := newTestService(repo)
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)

	repo.On("TotalCost", mock.Anything, (*uuid.UUID)(nil), (*string)(nil), from, to).Return(0, errors.New("db error"))

	total, err := svc.TotalCost(context.Background(), nil, nil, from, to)
	assert.Error(t, err)
	assert.Equal(t, 0, total)
}
