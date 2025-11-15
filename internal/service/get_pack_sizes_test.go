package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nsaltun/packman/internal/mocks"
	"github.com/nsaltun/packman/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPackSizes(t *testing.T) {
	tests := []struct {
		name            string
		expected        *model.GetPackSizesResponse
		inRepo          *model.PackConfiguration
		repoErr         error
		wantErrorAssert func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool
	}{
		{
			name: "Successful retrieval of pack sizes",
			expected: &model.GetPackSizesResponse{
				PackSizes: []int{250, 500, 1000},
				Version:   1,
				UpdatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
				UpdatedBy: "admin",
			},
			inRepo: &model.PackConfiguration{
				PackSizes: []int{250, 500, 1000},
				Version:   1,
				UpdatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
				UpdatedBy: "admin",
			},
			wantErrorAssert: assert.NoError,
		},
		{
			name:    "Repository error",
			inRepo:  nil,
			repoErr: fmt.Errorf("database error"),
			wantErrorAssert: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.EqualError(t, err, "database error", msgAndArgs...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//setup
			repoMock := mocks.MockPackRepository{}
			service := packService{packRepo: &repoMock}
			repoMock.On("GetPackConfiguration", mock.Anything).Return(tt.inRepo, tt.repoErr)

			//execute
			res, err := service.GetPackSizes(context.Background())

			//verify
			tt.wantErrorAssert(t, err)
			assert.Equal(t, tt.expected, res)
		})
	}
}
