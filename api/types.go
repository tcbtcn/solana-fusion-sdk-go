package api

// ApiConfig represents API configuration
type ApiConfig struct {
	BaseURL string
	AuthKey string
	Version string // "v1.0"
}

// Pagination represents paginated results
type Pagination[T any] struct {
	Meta  PaginationMeta `json:"meta"`
	Items []T            `json:"items"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	TotalItems   int `json:"totalItems"`
	ItemsPerPage int `json:"itemsPerPage"`
	TotalPages   int `json:"totalPages"`
	CurrentPage  int `json:"currentPage"`
}
