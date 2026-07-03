package core_pagination

import (
	"fmt"

	core_errors "github.com/vladislav-koval/todo-app/internal/core/errors"
)

type Pagination struct {
	Limit  int
	Offset int
}

const (
	DefaultLimit  = 10
	DefaultOffset = 0
	MaxLimit      = 100
)

func NewPagination(limit *int, offset *int) (Pagination, error) {
	result := Pagination{
		Limit:  DefaultLimit,
		Offset: DefaultOffset,
	}

	if limit != nil {
		if *limit < 0 {
			return Pagination{}, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
		}

		if *limit > MaxLimit {
			return Pagination{}, fmt.Errorf("limit must be less than or equal to %d: %w", MaxLimit, core_errors.ErrInvalidArgument)
		}

		result.Limit = *limit
	}

	if offset != nil {
		if *offset < 0 {
			return Pagination{}, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
		}

		result.Offset = *offset
	}

	return result, nil
}
