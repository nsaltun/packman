package service

import (
	"context"
	"errors"
	"sort"

	"github.com/nsaltun/packman/internal/apperror"
	"github.com/nsaltun/packman/internal/model"
	"github.com/nsaltun/packman/internal/repository"
)

// PackService defines the interface for pack-related operations
type PackService interface {
	CalculatePacks(ctx context.Context, quantity int) (*model.PackCalculationResponse, error)
	GetPackSizes(ctx context.Context) (*model.GetPackSizesResponse, error)
	UpdatePackSizes(ctx context.Context, sizes []int, updatedBy string) error
}

// packService is the concrete implementation of PackService
type packService struct {
	packRepo repository.PackRepository
}

// NewPackService creates a new instance of PackService
func NewPackService(packRepo repository.PackRepository) PackService {
	return &packService{packRepo: packRepo}
}

// CalculatePacks calculates the optimal combination of packs for a given quantity
func (s *packService) CalculatePacks(ctx context.Context, quantity int) (*model.PackCalculationResponse, error) {
	// get pack sizes from repository
	packSizes, err := s.packRepo.GetPackSizes(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.NotFoundError("Pack configuration not found", err)
		}
		return nil, apperror.InternalError("Failed to retrieve pack sizes", err)
	}

	// implement calculation logic
	if len(packSizes) == 0 {
		return nil, apperror.InternalError("Pack sizes configuration is empty", nil)
	}

	// sort pack sizes in descending order
	sort.Slice(packSizes, func(i, j int) bool {
		return packSizes[i] > packSizes[j]
	})

	//-- Callculate packs --

	//Create a copy of quantity to return later
	originalQuantity := quantity

	// get the smallest pack size
	minPackSize := packSizes[len(packSizes)-1]

	// greedy algorithm to find the combination of packs
	packsNumberResult := make(map[int]int)
	// iterate over pack sizes
	for _, packSize := range packSizes {
		// if quantity is zero, break
		if quantity == 0 {
			break
		}
		// calculate number of packs for current pack size
		if quantity >= packSize {
			// integer division to get number of packs
			numPacks := quantity / packSize
			// update result map
			packsNumberResult[packSize] = numPacks
			// update remaining quantity
			quantity = quantity % packSize
		}
	}

	// if there is remaining quantity less than the smallest pack size, add one smallest pack
	if quantity > 0 {
		packsNumberResult[minPackSize]++
	}

	// return result
	return &model.PackCalculationResponse{Packs: packsNumberResult, Quantity: originalQuantity}, nil
}

// GetPackSizes retrieves the current pack sizes from the repository
func (s *packService) GetPackSizes(ctx context.Context) (*model.GetPackSizesResponse, error) {
	res, err := s.packRepo.GetPackConfiguration(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.NotFoundError("Pack configuration not found", err)
		}
		return nil, apperror.InternalError("Failed to retrieve pack configuration", err)
	}

	return &model.GetPackSizesResponse{
		PackSizes: res.PackSizes,
		UpdatedAt: res.UpdatedAt,
		UpdatedBy: res.UpdatedBy,
		Version:   res.Version,
	}, nil
}

// UpdatePackSizes updates the pack sizes in the repository
func (s *packService) UpdatePackSizes(ctx context.Context, sizes []int, updatedBy string) error {
	err := s.packRepo.UpdatePackSizes(ctx, sizes, updatedBy)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.NotFoundError("Pack configuration not found", err)
		}
		return apperror.InternalError("Failed to update pack sizes", err)
	}

	return nil
}
