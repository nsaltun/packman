package service

import (
	"context"
	"testing"

	"github.com/nsaltun/packman/internal/apperror"
	"github.com/nsaltun/packman/internal/mocks"
	"github.com/nsaltun/packman/internal/model"
	"github.com/nsaltun/packman/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdatePackSizes(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		mockRepo := mocks.MockPackRepository{}
		service := packService{packRepo: &mockRepo}

		sizesToUpdate := []int{250, 500, 1000}
		updatedBy := "tester"
		expectedConfig := &model.PackConfiguration{
			ID:        1,
			Version:   2,
			PackSizes: sizesToUpdate,
			UpdatedBy: updatedBy,
		}

		mockRepo.On("UpdatePackSizes", mock.Anything, sizesToUpdate, updatedBy).Return(expectedConfig, nil)

		res, err := service.UpdatePackSizes(context.Background(), sizesToUpdate, updatedBy)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, sizesToUpdate, res.PackSizes)
		assert.Equal(t, updatedBy, res.UpdatedBy)
	})
	t.Run("repository error", func(t *testing.T) {
		mockRepo := mocks.MockPackRepository{}
		service := packService{packRepo: &mockRepo}
		sizesToUpdate := []int{250, 500, 1000}
		updatedBy := "tester"

		mockRepo.On("UpdatePackSizes", mock.Anything, sizesToUpdate, updatedBy).Return(nil, assert.AnError)
		res, err := service.UpdatePackSizes(context.Background(), sizesToUpdate, updatedBy)
		assert.Nil(t, res)
		assert.EqualError(t, err, apperror.InternalError("Failed to update pack sizes", assert.AnError).Error())
	})
	t.Run("repository not found error", func(t *testing.T) {
		mockRepo := mocks.MockPackRepository{}
		service := packService{packRepo: &mockRepo}
		sizesToUpdate := []int{250, 500, 1000}
		updatedBy := "tester"

		mockRepo.On("UpdatePackSizes", mock.Anything, sizesToUpdate, updatedBy).Return(nil, repository.ErrNotFound)
		res, err := service.UpdatePackSizes(context.Background(), sizesToUpdate, updatedBy)
		assert.Nil(t, res)
		assert.EqualError(t, err, apperror.NotFoundError("Pack configuration not found", repository.ErrNotFound).Error())
	})
}
