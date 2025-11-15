package service

import (
	"context"
	"errors"

	"github.com/nsaltun/packman/internal/model"
	"github.com/nsaltun/packman/internal/repository"
)

type PackService interface {
	CalculatePacks(ctx context.Context, quantity int) (*model.PackCalculationResponse, error)
	GetPackSizes(ctx context.Context) (*model.GetPackSizesResponse, error)
	UpdatePackSizes(ctx context.Context, sizes []int, updatedBy string) error
}

type packService struct {
	packRepo repository.PackRepository
}

func NewPackService(packRepo repository.PackRepository) PackService {
	return &packService{packRepo: packRepo}
}

func (s *packService) CalculatePacks(ctx context.Context, quantity int) (*model.PackCalculationResponse, error) {
	// get pack sizes from repository
	packSizes, err := s.packRepo.GetPackSizes(ctx)
	if err != nil {
		return nil, err
	}

	// implement calculation logic
	if len(packSizes) == 0 {
		return nil, errors.New("no pack sizes available")
	}

	// TODO: Implement pack calculation logic

	return nil, nil
}

func (s *packService) GetPackSizes(ctx context.Context) (*model.GetPackSizesResponse, error) {
	return nil, nil
}

func (s *packService) UpdatePackSizes(ctx context.Context, sizes []int, updatedBy string) error {
	return nil
}
