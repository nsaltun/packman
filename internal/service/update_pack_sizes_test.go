package service

import (
	"context"
	"testing"

	"github.com/nsaltun/packman/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdatePackSizes(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		mockRepo := mocks.MockPackRepository{}
		service := packService{packRepo: &mockRepo}

		sizesToUpdate := []int{250, 500, 1000}
		updatedBy := "tester"

		mockRepo.On("UpdatePackSizes", mock.Anything, sizesToUpdate, updatedBy).Return(nil)

		err := service.UpdatePackSizes(context.Background(), sizesToUpdate, updatedBy)
		assert.NoError(t, err)
	})
	t.Run("repository error", func(t *testing.T) {
		mockRepo := mocks.MockPackRepository{}
		service := packService{packRepo: &mockRepo}
		sizesToUpdate := []int{250, 500, 1000}
		updatedBy := "tester"

		mockRepo.On("UpdatePackSizes", mock.Anything, sizesToUpdate, updatedBy).Return(assert.AnError)
		err := service.UpdatePackSizes(context.Background(), sizesToUpdate, updatedBy)
		assert.EqualError(t, err, assert.AnError.Error())
	})
}
