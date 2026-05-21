package models

import "time"

// ID is a unique identifier for an entity.
type ID string

type PaginationType string

const (
	PaginationTypeTime      PaginationType = "time"
	PaginationTypePage      PaginationType = "page"
	PaginationTypeOffset    PaginationType = "offset"
	PaginationTypeCursor    PaginationType = "cursor"
	PaginationTypeKeySet    PaginationType = "keyset"
	PaginationTypeUndefined PaginationType = "undefined"
	PaginationTypeEmpty     PaginationType = ""
)

var ValidPaginationTypes = []PaginationType{PaginationTypeTime, PaginationTypePage, PaginationTypeOffset, PaginationTypeCursor, PaginationTypeKeySet}

type Sort string

const (
	SortAsc       Sort = "asc"
	SortDesc      Sort = "desc"
	SortUndefined Sort = "undefined"
	SortEmpty     Sort = ""
)

type Align string

const (
	AlignStart     Align = "start"
	AlignEnd       Align = "end"
	AlignUndefined Align = "undefined"
	AlignEmpty     Align = ""
)

type PaginationMode string

const (
	PaginationModeAuto      PaginationMode = "auto"
	PaginationModeManual    PaginationMode = "manual"
	PaginationModeUndefined PaginationMode = "undefined"
)

type Page struct {
	Type PaginationType `json:"type"`

	// page-based pagination
	CurrentPage *int `json:"currentPage,omitempty"`
	PageSize    int  `json:"pageSize,omitempty"`
	TotalPages  int  `json:"totalPages,omitempty"`

	// offset-based pagination
	Offset int `json:"offset,omitempty"`

	// cursor-based pagination
	PrevCursor string `json:"prevCursor,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`

	// keyset-based pagination
	ID string `json:"id,omitempty"`

	// common
	Limit      int `json:"limit,omitempty"`
	TotalItems int `json:"totalItems,omitempty"`
}

type Pagination[T any] struct {
	Data []T   `json:"items"`
	Page *Page `json:"page"`
}

type Paginatable interface {
	GetTime() time.Time
	GetID() ID
}

type Paginatables []Paginatable

func (e Paginatables) Len() int {
	return len(e)
}

func (e Paginatables) Less(i, j int) bool {
	return e[j].GetTime().After(e[i].GetTime())
}

func (e Paginatables) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
