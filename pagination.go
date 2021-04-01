package api

import (
	"context"
	"fmt"
	"io"
)

type Pagination struct {
	Page    int `json:"page"`
	Pages   int `json:"pages"`
	PerPage int `json:"perpage"`
	Total   int `json:total"`
}

func DerivePagination(ctx context.Context, fh io.ReadSeekCloser) (*Pagination, error) {

	return nil, fmt.Errorf("Not implemented")
}
