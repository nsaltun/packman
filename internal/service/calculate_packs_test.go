package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/nsaltun/packman/internal/apperror"
	"github.com/nsaltun/packman/internal/mocks"
	"github.com/nsaltun/packman/internal/model"
	"github.com/nsaltun/packman/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCalculatePacks(t *testing.T) {
	tests := []struct {
		name      string
		packSizes []int
		quantity  int
		expected  *model.PackCalculationResponse
		repoErr   error
		wantErr   error
	}{
		{
			name:      "Exact match with multiple pack sizes",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			quantity:  12000,
			expected: &model.PackCalculationResponse{
				Quantity: 12000,
				Packs: map[int]int{
					5000: 2,
					2000: 1,
				},
			},
		},
		{
			name:      "Exact match with single pack size",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			quantity:  10000,
			expected: &model.PackCalculationResponse{
				Quantity: 10000,
				Packs: map[int]int{
					5000: 2,
				},
			},
		},
		{
			name:      "No exact match",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			quantity:  7500,
			expected: &model.PackCalculationResponse{
				Quantity: 7500,
				Packs: map[int]int{
					5000: 1,
					2000: 1,
					500:  1,
				},
			},
		},
		{
			name:      "Quantity less than smallest pack size",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			quantity:  100,
			expected: &model.PackCalculationResponse{
				Quantity: 100,
				Packs:    map[int]int{250: 1},
			},
		},
		{
			name:      "edge case",
			packSizes: []int{23, 31, 53},
			quantity:  500000,
			expected: &model.PackCalculationResponse{
				Quantity: 500000,
				Packs: map[int]int{
					53: 9433,
					31: 1,
					23: 1,
				},
			},
		},
		{
			name:      "pack sizes empty",
			packSizes: []int{},
			quantity:  100,
			wantErr:   apperror.InternalError("Pack sizes configuration is empty", nil),
		},
		{
			name:      "repository error",
			packSizes: nil,
			quantity:  1000,
			repoErr:   fmt.Errorf("database error"),
			wantErr:   apperror.InternalError("Failed to retrieve pack sizes", fmt.Errorf("database error")),
		},
		{
			name:      "repository notfound error",
			packSizes: nil,
			quantity:  1000,
			repoErr:   repository.ErrNotFound,
			wantErr:   apperror.NotFoundError("Pack configuration not found", repository.ErrNotFound),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//setup
			repoMock := mocks.MockPackRepository{}
			service := packService{packRepo: &repoMock}
			repoMock.On("GetPackSizes", mock.Anything).Return(tt.packSizes, tt.repoErr) // Reset mock for each test

			//execute
			res, err := service.CalculatePacks(context.Background(), tt.quantity)

			//verify
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.expected, res)
		})
	}
}
