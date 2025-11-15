package handler

import (
	"fmt"

	"github.com/nsaltun/packman/internal/model"
)

const (
	maxQuantityLimit   = 10000000
	maxPackSizeLimit   = 1000000
	maxUpdatedByLength = 100
)

func validateCalculatePacksRequest(req *model.PackCalculationRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// validate quantity
	if req.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}

	// set a reasonable upper limit for quantity
	if req.Quantity > maxQuantityLimit {
		return fmt.Errorf("quantity must be less than or equal to %d", maxQuantityLimit)
	}

	return nil
}

func validateUpdatePackSizesRequest(req *model.UpdatePackSizesRequest) error {
	// validate request
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	// validate pack sizes
	if len(req.PackSizes) == 0 {
		return fmt.Errorf("pack_sizes cannot be empty")
	}
	// validate each pack size
	for _, size := range req.PackSizes {
		if size <= 0 {
			return fmt.Errorf("pack sizes must be greater than zero")
		}
		if size > maxPackSizeLimit {
			return fmt.Errorf("pack sizes must be less than or equal to %d", maxPackSizeLimit)
		}
	}
	// validate updated_by
	if len(req.UpdatedBy) > maxUpdatedByLength {
		return fmt.Errorf("updated_by must be less than or equal to %d characters", maxUpdatedByLength)
	}

	return nil
}
