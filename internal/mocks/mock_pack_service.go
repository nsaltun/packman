package mocks

import (
	"context"

	"github.com/nsaltun/packman/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockPackService struct {
	mock.Mock
}

func (m *MockPackService) CalculatePacks(ctx context.Context, quantity int) (*model.PackCalculationResponse, error) {
	args := m.Called(ctx, quantity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PackCalculationResponse), args.Error(1)
}

func (m *MockPackService) GetPackSizes(ctx context.Context) (*model.GetPackSizesResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.GetPackSizesResponse), args.Error(1)
}

func (m *MockPackService) UpdatePackSizes(ctx context.Context, sizes []int, updatedBy string) (*model.UpdatePackSizesResponse, error) {
	args := m.Called(ctx, sizes, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UpdatePackSizesResponse), args.Error(1)
}
