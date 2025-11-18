package mocks

import (
	"context"

	"github.com/nsaltun/packman/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockPackRepository is a mock implementation of repository.PackRepository
type MockPackRepository struct {
	mock.Mock
}

// GetPackSizes mocks the GetPackSizes method
func (m *MockPackRepository) GetPackSizes(ctx context.Context) ([]int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

// GetPackConfiguration mocks the GetPackConfiguration method
func (m *MockPackRepository) GetPackConfiguration(ctx context.Context) (*model.PackConfiguration, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PackConfiguration), args.Error(1)
}

// UpdatePackSizes mocks the UpdatePackSizes method
func (m *MockPackRepository) UpdatePackSizes(ctx context.Context, sizes []int, updatedBy string) (*model.PackConfiguration, error) {
	args := m.Called(ctx, sizes, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PackConfiguration), args.Error(1)
}

// GetPackConfigurationHistory mocks the GetPackConfigurationHistory method
func (m *MockPackRepository) GetPackConfigurationHistory(ctx context.Context, limit int) ([]*model.PackConfiguration, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.PackConfiguration), args.Error(1)
}
