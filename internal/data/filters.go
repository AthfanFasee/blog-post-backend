package data

import (
	"math"
	"strings"

	"github.com/AthfanFasee/blog-post-backend/internal/validator"
)

type Filters struct {
	ID           int
	Page         int
	Limit        int
	Sort         string
	SortSafeList []string
}

func (f Filters) sortParam() string {
	var sortParam string
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			switch {
			case f.Sort == "-likescount":
				sortParam = "ARRAY_LENGTH(liked_by, 1)"
				return sortParam
			default:
				sortParam = strings.TrimPrefix(f.Sort, "-")
				return sortParam
			}
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	switch {
	case f.Sort == "-likescount":
		return "DESC NULLS LAST"
	case strings.HasPrefix(f.Sort, "-"):
		return "DESC"
	default:
		return "ASC"
	}
}

func (f Filters) limit() int {
	return f.Limit
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.Limit
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.ID >= 0, "id", "must be greater than zero")
	// Limit the maximum page number to prevent integer overflow
	v.Check(f.Page <= 1_000_000, "page", "must be a maximum of 1 million")
	v.Check(f.Limit > 0, "limit", "must be greater than zero")
	v.Check(f.Limit <= 100, "limit", "must be a maximum of 100")

	// Check that the sort parameter matches a value in safelist
	v.Check(validator.In(f.Sort, f.SortSafeList...), "sort", "invalid sort value")

}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func calculateMetadata(totalRecords, page, limit int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     limit,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(limit))),
		TotalRecords: totalRecords,
	}
}
