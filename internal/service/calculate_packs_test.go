package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/nsaltun/packman/internal/mocks"
	"github.com/nsaltun/packman/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCalculatePacks(t *testing.T) {
	tests := []struct {
		name      string
		packSizes []int
		quantity  int
		expected  model.PackCalculationResponse
		wantErr   error
	}{
		{
			name:      "Exact match with multiple pack sizes",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			quantity:  12000,
			expected: model.PackCalculationResponse{
				Quantity: 12000,
				Packs: map[int]int{
					5000: 2,
					2000: 1,
					1000: 1,
				},
			},
		},
		{
			name:      "Exact match with single pack size",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			quantity:  10000,
			expected: model.PackCalculationResponse{
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
			expected: model.PackCalculationResponse{
				Quantity: 7500,
				Packs: map[int]int{
					5000: 1,
					2000: 1,
					1000: 1,
				},
			},
		},
		{
			name:      "Quantity less than smallest pack size",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			quantity:  100,
			expected: model.PackCalculationResponse{
				Quantity: 100,
				Packs:    map[int]int{100: 1},
			},
		},
		{
			name:      "Zero quantity",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			quantity:  0,
			wantErr:   fmt.Errorf("quantity must be greater than zero"),
		},
		{
			name:      "pack sizes empty",
			packSizes: []int{},
			quantity:  100,
			wantErr:   fmt.Errorf("pack sizes empty. Cannot calculate packs"),
		},
		{
			name:      "edge case",
			packSizes: []int{23, 31, 53},
			quantity:  500000,
			expected: model.PackCalculationResponse{
				Quantity: 500000,
				Packs: map[int]int{
					53: 9429,
					31: 7,
					23: 2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//setup
			repoMock := mocks.MockPackRepository{}
			service := packService{packRepo: &repoMock}
			repoMock.On("GetPackSizes", mock.Anything).Return(tt.packSizes, nil) // Reset mock for each test

			//execute
			res, err := service.CalculatePacks(context.Background(), tt.quantity)

			//verify
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.expected, res)
		})
	}
}
