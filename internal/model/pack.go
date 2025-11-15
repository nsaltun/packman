package model

import "time"

type PackCalculationRequest struct {
	Quantity int `json:"quantity"`
}

// PackCalculationResponse represents the result of pack calculation
type PackCalculationResponse struct {
	Quantity int         `json:"quantity"`
	Packs    map[int]int `json:"packs"`
}

// GetPackSizesResponse represents the response for getting pack sizes
type GetPackSizesResponse struct {
	PackSizes []int     `json:"pack_sizes"`
	Version   int       `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by,omitempty"`
}

// UpdatePackSizesRequest represents a request to update pack sizes
type UpdatePackSizesRequest struct {
	PackSizes []int `json:"pack_sizes"`
}

// PackConfiguration represents the current pack size configuration
type PackConfiguration struct {
	ID        int       `json:"id" db:"id"`
	Version   int       `json:"version" db:"version"`
	PackSizes []int     `json:"pack_sizes"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy string    `json:"updated_by,omitempty" db:"updated_by"`
}
