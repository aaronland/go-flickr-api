package api

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
)

type Pagination struct {
	Page    int `json:"page"`
	Pages   int `json:"pages"`
	PerPage int `json:"perpage"`
	Total   int `json:total"`
}

func DerivePagination(ctx context.Context, fh io.ReadSeekCloser) (*Pagination, error) {

	body, err := io.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	page_rsp := gjson.GetBytes(body, "*.page")

	if !page_rsp.Exists() {
		return nil, fmt.Errorf("Unable to determine pagination properties (page) in response")
	}

	pages_rsp := gjson.GetBytes(body, "*.pages")

	if !pages_rsp.Exists() {
		return nil, fmt.Errorf("Unable to determine pagination properties (pages) in response")
	}

	perpage_rsp := gjson.GetBytes(body, "*.perpage")

	if !perpage_rsp.Exists() {
		return nil, fmt.Errorf("Unable to determine pagination properties (perpage) in response")
	}

	total_rsp := gjson.GetBytes(body, "*.total")

	if !total_rsp.Exists() {
		return nil, fmt.Errorf("Unable to determine pagination properties (total) in response")
	}

	pg := &Pagination{
		Page:    int(page_rsp.Int()),
		Pages:   int(pages_rsp.Int()),
		PerPage: int(perpage_rsp.Int()),
		Total:   int(total_rsp.Int()),
	}

	return pg, nil
}
